package csharp

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

// Converters emits one `<T>JsonConverter.cs` per disjunction-derived
// struct, providing System.Text.Json (de)serialization for types that
// can't be expressed as plain POCOs (sum types).
//
// Three disjunction shapes are supported, mirroring the Java jenny:
//   - HintDisjunctionOfScalars         (e.g. `string | bool | int64`)
//   - HintDiscriminatedDisjunctionOfRefs (e.g. `KindA | KindB`)
//   - HintDisjunctionOfScalarsAndRefs   (e.g. `string | KindA`)
type Converters struct {
	config        Config
	tmpl          *template.Template
	context       languages.Context
	typeFormatter *typeFormatter
}

func (jenny *Converters) JennyName() string {
	return "CSharpJsonConverters"
}

func (jenny *Converters) Generate(context languages.Context) (codejen.Files, error) {
	jenny.context = context
	jenny.typeFormatter = newTypeFormatter(context, jenny.config, nil)

	files := make(codejen.Files, 0)
	for _, schema := range context.Schemas {
		var iterErr error
		schema.Objects.Iterate(func(_ string, object ast.Object) {
			if iterErr != nil || !needsCustomConverter(object) {
				return
			}
			f, err := jenny.genConverter(schema.Package, object)
			if err != nil {
				iterErr = err
				return
			}
			files = append(files, *f)
		})
		if iterErr != nil {
			return nil, iterErr
		}
	}
	return files, nil
}

func (jenny *Converters) genConverter(pkg string, object ast.Object) (*codejen.File, error) {
	namespace := jenny.config.formatNamespace(pkg)
	pkgFolder := formatPackageName(pkg)
	name := formatObjectName(object.Name)
	filename := filepath.Join(jenny.config.ProjectPath, pkgFolder, name+"JsonConverter.cs")

	// Imports collected during type formatting are discarded for now —
	// converter templates emit their `using` directives statically.
	tf := jenny.typeFormatter.withImports(newImportMap(jenny.config.NamespaceRoot, namespace))

	t := object.Type
	switch {
	case t.HasHint(ast.HintDisjunctionOfScalars):
		out, err := jenny.renderScalars(namespace, name, object, tf)
		if err != nil {
			return nil, err
		}
		return codejen.NewFile(filename, out, jenny), nil
	case t.HasHint(ast.HintDiscriminatedDisjunctionOfRefs):
		out, err := jenny.renderRefs(namespace, name, object)
		if err != nil {
			return nil, err
		}
		return codejen.NewFile(filename, out, jenny), nil
	case t.HasHint(ast.HintDisjunctionOfScalarsAndRefs):
		out, err := jenny.renderScalarsAndRefs(namespace, name, object, tf)
		if err != nil {
			return nil, err
		}
		return codejen.NewFile(filename, out, jenny), nil
	}
	return nil, fmt.Errorf("converter requested for non-disjunction object %q", object.Name)
}

// scalarBranchModel describes one branch of a scalar disjunction for
// the converter template. After grouping, multiple source branches
// may share the same token (e.g. several numeric kinds → `Number`),
// in which case the int/float branches are folded into the same case
// body and dispatched at runtime via `reader.TryGetInt64`.
type scalarBranchModel struct {
	Field    string   // formatted C# field name on the parent struct
	Tokens   []string // JsonTokenType members that select this branch (Utf8JsonReader path)
	ReadExpr string   // expression used to read the value from the reader
	// WriteAccess is appended to the field reference in the Write
	// path, e.g. ".Value" so we serialise the underlying value of a
	// `bool?` rather than the Nullable<bool> wrapper.
	WriteAccess string
	// Numeric grouping companions for the Utf8JsonReader path.
	IntField    string
	IntReadExpr string
	IntWrite    string

	// JsonElement-path fields, populated by scalarBranchesForElement
	// for templates that read via a JsonDocument tree.
	ValueKinds    []string // JsonValueKind members that select this branch
	ElementGet    string   // method-call segment after `root.` (e.g. "GetString()")
	IntElementGet string   // sibling Get* used for the integer branch in numeric grouping
	DeserializeAs string   // when ElementGet is empty, deserialise raw JSON into this type
}

func (jenny *Converters) renderScalars(namespace, name string, object ast.Object, tf *typeFormatter) ([]byte, error) {
	branches := scalarBranches(object.Type.AsStruct().Fields, tf)
	return jenny.tmpl.RenderAsBytes("marshalling/disjunction_of_scalars.tmpl", map[string]any{
		"Namespace": namespace,
		"Name":      name,
		"Branches":  branches,
	})
}

// refCaseModel describes one entry in a discriminated-disjunction
// switch.
type refCaseModel struct {
	Discriminator string // discriminator value, ignored when CatchAll is true
	CatchAll      bool
	Field         string
	TypeRef       string // C# type expression for the branch
}

func (jenny *Converters) renderRefs(namespace, name string, object ast.Object) ([]byte, error) {
	disjunctionHint, ok := object.Type.Hints[ast.HintDiscriminatedDisjunctionOfRefs].(ast.DisjunctionType)
	if !ok {
		return nil, fmt.Errorf("disjunction hint missing or of wrong type for %q", object.Name)
	}

	cases := make([]refCaseModel, 0, len(disjunctionHint.DiscriminatorMapping))
	hasCatchAll := false
	for discValue, refType := range disjunctionHint.DiscriminatorMapping {
		isCatchAll := discValue == "cog_discriminator_catch_all"
		if isCatchAll {
			hasCatchAll = true
		}
		cases = append(cases, refCaseModel{
			Discriminator: discValue,
			CatchAll:      isCatchAll,
			Field:         formatFieldName(refType),
			TypeRef:       formatObjectName(refType),
		})
	}
	// Catch-all entry sorts last so the `default:` label appears at the
	// bottom of the switch; remaining entries sorted alphabetically for
	// stable output across runs.
	sort.SliceStable(cases, func(i, j int) bool {
		if cases[i].CatchAll != cases[j].CatchAll {
			return !cases[i].CatchAll
		}
		return cases[i].Discriminator < cases[j].Discriminator
	})

	return jenny.tmpl.RenderAsBytes("marshalling/disjunction_of_refs.tmpl", map[string]any{
		"Namespace":     namespace,
		"Name":          name,
		"Discriminator": disjunctionHint.Discriminator,
		"Cases":         cases,
		"HasCatchAll":   hasCatchAll,
	})
}

func (jenny *Converters) renderScalarsAndRefs(namespace, name string, object ast.Object, tf *typeFormatter) ([]byte, error) {
	scalarFields := make([]ast.StructField, 0)
	refFields := make([]ast.StructField, 0)
	var anyField *ast.StructField
	for _, field := range object.Type.AsStruct().Fields {
		f := field
		switch {
		case f.Type.IsScalar() && f.Type.AsScalar().ScalarKind == ast.KindAny:
			// `any`-typed branches are dispatched as a JSON-object
			// catch-all rather than via TokenType, since `object`
			// would otherwise shadow every ref branch.
			anyField = &f
		case f.Type.IsScalar() || f.Type.IsArray():
			scalarFields = append(scalarFields, f)
		case f.Type.IsRef():
			refFields = append(refFields, f)
		}
	}

	scalarBranches := scalarBranchesForElement(scalarFields, tf)

	refModels := make([]refCaseModel, 0, len(refFields))
	for _, field := range refFields {
		refModels = append(refModels, refCaseModel{
			Field:   formatFieldName(field.Name),
			TypeRef: tf.formatFieldType(field.Type),
		})
	}

	data := map[string]any{
		"Namespace":      namespace,
		"Name":           name,
		"ScalarBranches": scalarBranches,
		"RefBranches":    refModels,
		"HasAnyBranch":   anyField != nil,
	}
	if anyField != nil {
		data["AnyField"] = formatFieldName(anyField.Name)
	}

	return jenny.tmpl.RenderAsBytes("marshalling/disjunction_of_scalars_and_refs.tmpl", data)
}

// scalarBranches builds dispatch metadata for scalar/array fields.
// Branches sharing the same JsonTokenType are merged so the generated
// `switch` has unique case labels. Numeric kinds (`Number`) are
// disambiguated at runtime via `reader.TryGetInt64`: an integer-typed
// branch wins when present, with the floating-point branch as the
// fallback. Subsequent same-token branches beyond the first two are
// dropped — disjunctions of three+ numeric kinds aren't expressible
// without losing information from JSON anyway.
func scalarBranches(fields []ast.StructField, tf *typeFormatter) []scalarBranchModel {
	out := make([]scalarBranchModel, 0, len(fields))
	tokenIdx := make(map[string]int, len(fields))

	for _, field := range fields {
		fieldName := formatFieldName(field.Name)
		var b scalarBranchModel
		switch {
		case field.Type.IsArray():
			elementType := tf.formatFieldType(field.Type.AsArray().ValueType)
			b = scalarBranchModel{
				Field:    fieldName,
				Tokens:   []string{"StartArray"},
				ReadExpr: fmt.Sprintf("JsonSerializer.Deserialize<List<%s>>(ref reader, options)", elementType),
			}
		case field.Type.IsScalar():
			tokens, readExpr, writeAccess := scalarBranchInfo(field.Type.AsScalar().ScalarKind)
			b = scalarBranchModel{
				Field:       fieldName,
				Tokens:      tokens,
				ReadExpr:    readExpr,
				WriteAccess: writeAccess,
			}
		default:
			continue
		}

		key := tokenKey(b.Tokens)
		if idx, exists := tokenIdx[key]; exists && key == "Number" {
			existing := out[idx]
			// Promote the integer branch into the primary slot when the
			// existing branch is a float, so `TryGetInt64` is checked
			// first.
			if isIntegerRead(b.ReadExpr) && !isIntegerRead(existing.ReadExpr) {
				existing.IntField = b.Field
				existing.IntReadExpr = b.ReadExpr
				existing.IntWrite = b.WriteAccess
			} else {
				existing.IntField = existing.Field
				existing.IntReadExpr = existing.ReadExpr
				existing.IntWrite = existing.WriteAccess
				existing.Field = b.Field
				existing.ReadExpr = b.ReadExpr
				existing.WriteAccess = b.WriteAccess
			}
			out[idx] = existing
			continue
		}
		if _, exists := tokenIdx[key]; exists {
			// Non-numeric duplicate (e.g. two `String` branches): drop
			// silently — first branch wins.
			continue
		}
		tokenIdx[key] = len(out)
		out = append(out, b)
	}
	return out
}

func tokenKey(tokens []string) string {
	if len(tokens) == 1 {
		return tokens[0]
	}
	// Multi-token branches (e.g. True+False for bool) get a unique key.
	key := tokens[0]
	for _, t := range tokens[1:] {
		key += "|" + t
	}
	return key
}

func isIntegerRead(expr string) bool {
	switch expr {
	case
		"reader.GetSByte()",
		"reader.GetByte()",
		"reader.GetInt16()",
		"reader.GetUInt16()",
		"reader.GetInt32()",
		"reader.GetUInt32()",
		"reader.GetInt64()",
		"reader.GetUInt64()":
		return true
	}
	return false
}

// scalarBranchesForElement is the JsonElement-flavoured analogue of
// scalarBranches. Used by the scalars-and-refs template, which reads
// the input into a JsonDocument before dispatching.
func scalarBranchesForElement(fields []ast.StructField, tf *typeFormatter) []scalarBranchModel {
	out := make([]scalarBranchModel, 0, len(fields))
	kindIdx := make(map[string]int, len(fields))

	for _, field := range fields {
		fieldName := formatFieldName(field.Name)
		var b scalarBranchModel
		switch {
		case field.Type.IsArray():
			elementType := tf.formatFieldType(field.Type.AsArray().ValueType)
			b = scalarBranchModel{
				Field:         fieldName,
				ValueKinds:    []string{"Array"},
				DeserializeAs: fmt.Sprintf("List<%s>", elementType),
			}
		case field.Type.IsScalar():
			kinds, elemGet, writeAccess := elementBranchInfo(field.Type.AsScalar().ScalarKind)
			b = scalarBranchModel{
				Field:       fieldName,
				ValueKinds:  kinds,
				ElementGet:  elemGet,
				WriteAccess: writeAccess,
			}
		default:
			continue
		}

		key := tokenKey(b.ValueKinds)
		if idx, exists := kindIdx[key]; exists && key == "Number" {
			existing := out[idx]
			if isIntegerElementGet(b.ElementGet) && !isIntegerElementGet(existing.ElementGet) {
				existing.IntField = b.Field
				existing.IntElementGet = b.ElementGet
				existing.IntWrite = b.WriteAccess
			} else {
				existing.IntField = existing.Field
				existing.IntElementGet = existing.ElementGet
				existing.IntWrite = existing.WriteAccess
				existing.Field = b.Field
				existing.ElementGet = b.ElementGet
				existing.WriteAccess = b.WriteAccess
			}
			out[idx] = existing
			continue
		}
		if _, exists := kindIdx[key]; exists {
			continue
		}
		kindIdx[key] = len(out)
		out = append(out, b)
	}
	return out
}

// elementBranchInfo maps a scalar kind to the JsonValueKind(s) it
// reacts to and the JsonElement getter used to read the value.
func elementBranchInfo(kind ast.ScalarKind) (kinds []string, elementGet string, writeAccess string) {
	switch kind {
	case ast.KindString:
		return []string{"String"}, "GetString()", ""
	case ast.KindBool:
		return []string{"True", "False"}, "GetBoolean()", ".Value"
	case ast.KindInt8:
		return []string{"Number"}, "GetSByte()", ".Value"
	case ast.KindUint8, ast.KindBytes:
		return []string{"Number"}, "GetByte()", ".Value"
	case ast.KindInt16:
		return []string{"Number"}, "GetInt16()", ".Value"
	case ast.KindUint16:
		return []string{"Number"}, "GetUInt16()", ".Value"
	case ast.KindInt32:
		return []string{"Number"}, "GetInt32()", ".Value"
	case ast.KindUint32:
		return []string{"Number"}, "GetUInt32()", ".Value"
	case ast.KindInt64:
		return []string{"Number"}, "GetInt64()", ".Value"
	case ast.KindUint64:
		return []string{"Number"}, "GetUInt64()", ".Value"
	case ast.KindFloat32:
		return []string{"Number"}, "GetSingle()", ".Value"
	case ast.KindFloat64:
		return []string{"Number"}, "GetDouble()", ".Value"
	}
	return nil, "", ""
}

func isIntegerElementGet(expr string) bool {
	switch expr {
	case
		"GetSByte()",
		"GetByte()",
		"GetInt16()",
		"GetUInt16()",
		"GetInt32()",
		"GetUInt32()",
		"GetInt64()",
		"GetUInt64()":
		return true
	}
	return false
}

// scalarBranchInfo maps a scalar kind to the JsonTokenType(s) it
// reacts to, the reader expression used to read it, and whether
// `.Value` must be appended on writes (true for nullable value-types).
func scalarBranchInfo(kind ast.ScalarKind) (tokens []string, readExpr string, writeAccess string) {
	switch kind {
	case ast.KindString:
		return []string{"String"}, "reader.GetString()", ""
	case ast.KindBool:
		return []string{"True", "False"}, "reader.GetBoolean()", ".Value"
	case ast.KindInt8:
		return []string{"Number"}, "reader.GetSByte()", ".Value"
	case ast.KindUint8, ast.KindBytes:
		return []string{"Number"}, "reader.GetByte()", ".Value"
	case ast.KindInt16:
		return []string{"Number"}, "reader.GetInt16()", ".Value"
	case ast.KindUint16:
		return []string{"Number"}, "reader.GetUInt16()", ".Value"
	case ast.KindInt32:
		return []string{"Number"}, "reader.GetInt32()", ".Value"
	case ast.KindUint32:
		return []string{"Number"}, "reader.GetUInt32()", ".Value"
	case ast.KindInt64:
		return []string{"Number"}, "reader.GetInt64()", ".Value"
	case ast.KindUint64:
		return []string{"Number"}, "reader.GetUInt64()", ".Value"
	case ast.KindFloat32:
		return []string{"Number"}, "reader.GetSingle()", ".Value"
	case ast.KindFloat64:
		return []string{"Number"}, "reader.GetDouble()", ".Value"
	case ast.KindAny:
		return []string{"StartObject"}, "JsonSerializer.Deserialize<object>(ref reader, options)", ""
	}
	return []string{"None"}, "default!", ""
}

// sortRefCases left intentionally removed — sort.SliceStable is used
// directly inside renderRefs.
