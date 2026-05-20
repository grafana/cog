package zod

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	config Config
}

func (jenny RawTypes) JennyName() string {
	return "ZodRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	// flat_layout puts every schema at <prefix>/<filename>, so multiple
	// schemas would collide. Fail rather than silently overwrite.
	if jenny.config.FlatLayout && len(context.Schemas) > 1 {
		pkgs := make([]string, 0, len(context.Schemas))
		for _, schema := range context.Schemas {
			pkgs = append(pkgs, schema.Package)
		}
		sort.Strings(pkgs)
		return nil, fmt.Errorf(
			"zod jenny: flat_layout cannot be combined with multiple input schemas — files would collide at the same path (got %d schemas: %s)",
			len(context.Schemas), strings.Join(pkgs, ", "),
		)
	}

	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := jenny.config.outputFilePath(formatPackageName(schema.Package))

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(context languages.Context, schema *ast.Schema) ([]byte, error) {
	imports := newImportMap()

	// flat_layout drops per-package dirs, breaking the "../<alias>" import
	// pattern. Collect cross-package refs and fail later, rather than emit
	// an unresolvable import.
	var crossPkgRefs []string
	seenCrossPkgRefs := make(map[string]struct{})

	pkgMapper := func(pkg string) string {
		if pkg == "" || pkg == schema.Package {
			return ""
		}
		alias := formatPackageName(pkg)
		if jenny.config.FlatLayout {
			if _, seen := seenCrossPkgRefs[pkg]; !seen {
				seenCrossPkgRefs[pkg] = struct{}{}
				crossPkgRefs = append(crossPkgRefs, pkg)
			}
		}
		return imports.Add(alias, "../"+alias)
	}

	em := &emitter{
		context:       context,
		config:        jenny.config,
		pkg:           schema.Package,
		packageMapper: pkgMapper,
	}

	var body strings.Builder
	schema.Objects.Iterate(func(_ string, object ast.Object) {
		body.WriteString(em.formatObject(object))
		body.WriteString("\n")
	})

	if len(crossPkgRefs) > 0 {
		sort.Strings(crossPkgRefs)
		return nil, fmt.Errorf(
			"zod jenny: flat_layout is incompatible with cross-package references; schema %q references %s, and import paths cannot be resolved without per-package directories — disable flat_layout or eliminate the cross-package reference",
			schema.Package, strings.Join(crossPkgRefs, ", "),
		)
	}

	header := "// @ts-nocheck\n"
	header += "import { z } from 'zod';\n"
	if importStatements := imports.String(); importStatements != "" {
		header += importStatements
	}
	header += "\n"

	return []byte(header + body.String()), nil
}

// emitter walks the AST for one schema and produces Zod source as strings.
type emitter struct {
	context       languages.Context
	config        Config
	pkg           string
	packageMapper func(pkg string) string
}

// stringLiteralValue returns (value, true) when t is emitted as
// `z.literal("X")` — covers inline string scalars (`kind: "Panel"`) and
// refs to string constants (`mode: RepeatMode` with `RepeatMode = "variable"`).
func stringLiteralValue(t ast.Type) (any, bool) {
	switch t.Kind {
	case ast.KindScalar:
		s := t.AsScalar()
		if s.IsConcrete() && s.ScalarKind == ast.KindString {
			return s.Value, true
		}
	case ast.KindConstantRef:
		v := t.AsConstantRef().ReferenceValue
		if _, ok := v.(string); ok {
			return v, true
		}
	}
	return nil, false
}

func (e *emitter) formatObject(def ast.Object) string {
	var buf strings.Builder

	objectName := tools.CleanupNames(def.Name)

	// Constants emit as a plain TS const + typeof alias so consumers can use
	// them as both value and type. No schema means no .describe(), so doc
	// comments stay as plain JS comments.
	if def.Type.IsConcreteScalar() {
		for _, comment := range def.Comments {
			buf.WriteString("// ")
			buf.WriteString(comment)
			buf.WriteString("\n")
		}
		val := def.Type.AsScalar().Value
		buf.WriteString(fmt.Sprintf("export const %s = %s;\n", objectName, formatLiteral(val)))
		buf.WriteString(fmt.Sprintf("export type %s = typeof %s;\n", objectName, objectName))
		return buf.String()
	}

	expr := e.formatType(def.Type)
	if desc := joinDescription(def.Comments); desc != "" {
		expr += fmt.Sprintf(".describe(%s)", formatLiteral(desc))
	}

	buf.WriteString(fmt.Sprintf("export const %sSchema = %s;\n", objectName, expr))
	return buf.String()
}

func (e *emitter) formatType(t ast.Type) string {
	expr := e.formatBareType(t)
	if t.Nullable {
		expr += ".nullable()"
	}
	return expr
}

func (e *emitter) formatBareType(t ast.Type) string {
	switch t.Kind {
	case ast.KindStruct:
		return e.formatStruct(t)
	case ast.KindScalar:
		return e.formatScalar(t.AsScalar())
	case ast.KindRef:
		return e.formatRef(t.AsRef())
	case ast.KindEnum:
		return e.formatEnum(t.AsEnum())
	case ast.KindArray:
		return e.formatArray(t.AsArray())
	case ast.KindMap:
		return e.formatMap(t.AsMap())
	case ast.KindDisjunction:
		return e.formatDisjunction(t.AsDisjunction())
	case ast.KindIntersection:
		return e.formatIntersection(t.AsIntersection())
	case ast.KindComposableSlot:
		return "z.unknown()"
	case ast.KindConstantRef:
		return e.formatConstantRef(t.AsConstantRef())
	}
	return fmt.Sprintf("z.unknown() /* unhandled kind: %s */", t.Kind)
}

func (e *emitter) formatStruct(t ast.Type) string {
	var buf strings.Builder

	buf.WriteString("z.object({\n")
	for _, field := range t.AsStruct().Fields {
		fieldExpr := e.formatType(field.Type)

		// OptionalKindLiterals implies VerboseDefaults so the file uses one
		// form for every defaulted field.
		verboseDefaults := e.config.VerboseDefaults || e.config.OptionalKindLiterals

		litVal, isStringLit := stringLiteralValue(field.Type)
		switch {
		case e.config.OptionalKindLiterals &&
			field.Required &&
			field.Type.Default == nil &&
			isStringLit:
			fieldExpr += fmt.Sprintf(".optional().default(%s)", formatLiteral(litVal))

		case field.Type.Default != nil && verboseDefaults:
			fieldExpr += fmt.Sprintf(".optional().default(%s)", formatLiteral(field.Type.Default))

		case field.Type.Default != nil:
			fieldExpr += fmt.Sprintf(".default(%s)", formatLiteral(field.Type.Default))
			if !field.Required {
				fieldExpr += ".optional()"
			}

		case !field.Required:
			fieldExpr += ".optional()"
		}

		if desc := joinDescription(field.Comments); desc != "" {
			fieldExpr += fmt.Sprintf(".describe(%s)", formatLiteral(desc))
		}

		buf.WriteString(fmt.Sprintf("\t%s: %s,\n", quoteFieldName(field.Name), fieldExpr))
	}
	buf.WriteString("})")

	// CUE open structs (`...`, `[string]: _`) materialize upstream as
	// KindMap / KindAny, so anything reaching here is closed.
	switch e.config.ObjectMode {
	case "strip", "":
	case "loose":
		buf.WriteString(".loose()")
	case "strict":
		buf.WriteString(".strict()")
	default:
		// Unknown mode → strict, so a typo surfaces at parse time.
		buf.WriteString(".strict()")
	}

	return buf.String()
}

func (e *emitter) formatScalar(s ast.ScalarType) string {
	if s.IsConcrete() {
		return fmt.Sprintf("z.literal(%s)", formatLiteral(s.Value))
	}

	for _, c := range s.Constraints {
		if c.Op == ast.EqualOp {
			return fmt.Sprintf("z.literal(%s)", formatLiteral(c.Args[0]))
		}
	}

	switch s.ScalarKind {
	case ast.KindNull:
		return "z.null()"
	case ast.KindAny:
		if e.config.PreferUnknownForAny {
			return "z.unknown()"
		}
		return "z.any()"
	case ast.KindBool:
		return "z.boolean()"
	case ast.KindBytes, ast.KindString:
		return applyStringConstraints("z.string()", s.Constraints)
	case ast.KindFloat32, ast.KindFloat64:
		return applyNumberConstraints("z.number()", s.Constraints)
	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64,
		ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		return applyNumberConstraints("z.number().int()", s.Constraints)
	}
	return fmt.Sprintf("z.unknown() /* unhandled scalar: %s */", s.ScalarKind)
}

func applyStringConstraints(base string, cs []ast.TypeConstraint) string {
	for _, c := range cs {
		switch c.Op {
		case ast.MinLengthOp:
			base += fmt.Sprintf(".min(%s)", formatLiteral(c.Args[0]))
		case ast.MaxLengthOp:
			base += fmt.Sprintf(".max(%s)", formatLiteral(c.Args[0]))
		case ast.RegexMatchOp:
			base += fmt.Sprintf(".regex(new RegExp(%s))", formatLiteral(c.Args[0]))
		case ast.NotRegexMatchOp:
			base += fmt.Sprintf(".refine((value) => !(new RegExp(%s)).test(value), { message: %q })", formatLiteral(c.Args[0]), "String must not match forbidden pattern")
		}
	}
	return base
}

func applyNumberConstraints(base string, cs []ast.TypeConstraint) string {
	for _, c := range cs {
		switch c.Op {
		case ast.LessThanOp:
			base += fmt.Sprintf(".lt(%s)", formatLiteral(c.Args[0]))
		case ast.LessThanEqualOp:
			base += fmt.Sprintf(".lte(%s)", formatLiteral(c.Args[0]))
		case ast.GreaterThanOp:
			base += fmt.Sprintf(".gt(%s)", formatLiteral(c.Args[0]))
		case ast.GreaterThanEqualOp:
			base += fmt.Sprintf(".gte(%s)", formatLiteral(c.Args[0]))
		case ast.MultipleOfOp:
			base += fmt.Sprintf(".multipleOf(%s)", formatLiteral(c.Args[0]))
		case ast.NotEqualOp:
			base += fmt.Sprintf(".refine((value) => value !== %s, { message: %q })", formatLiteral(c.Args[0]), "Value must not equal forbidden constant")
		case ast.EqualOp:
			base += fmt.Sprintf(".refine((value) => value === %s, { message: %q })", formatLiteral(c.Args[0]), "Value must equal required constant")
		}
	}
	return base
}

func (e *emitter) formatRef(ref ast.RefType) string {
	name := tools.CleanupNames(ref.ReferredType) + "Schema"
	if alias := e.packageMapper(ref.ReferredPkg); alias != "" {
		name = alias + "." + name
	}
	// Wrap every cross-reference in z.lazy() to sidestep declaration order.
	// TODO: topo-sort objects so only real cycles need z.lazy.
	return fmt.Sprintf("z.lazy(() => %s)", name)
}

func (e *emitter) formatEnum(en ast.EnumType) string {
	allString := true
	for _, v := range en.Values {
		if v.Type.Scalar == nil || v.Type.Scalar.ScalarKind != ast.KindString {
			allString = false
			break
		}
	}

	if allString {
		parts := tools.Map(en.Values, func(v ast.EnumValue) string {
			return formatLiteral(v.Value)
		})
		return fmt.Sprintf("z.enum([%s])", strings.Join(parts, ", "))
	}

	parts := tools.Map(en.Values, func(v ast.EnumValue) string {
		return fmt.Sprintf("z.literal(%s)", formatLiteral(v.Value))
	})
	return fmt.Sprintf("z.union([%s])", strings.Join(parts, ", "))
}

func (e *emitter) formatArray(arr ast.ArrayType) string {
	base := fmt.Sprintf("z.array(%s)", e.formatType(arr.ValueType))
	for _, c := range arr.Constraints {
		switch c.Op {
		case ast.MinItemsOp:
			base += fmt.Sprintf(".min(%s)", formatLiteral(c.Args[0]))
		case ast.MaxItemsOp:
			base += fmt.Sprintf(".max(%s)", formatLiteral(c.Args[0]))
		case ast.UniqueItemsOp:
			base += `.refine((items) => new Set(items.map((item) => JSON.stringify(item))).size === items.length, { message: "Array items must be unique" })`
		}
	}
	return base
}

func (e *emitter) formatMap(m ast.MapType) string {
	return fmt.Sprintf("z.record(%s, %s)", e.formatType(m.IndexType), e.formatType(m.ValueType))
}

func (e *emitter) formatDisjunction(d ast.DisjunctionType) string {
	branches := tools.Map(d.Branches, func(b ast.Type) string {
		return e.formatType(b)
	})

	// z.discriminatedUnion needs every branch to be a z.object (possibly
	// wrapped in z.lazy). HasOnlyRefs alone isn't enough — refs to enums or
	// nested unions would make Zod throw at module load. Fall back to
	// z.union when we can't prove every branch is a struct ref.
	if d.Discriminator != "" && d.Branches.HasOnlyRefs() && e.allBranchesResolveToStruct(d.Branches) {
		return fmt.Sprintf("z.discriminatedUnion(%q, [%s])", d.Discriminator, strings.Join(branches, ", "))
	}

	return fmt.Sprintf("z.union([%s])", strings.Join(branches, ", "))
}

func (e *emitter) allBranchesResolveToStruct(branches ast.Types) bool {
	for _, b := range branches {
		if !e.context.ResolveToStruct(b) {
			return false
		}
	}
	return true
}

func (e *emitter) formatIntersection(i ast.IntersectionType) string {
	branches := tools.Map(i.Branches, func(b ast.Type) string {
		return e.formatType(b)
	})
	switch len(branches) {
	case 0:
		return "z.unknown()"
	case 1:
		return branches[0]
	}

	expr := branches[0]
	for _, b := range branches[1:] {
		expr = fmt.Sprintf("z.intersection(%s, %s)", expr, b)
	}
	return expr
}

func (e *emitter) formatConstantRef(c ast.ConstantReferenceType) string {
	return fmt.Sprintf("z.literal(%s)", formatLiteral(c.ReferenceValue))
}

func newImportMap() *common.DirectImportMap {
	return common.NewDirectImportMap(
		common.WithFormatter(func(im common.DirectImportMap) string {
			if im.Imports.Len() == 0 {
				return ""
			}
			var b strings.Builder
			im.Imports.Iterate(func(alias string, importPath string) {
				b.WriteString(fmt.Sprintf("import * as %s from '%s';\n", alias, importPath))
			})
			return b.String()
		}),
	)
}
