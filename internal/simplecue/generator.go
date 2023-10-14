package simplecue

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/pkg/strconv"
	"github.com/grafana/cog/internal/ast"
)

const cogAnnotationName = "cog"
const cuetsyAnnotationName = "cuetsy"
const hintKindEnum = "enum"
const annotationKindFieldName = "kind"
const enumMembersAttr = "memberNames"

type Config struct {
	// Package name used to generate code into.
	Package string

	// For dataqueries, the schema doesn't define any top-level object
	ForceVariantEnvelope bool

	SchemaMetadata ast.SchemaMeta
}

type generator struct {
	schema *ast.Schema
}

func GenerateAST(val cue.Value, c Config) (*ast.Schema, error) {
	g := &generator{
		schema: &ast.Schema{
			Package:  c.Package,
			Metadata: c.SchemaMetadata,
		},
	}

	if c.ForceVariantEnvelope {
		if err := g.walkCueSchemaWithVariantEnvelope(val); err != nil {
			return nil, err
		}
	} else {
		if err := g.walkCueSchema(val); err != nil {
			return nil, err
		}
	}

	return g.schema, nil
}

func (g *generator) walkCueSchemaWithVariantEnvelope(v cue.Value) error {
	i, err := v.Fields(cue.Definitions(true), cue.All())
	if err != nil {
		return err
	}

	var rootObjectFields []ast.StructField
	for i.Next() {
		sel := i.Selector()
		name := selectorLabel(sel)

		if i.Selector().IsDefinition() {
			n, err := g.declareObject(name, i.Value())
			if err != nil {
				return err
			}

			g.schema.Objects = append(g.schema.Objects, n)
			continue
		}

		nodeType, err := g.declareNode(i.Value())
		if err != nil {
			return err
		}

		rootObjectFields = append(rootObjectFields, ast.NewStructField(name, nodeType))
	}

	g.schema.Objects = append(g.schema.Objects, ast.Object{
		Name:     string(g.schema.Metadata.Variant),
		Comments: commentsFromCueValue(v),
		Type:     ast.NewStruct(rootObjectFields...),
		SelfRef: ast.RefType{
			ReferredPkg:  g.schema.Package,
			ReferredType: string(g.schema.Metadata.Variant),
		},
	})

	return nil
}

func (g *generator) walkCueSchema(v cue.Value) error {
	i, err := v.Fields(cue.Definitions(true))
	if err != nil {
		return err
	}

	for i.Next() {
		sel := i.Selector()

		n, err := g.declareObject(selectorLabel(sel), i.Value())
		if err != nil {
			return err
		}

		g.schema.Objects = append(g.schema.Objects, n)
	}

	return nil
}

func (g *generator) declareObject(name string, v cue.Value) (ast.Object, error) {
	// Try to guess if `v` can be represented as an enum
	implicitEnum, err := isImplicitEnum(v)
	if err != nil {
		return ast.Object{}, err
	}
	if implicitEnum {
		return g.declareEnum(name, v)
	}

	nodeType, err := g.declareNode(v)
	if err != nil {
		return ast.Object{}, err
	}

	objectDef := ast.Object{
		Name:     name,
		Comments: commentsFromCueValue(v),
		Type:     nodeType,
		SelfRef: ast.RefType{
			ReferredPkg:  g.schema.Package,
			ReferredType: name,
		},
	}

	return objectDef, nil
}

func (g *generator) declareEnum(name string, v cue.Value) (ast.Object, error) {
	defVal, err := g.extractDefault(v)
	if err != nil {
		return ast.Object{}, err
	}

	enumType, err := g.declareAnonymousEnum(v, defVal)
	if err != nil {
		return ast.Object{}, err
	}

	return ast.Object{
		Name:     name,
		Comments: commentsFromCueValue(v),
		Type:     enumType,
		SelfRef: ast.RefType{
			ReferredPkg:  g.schema.Package,
			ReferredType: name,
		},
	}, nil
}

func (g *generator) extractEnumValues(v cue.Value) ([]ast.EnumValue, error) {
	_, dvals := v.Expr()
	a := v.Attribute(cogAnnotationName)
	// try to fall-back on the cuetsy annotation
	if a.Err() != nil {
		a = v.Attribute(cuetsyAnnotationName)
	}

	var attrMemberNameExist bool
	var evals []string
	if a.Err() == nil {
		val, found, err := a.Lookup(0, enumMembersAttr)
		if err == nil && found {
			attrMemberNameExist = true
			evals = strings.Split(val, "|")
			if len(evals) != len(dvals) {
				return nil, errorWithCueRef(v, "enums and %s attributes size doesn't match", enumMembersAttr)
			}
		}
	}

	// We only allowed String Enum to be generated without memberName attribute
	if v.IncompleteKind() != cue.StringKind && !attrMemberNameExist {
		return nil, errorWithCueRef(v, "numeric enums may only be generated from memberNames attribute")
	}

	subType := ast.String()
	if v.IncompleteKind() == cue.IntKind {
		subType = ast.NewScalar(ast.KindInt64)
	}

	var fields []ast.EnumValue
	for idx, dv := range dvals {
		var text string
		if attrMemberNameExist {
			text = evals[idx]
		} else {
			text, _ = dv.String()
		}

		if !dv.IsConcrete() {
			return nil, errorWithCueRef(v, "enums may only be generated from a disjunction of concrete strings or numbers")
		}

		val, err := cueConcreteToScalar(dv)
		if err != nil {
			return nil, err
		}

		fields = append(fields, ast.EnumValue{
			Type:  subType,
			Name:  text,
			Value: val,
		})
	}

	return fields, nil
}

func (g *generator) structFields(v cue.Value) ([]ast.StructField, error) {
	// This check might be too restrictive
	if v.IncompleteKind() != cue.StructKind {
		return nil, errorWithCueRef(v, "top-level type definitions may only be generated from structs")
	}

	var fields []ast.StructField

	// explore struct fields
	for i, _ := v.Fields(cue.Optional(true), cue.Definitions(true)); i.Next(); {
		fieldLabel := selectorLabel(i.Selector())

		node, err := g.declareNode(i.Value())
		if err != nil {
			return nil, err
		}

		fields = append(fields, ast.StructField{
			Name:     fieldLabel,
			Comments: commentsFromCueValue(i.Value()),
			Required: !i.IsOptional(),
			Type:     node,
		})
	}

	return fields, nil
}

// FIXME: this is probably very brittle and not always correct :|
func (g *generator) referencePackage(source cueast.Node) (string, error) {
	switch source.(type) { //nolint: gocritic
	case *cueast.SelectorExpr:
		selector := source.(*cueast.SelectorExpr)

		x := selector.X.(*cueast.Ident)

		return x.Name, nil
	case *cueast.Field:
		field := source.(*cueast.Field)

		if _, ok := field.Value.(*cueast.SelectorExpr); ok {
			return g.referencePackage(field.Value)
		}

		ident := field.Value.(*cueast.Ident)

		if ident.Scope == nil {
			return g.schema.Package, nil
		}

		if _, ok := ident.Scope.(*cueast.File); !ok {
			return g.schema.Package, nil
		}

		scope := ident.Scope.(*cueast.File)
		if len(scope.Decls) == 0 {
			return g.schema.Package, nil
		}

		if _, ok := scope.Decls[0].(*cueast.Package); !ok {
			return g.schema.Package, nil
		}

		referredTypePkg := scope.Decls[0].(*cueast.Package).Name

		return referredTypePkg.Name, nil
	case *cueast.Ident:
		ident := source.(*cueast.Ident)
		if ident.Scope == nil {
			return g.schema.Package, nil
		}

		if _, ok := ident.Scope.(*cueast.File); !ok {
			return g.schema.Package, nil
		}

		scope := ident.Scope.(*cueast.File)
		if len(scope.Decls) == 0 {
			return g.schema.Package, nil
		}

		referredTypePkg := ident.Scope.(*cueast.File).Decls[0].(*cueast.Package).Name

		return referredTypePkg.Name, nil
	case *cueast.Ellipsis: // TODO: this makes no sense
		return g.schema.Package, nil
	default:
		return "", fmt.Errorf("can't extract reference package")
	}
}

func (g *generator) declareNode(v cue.Value) (ast.Type, error) {
	// This node is referring to another definition
	_, path := v.ReferencePath()
	if path.String() != "" {
		selectors := path.Selectors()

		refPkg, err := g.referencePackage(v.Source())
		if err != nil {
			return ast.Type{}, errorWithCueRef(v, err.Error())
		}

		return ast.NewRef(
			refPkg,
			selectorLabel(selectors[len(selectors)-1]),
		), nil
	}

	defVal, err := g.extractDefault(v)
	if err != nil {
		return ast.Type{}, err
	}

	disjunctions := appendSplit(nil, cue.OrOp, v)
	// Possible cases:
	// 1. "value" | "other_value" | "concrete_value" → we want to parse that as an "enum" type
	// 2. SomeType | SomeOtherType | string → we want to parse that as a disjunction
	if len(disjunctions) != 1 {
		// Try to guess if `v` can be represented as an enum (includes checking for a type hint) (1)
		implicitEnum, err := isImplicitEnum(v)
		if err != nil {
			return ast.Type{}, err
		}
		if implicitEnum {
			return g.declareAnonymousEnum(v, defVal)
		}

		// We must be looking at a disjunction then (2)
		branches := make([]ast.Type, 0, len(disjunctions))
		for _, subTypeValue := range disjunctions {
			subType, err := g.declareNode(subTypeValue)
			if err != nil {
				return ast.Type{}, err
			}

			branches = append(branches, subType)
		}

		return ast.NewDisjunction(branches, ast.Default(defVal)), nil
	}

	switch v.IncompleteKind() {
	case cue.TopKind:
		return ast.Any(), nil
	case cue.NullKind:
		return ast.Null(), nil
	case cue.BoolKind:
		return ast.Bool(ast.Default(defVal)), nil
	case cue.BytesKind:
		return ast.Bytes(ast.Default(defVal)), nil
	case cue.StringKind:
		return g.declareString(v, defVal)
	case cue.FloatKind, cue.NumberKind, cue.IntKind:
		return g.declareNumber(v, defVal)
	case cue.ListKind:
		return g.declareList(v)
	case cue.StructKind:
		// in cue: {...}, {[string]: type}, or inline struct
		if op, _ := v.Expr(); op == cue.NoOp {
			t := v.LookupPath(cue.MakePath(cue.AnyString))
			if t.Exists() && t.IncompleteKind() != cue.TopKind {
				typeDef, err := g.declareNode(t)
				if err != nil {
					return ast.Type{}, err
				}

				return ast.NewMap(ast.String(), typeDef), nil
			}
		}

		fields, err := g.structFields(v)
		if err != nil {
			return ast.Type{}, err
		}

		// {...}
		if len(fields) == 0 {
			return ast.Any(), nil
		}

		return ast.NewStruct(fields...), nil
	default:
		return ast.Type{}, errorWithCueRef(v, "unexpected node with kind '%s'", v.IncompleteKind().String())
	}
}

func (g *generator) declareAnonymousEnum(v cue.Value, defValue any) (ast.Type, error) {
	allowed := cue.StringKind | cue.IntKind
	ik := v.IncompleteKind()
	if ik&allowed != ik {
		return ast.Type{}, errorWithCueRef(v, "enums may only be generated from concrete strings, or ints")
	}

	values, err := g.extractEnumValues(v)
	if err != nil {
		return ast.Type{}, err
	}

	return ast.NewEnum(values, ast.Default(defValue)), nil
}

func (g *generator) declareString(v cue.Value, defVal any) (ast.Type, error) {
	typeDef := ast.String(ast.Default(defVal))

	// v.IsConcrete() being true means we're looking at a constant/known value
	if v.IsConcrete() {
		val, err := cueConcreteToScalar(v)
		if err != nil {
			return typeDef, err
		}

		typeDef.Scalar.Value = val
	}

	// Extract constraints
	constraints, err := g.declareStringConstraints(v)
	if err != nil {
		return typeDef, err
	}

	typeDef.Scalar.Constraints = constraints

	return typeDef, nil
}

func (g *generator) extractDefault(v cue.Value) (any, error) {
	defaultVal, ok := v.Default()
	if !ok {
		//nolint: nilnil
		return nil, nil
	}

	def, err := cueConcreteToScalar(defaultVal)
	if err != nil {
		return nil, err
	}

	return def, nil
}

func (g *generator) declareStringConstraints(v cue.Value) ([]ast.TypeConstraint, error) {
	typeAndConstraints := appendSplit(nil, cue.AndOp, v)

	// nothing to do
	if len(typeAndConstraints) == 1 {
		return nil, nil
	}

	// the constraint allows cue to infer a concrete value
	// ex: #SomeEnumType & "some value from the enum"
	if v.IsConcrete() {
		stringVal, err := v.String()
		if err != nil {
			return nil, errorWithCueRef(v, "could not convert concrete value to string")
		}

		return []ast.TypeConstraint{
			{
				Op:   ast.EqualOp,
				Args: []any{stringVal},
			},
		}, nil
	}

	constraints := make([]ast.TypeConstraint, 0, len(typeAndConstraints))

	for _, andExpr := range typeAndConstraints {
		op, args := andExpr.Expr()

		// TODO: support more OPs?
		if op != cue.CallOp {
			continue
		}

		// TODO: support more constraints?
		switch fmt.Sprint(args[0]) {
		case "strings.MinRunes":
			scalar, err := cueConcreteToScalar(args[1])
			if err != nil {
				return nil, err
			}

			constraints = append(constraints, ast.TypeConstraint{
				Op:   ast.MinLengthOp,
				Args: []any{scalar},
			})

		case "strings.MaxRunes":
			scalar, err := cueConcreteToScalar(args[1])
			if err != nil {
				return nil, err
			}

			constraints = append(constraints, ast.TypeConstraint{
				Op:   ast.MaxLengthOp,
				Args: []any{scalar},
			})
		}
	}

	return constraints, nil
}

func (g *generator) declareNumber(v cue.Value, defVal any) (ast.Type, error) {
	numberTypeWithConstraintsAsString, err := format.Node(v.Syntax())
	if err != nil {
		return ast.Type{}, err
	}
	parts := strings.Split(string(numberTypeWithConstraintsAsString), " ")
	if len(parts) == 0 {
		return ast.Type{}, errorWithCueRef(v, "something went very wrong while formatting a number expression into a string")
	}

	// dirty way of preserving the actual type from cue
	// note that CUE has predefined "types": https://cuelang.org/docs/tutorials/tour/types/bounddef/
	var numberType ast.ScalarKind
	for _, numberTypeCandidate := range parts {
		switch ast.ScalarKind(numberTypeCandidate) {
		case ast.KindFloat32, ast.KindFloat64:
			numberType = ast.ScalarKind(numberTypeCandidate)
		case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64:
			numberType = ast.ScalarKind(numberTypeCandidate)
		case ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
			numberType = ast.ScalarKind(numberTypeCandidate)
		case "uint":
			numberType = ast.KindUint64
		case "int":
			numberType = ast.KindInt64
		case "float":
			numberType = ast.KindFloat64
		case "number":
			numberType = ast.KindFloat64
		}
	}

	if numberType == "" {
		return ast.Type{}, errorWithCueRef(v, "could not infer number type from expression '%s'", numberTypeWithConstraintsAsString)
	}

	typeDef := ast.NewScalar(numberType, ast.Default(defVal))

	// v.IsConcrete() being true means we're looking at a constant/known value
	if v.IsConcrete() {
		val, err := cueConcreteToScalar(v)
		if err != nil {
			return typeDef, err
		}

		typeDef.Scalar.Value = val
	}

	// If the default (all lists have a default, usually self, ugh) differs from the
	// input list, peel it off. Otherwise our AnyIndex lookup may end up getting
	// sent on the wrong path.

	// extract constraints
	constraints, err := g.declareNumberConstraints(v)
	if err != nil {
		return ast.Type{}, err
	}

	typeDef.Scalar.Constraints = constraints

	return typeDef, nil
}

// having written this makes my soul hurt.
func (g *generator) declareNumberConstraints(v cue.Value) ([]ast.TypeConstraint, error) {
	// if the number has a default value, strip it from `v` before trying to extract constraints.
	_, hasDefault := v.Default()
	if hasDefault {
		_, dvals := v.Expr()
		v = dvals[0]
	}

	numberTypeWithConstraintsAsString, err := format.Node(v.Syntax())
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(numberTypeWithConstraintsAsString), " & ")

	extractOperatorAndArg := func(input string, numberKind cue.Kind) (ast.Op, any, error) {
		argStartIndex := -1

		var op ast.Op

		if input[0] == '>' {
			op = ast.GreaterThanOp
			argStartIndex = 1

			if input[1] == '=' {
				op = ast.GreaterThanEqualOp
				argStartIndex = 2
			}
		}
		if input[0] == '<' {
			op = ast.LessThanOp
			argStartIndex = 1

			if input[1] == '=' {
				op = ast.LessThanEqualOp
				argStartIndex = 2
			}
		}
		if input[0] == '=' && input[1] == '=' {
			op = ast.EqualOp
			argStartIndex = 2
		}
		if input[0] == '!' && input[1] == '=' {
			op = ast.NotEqualOp
			argStartIndex = 2
		}

		if op == "" {
			return op, nil, fmt.Errorf("could not infer operator from '%s'", input)
		}

		if numberKind.IsAnyOf(cue.FloatKind) {
			arg, err := strconv.ParseFloat(input[argStartIndex:], 64)
			return op, arg, err
		}

		arg, err := strconv.ParseInt(input[argStartIndex:], 10, 64)
		return op, arg, err
	}

	var constraints []ast.TypeConstraint
	for _, part := range parts {
		if part[0] != '<' && part[0] != '>' {
			continue
		}

		op, arg, err := extractOperatorAndArg(part, v.IncompleteKind())
		if err != nil {
			return nil, errorWithCueRef(v, err.Error())
		}

		constraints = append(constraints, ast.TypeConstraint{
			Op:   op,
			Args: []any{arg},
		})
	}

	return constraints, nil
}

func (g *generator) declareList(v cue.Value) (ast.Type, error) {
	i, err := v.List()
	if err != nil {
		return ast.Type{}, err
	}

	typeDef := ast.NewArray(ast.Any())

	// works only for a closed/concrete list
	if v.IsConcrete() {
		// fixme: this is wrong
		for i.Next() {
			node, err := g.declareNode(i.Value())
			if err != nil {
				return ast.Type{}, err
			}

			typeDef.Array.ValueType = node
		}

		return typeDef, nil
	}

	// open list

	// If the default (all lists have a default, usually self, ugh) differs from the
	// input list, peel it off. Otherwise our AnyIndex lookup may end up getting
	// sent on the wrong path.
	defv, _ := v.Default()
	if !defv.Equals(v) {
		_, dvals := v.Expr()
		v = dvals[0]
	}

	e := v.LookupPath(cue.MakePath(cue.AnyIndex))
	if !e.Exists() {
		// unreachable?
		return ast.Type{}, errorWithCueRef(v, "open list must have a type")
	}

	expr, err := g.declareNode(e)
	if err != nil {
		return ast.Type{}, err
	}

	typeDef.Array.ValueType = expr

	return typeDef, nil
}
