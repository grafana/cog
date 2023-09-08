package simplecue

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
	"github.com/grafana/cog/internal/ast"
)

const annotationName = "cog"
const hintKindEnum = "enum"
const annotationKindFieldName = "kind"
const enumMembersAttr = "memberNames"

type Config struct {
	// Package name used to generate code into.
	Package string
}

type newGenerator struct {
	file *ast.File
}

func GenerateAST(val cue.Value, c Config) (*ast.File, error) {
	g := &newGenerator{
		file: &ast.File{
			Package: c.Package,
		},
	}

	i, err := val.Fields(cue.Definitions(true))
	if err != nil {
		return nil, err
	}
	for i.Next() {
		sel := i.Selector()

		n, err := g.declareTopLevelType(selectorLabel(sel), i.Value())
		if err != nil {
			return nil, err
		}

		g.file.Definitions = append(g.file.Definitions, n)
	}

	return g.file, nil
}

// Do we really need to distinguish top-level types with others?
func (g *newGenerator) declareTopLevelType(name string, v cue.Value) (ast.Object, error) {
	typeHint, err := getTypeHint(v)
	if err != nil {
		return ast.Object{}, err
	}

	// Hinted as an enum
	if typeHint == hintKindEnum {
		return g.declareEnum(name, v)
	}

	ik := v.IncompleteKind()

	// Is it a string disjunction that we can turn into an enum?
	disjunctions := appendSplit(nil, cue.OrOp, v)
	if len(disjunctions) != 1 && ik&cue.StringKind == ik {
		return g.declareEnum(name, v)
	}

	switch v.IncompleteKind() {
	case cue.StringKind:
		return g.declareTopLevelString(name, v)
	case cue.StructKind:
		return g.declareTopLevelStruct(name, v)
	default:
		return ast.Object{}, errorWithCueRef(v, "unexpected top-level kind '%s'", v.IncompleteKind().String())
	}
}

func (g *newGenerator) declareTopLevelString(name string, v cue.Value) (ast.Object, error) {
	ik := v.IncompleteKind()
	if ik&cue.StringKind != ik {
		return ast.Object{}, errorWithCueRef(v, "top-level strings may only be generated from concrete strings")
	}

	strType, err := g.declareString(v)
	if err != nil {
		return ast.Object{}, err
	}

	return ast.Object{
		Name:     name,
		Comments: commentsFromCueValue(v),
		Type:     strType,
	}, nil
}

func (g *newGenerator) declareEnum(name string, v cue.Value) (ast.Object, error) {
	// Restrict the expression of enums to ints or strings.
	allowed := cue.StringKind | cue.IntKind
	ik := v.IncompleteKind()
	if ik&allowed != ik {
		return ast.Object{}, errorWithCueRef(v, "enums may only be generated from concrete strings, or ints")
	}

	values, err := g.extractEnumValues(v)
	if err != nil {
		return ast.Object{}, err
	}

	return ast.Object{
		Name:     name,
		Comments: commentsFromCueValue(v),
		Type: ast.EnumType{
			Values: values,
		},
	}, nil
}

func (g *newGenerator) extractEnumValues(v cue.Value) ([]ast.EnumValue, error) {
	_, dvals := v.Expr()
	a := v.Attribute(annotationName)

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

	subType := ast.ScalarType{ScalarKind: ast.KindString}
	if v.IncompleteKind() == cue.IntKind {
		subType = ast.ScalarType{ScalarKind: ast.KindInt64}
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
			Type: subType,

			// Simple mapping of all enum values (which we are assuming are in
			// lowerCamelCase) to corresponding CamelCase
			Name:  text,
			Value: val,
		})
	}

	return fields, nil
}

func (g *newGenerator) declareTopLevelStruct(name string, v cue.Value) (ast.Object, error) {
	// This check might be too restrictive
	if v.IncompleteKind() != cue.StructKind {
		return ast.Object{}, errorWithCueRef(v, "top-level type definitions may only be generated from structs")
	}

	nodeType, err := g.declareNode(v)
	if err != nil {
		return ast.Object{}, err
	}

	typeDef := ast.Object{
		Name:     name,
		Comments: commentsFromCueValue(v),
		Type:     nodeType,
	}

	return typeDef, nil
}

func (g *newGenerator) structFields(v cue.Value) ([]ast.StructField, error) {
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

		// Extract the default value if it's there
		defVal, err := g.extractDefault(i.Value())
		if err != nil {
			return nil, err
		}

		fields = append(fields, ast.StructField{
			Name:     fieldLabel,
			Comments: commentsFromCueValue(i.Value()),
			Required: !i.IsOptional(),
			Type:     node,
			Default:  defVal,
		})
	}

	return fields, nil
}

func (g *newGenerator) declareNode(v cue.Value) (ast.Type, error) {
	// This node is referring to another definition
	_, path := v.ReferencePath()
	if path.String() != "" {
		selectors := path.Selectors()

		return ast.RefType{
			ReferredType: selectorLabel(selectors[len(selectors)-1]),
		}, nil
	}

	disjunctions := appendSplit(nil, cue.OrOp, v)
	if len(disjunctions) != 1 {
		allowedKindsForAnonymousEnum := cue.StringKind | cue.IntKind
		ik := v.IncompleteKind()
		if ik&allowedKindsForAnonymousEnum == ik {
			return g.declareAnonymousEnum(v)
		}

		branches := make([]ast.Type, 0, len(disjunctions))
		for _, subTypeValue := range disjunctions {
			subType, err := g.declareNode(subTypeValue)
			if err != nil {
				return nil, err
			}

			branches = append(branches, subType)
		}

		return ast.DisjunctionType{
			Branches: branches,
		}, nil
	}

	switch v.IncompleteKind() {
	case cue.TopKind:
		return ast.ScalarType{ScalarKind: ast.KindAny}, nil
	case cue.NullKind:
		return ast.ScalarType{ScalarKind: ast.KindNull}, nil
	case cue.BoolKind:
		return ast.ScalarType{
			ScalarKind: ast.KindBool,
		}, nil
	case cue.BytesKind:
		return ast.ScalarType{ScalarKind: ast.KindBytes}, nil
	case cue.StringKind:
		return g.declareString(v)
	case cue.FloatKind, cue.NumberKind, cue.IntKind:
		return g.declareNumber(v)
	case cue.ListKind:
		return g.declareList(v)
	case cue.StructKind:
		op, opArgs := v.Expr()

		// in cue: {...}, {[string]: type}, or inline struct

		if op == cue.NoOp {
			// looking for {[string]: type}
			// (don't know how to do this properly)
			//nolint: staticcheck
			t, ok := opArgs[0].Elem()
			if ok && t.IncompleteKind() != cue.TopKind {
				typeDef, err := g.declareNode(t)
				if err != nil {
					return nil, err
				}

				return ast.MapType{
					IndexType: ast.ScalarType{ScalarKind: ast.KindString},
					ValueType: typeDef,
				}, nil
			}
		}

		fields, err := g.structFields(v)
		if err != nil {
			return nil, err
		}

		// {...}
		if len(fields) == 0 {
			return ast.ScalarType{ScalarKind: ast.KindAny}, nil
		}

		return ast.StructType{Fields: fields}, nil
	default:
		return nil, errorWithCueRef(v, "unexpected node with kind '%s'", v.IncompleteKind().String())
	}
}

func (g *newGenerator) declareAnonymousEnum(v cue.Value) (ast.Type, error) {
	allowed := cue.StringKind | cue.IntKind
	ik := v.IncompleteKind()
	if ik&allowed != ik {
		return nil, errorWithCueRef(v, "enums may only be generated from concrete strings, or ints")
	}

	values, err := g.extractEnumValues(v)
	if err != nil {
		return nil, err
	}

	return ast.EnumType{
		Values: values,
	}, nil
}

func (g *newGenerator) declareString(v cue.Value) (ast.ScalarType, error) {
	typeDef := ast.ScalarType{
		ScalarKind: ast.KindString,
	}

	// v.IsConcrete() being true means we're looking at a constant/known value
	if v.IsConcrete() {
		val, err := cueConcreteToScalar(v)
		if err != nil {
			return typeDef, err
		}

		typeDef.Value = val
	}

	// Extract constraints
	constraints, err := g.declareStringConstraints(v)
	if err != nil {
		return typeDef, err
	}

	typeDef.Constraints = constraints

	return typeDef, nil
}

func (g *newGenerator) extractDefault(v cue.Value) (any, error) {
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

func (g *newGenerator) declareStringConstraints(v cue.Value) ([]ast.TypeConstraint, error) {
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
				Op:   "==",
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
				Op:   "minLength",
				Args: []any{scalar},
			})

		case "strings.MaxRunes":
			scalar, err := cueConcreteToScalar(args[1])
			if err != nil {
				return nil, err
			}

			constraints = append(constraints, ast.TypeConstraint{
				Op:   "maxLength",
				Args: []any{scalar},
			})
		}
	}

	return constraints, nil
}

func (g *newGenerator) declareNumber(v cue.Value) (ast.ScalarType, error) {
	numberTypeWithConstraintsAsString, err := format.Node(v.Syntax())
	if err != nil {
		return ast.ScalarType{}, err
	}
	parts := strings.Split(string(numberTypeWithConstraintsAsString), " ")
	if len(parts) == 0 {
		return ast.ScalarType{}, errorWithCueRef(v, "something went very wrong while formatting a number expression into a string")
	}

	// dirty way of preserving the actual type from cue
	// FIXME: fails if the type has a custom bound that further restricts the values
	// IE: uint8 & < 12 will be printed as "uint & < 12
	var numberType ast.Kind
	switch ast.Kind(parts[0]) {
	case ast.KindFloat32, ast.KindFloat64:
		numberType = ast.Kind(parts[0])
	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64:
		numberType = ast.Kind(parts[0])
	case ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		numberType = ast.Kind(parts[0])
	case "uint":
		numberType = ast.KindUint64
	case "int":
		numberType = ast.KindInt64
	case "number":
		numberType = ast.KindFloat64
	default:
		return ast.ScalarType{}, errorWithCueRef(v, "unknown number type '%s'", parts[0])
	}

	typeDef := ast.ScalarType{
		ScalarKind: numberType,
	}

	// v.IsConcrete() being true means we're looking at a constant/known value
	if v.IsConcrete() {
		val, err := cueConcreteToScalar(v)
		if err != nil {
			return typeDef, err
		}

		typeDef.Value = val
	}

	// If the default (all lists have a default, usually self, ugh) differs from the
	// input list, peel it off. Otherwise our AnyIndex lookup may end up getting
	// sent on the wrong path.
	defv, _ := v.Default()
	if !defv.Equals(v) {
		_, dvals := v.Expr()
		v = dvals[0]
	}

	// extract constraints
	constraints, err := g.declareNumberConstraints(v)
	if err != nil {
		return ast.ScalarType{}, err
	}

	typeDef.Constraints = constraints

	return typeDef, nil
}

func (g *newGenerator) declareNumberConstraints(v cue.Value) ([]ast.TypeConstraint, error) {
	// typeAndConstraints can contain the following cue expressions:
	// 	- number
	// 	- int|float, number, upper bound, lower bound
	typeAndConstraints := appendSplit(nil, cue.AndOp, v)

	// nothing to do
	if len(typeAndConstraints) == 1 {
		return nil, nil
	}

	constraints := make([]ast.TypeConstraint, 0, len(typeAndConstraints))

	constraintsStartIndex := 1

	// don't include type-related constraints
	if len(typeAndConstraints) > 1 && typeAndConstraints[0].IncompleteKind() != cue.NumberKind {
		constraintsStartIndex = 3
	}

	for _, s := range typeAndConstraints[constraintsStartIndex:] {
		constraint, err := g.extractConstraint(s)
		if err != nil {
			return nil, err
		}

		constraints = append(constraints, constraint)
	}

	return constraints, nil
}

func (g *newGenerator) extractConstraint(v cue.Value) (ast.TypeConstraint, error) {
	toConstraint := func(operator string, arg cue.Value) (ast.TypeConstraint, error) {
		scalar, err := cueConcreteToScalar(arg)
		if err != nil {
			return ast.TypeConstraint{}, err
		}

		return ast.TypeConstraint{
			Op:   operator,
			Args: []any{scalar},
		}, nil
	}

	switch op, a := v.Expr(); op {
	case cue.LessThanOp:
		return toConstraint("<", a[0])
	case cue.LessThanEqualOp:
		return toConstraint("<=", a[0])
	case cue.GreaterThanOp:
		return toConstraint(">", a[0])
	case cue.GreaterThanEqualOp:
		return toConstraint(">=", a[0])
	case cue.NotEqualOp:
		return toConstraint("!=", a[0])
	default:
		return ast.TypeConstraint{}, errorWithCueRef(v, "unsupported op for number %v", op)
	}
}

func (g *newGenerator) declareList(v cue.Value) (ast.Type, error) {
	i, err := v.List()
	if err != nil {
		return nil, err
	}

	typeDef := ast.ArrayType{
		// FIXME: we set a default type because our logic is broken
		ValueType: ast.ScalarType{
			ScalarKind: ast.KindAny,
		},
	}

	// works only for a closed/concrete list
	if v.IsConcrete() {
		// fixme: this is wrong
		for i.Next() {
			node, err := g.declareNode(i.Value())
			if err != nil {
				return nil, err
			}

			typeDef.ValueType = node
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
		return nil, errorWithCueRef(v, "open list must have a type")
	}

	expr, err := g.declareNode(e)
	if err != nil {
		return nil, err
	}

	typeDef.ValueType = expr

	return typeDef, nil
}
