package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DefaultAsTyped)(nil)

// DefaultAsTyped converts the raw `Default any` carried by every ast.Type into
// a typed ast.TypedDefault. The pass is recursive: an array default visits
// each item against the array's value type, a map default visits each entry
// against the map's value type, and a struct default visits each known field
// against its declared type.
//
// Code generators should prefer Type.TypedDefault over Type.Default to avoid
// unsafe type assertions like def.Default.(map[string]any).
type DefaultAsTyped struct{}

func (pass *DefaultAsTyped) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := &Visitor{
		OnScalar: pass.processScalar,
		OnEnum:   pass.processEnum,
		OnArray:  pass.processArray,
		OnMap:    pass.processMap,
		OnStruct: pass.processStruct,
		OnRef:    pass.processRef,
	}
	return visitor.VisitSchemas(schemas)
}

func (pass *DefaultAsTyped) processScalar(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if def.Default != nil {
		def.TypedDefault = ast.NewScalarDefault(def.Default)
	}
	return def, nil
}

func (pass *DefaultAsTyped) processEnum(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if def.Default != nil {
		def.TypedDefault = ast.NewEnumDefault(def.Default)
	}
	return def, nil
}

func (pass *DefaultAsTyped) processArray(visitor *Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
	valueType, err := visitor.VisitType(schema, def.AsArray().ValueType)
	if err != nil {
		return def, err
	}
	def.Array.ValueType = valueType

	if def.Default != nil {
		def.TypedDefault = ast.NewArrayDefault(toArrayItems(def.AsArray().ValueType, def.Default))
	}
	return def, nil
}

func (pass *DefaultAsTyped) processMap(visitor *Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
	valueType, err := visitor.VisitType(schema, def.AsMap().ValueType)
	if err != nil {
		return def, err
	}
	def.Map.ValueType = valueType

	if def.Default != nil {
		def.TypedDefault = ast.NewMapDefault(toMapEntries(def.AsMap().ValueType, def.Default))
	}
	return def, nil
}

func (pass *DefaultAsTyped) processStruct(visitor *Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
	for i, field := range def.Struct.Fields {
		visitedField, err := visitor.VisitStructField(schema, field)
		if err != nil {
			return def, err
		}
		def.Struct.Fields[i] = visitedField
	}

	if def.Default != nil {
		def.TypedDefault = ast.NewStructDefault(toStructFields(def.AsStruct(), def.Default))
	}
	return def, nil
}

func (pass *DefaultAsTyped) processRef(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if def.Default == nil {
		return def, nil
	}

	def.TypedDefault = ast.NewRefDefault(inferTypedDefault(def.Default))
	return def, nil
}

func toArrayItems(valueType ast.Type, raw any) []ast.TypedDefault {
	rawItems, ok := raw.([]any)
	if !ok {
		return nil
	}

	items := make([]ast.TypedDefault, 0, len(rawItems))
	for _, rawItem := range rawItems {
		items = append(items, valueToTypedDefault(valueType, rawItem))
	}
	return items
}

func toMapEntries(valueType ast.Type, raw any) map[string]ast.TypedDefault {
	rawEntries, ok := raw.(map[string]any)
	if !ok {
		return nil
	}

	entries := make(map[string]ast.TypedDefault, len(rawEntries))
	for key, val := range rawEntries {
		entries[key] = valueToTypedDefault(valueType, val)
	}
	return entries
}

func toStructFields(structType ast.StructType, raw any) map[string]ast.TypedDefault {
	rawFields, ok := raw.(map[string]any)
	if !ok {
		return nil
	}

	fields := make(map[string]ast.TypedDefault, len(rawFields))
	for _, field := range structType.Fields {
		rawVal, exists := rawFields[field.Name]
		if !exists {
			continue
		}
		fields[field.Name] = valueToTypedDefault(field.Type, rawVal)
	}
	return fields
}

// valueToTypedDefault converts a raw default value into a TypedDefault using
// the surrounding declared type to choose the right shape.
func valueToTypedDefault(declared ast.Type, raw any) ast.TypedDefault {
	switch declared.Kind {
	case ast.KindArray:
		return ast.TypedDefault{
			Kind:  ast.KindArray,
			Array: &ast.ArrayDefault{Items: toArrayItems(declared.AsArray().ValueType, raw)},
		}
	case ast.KindMap:
		return ast.TypedDefault{
			Kind: ast.KindMap,
			Map:  &ast.MapDefault{Entries: toMapEntries(declared.AsMap().ValueType, raw)},
		}
	case ast.KindStruct:
		return ast.TypedDefault{
			Kind:   ast.KindStruct,
			Struct: &ast.StructDefault{Fields: toStructFields(declared.AsStruct(), raw)},
		}
	case ast.KindEnum:
		return ast.TypedDefault{Kind: ast.KindEnum, Enum: &ast.EnumDefault{Value: raw}}
	case ast.KindRef:
		return ast.TypedDefault{
			Kind: ast.KindRef,
			Ref:  &ast.RefDefault{Inner: inferTypedDefault(raw)},
		}
	default:
		return ast.TypedDefault{Kind: ast.KindScalar, Scalar: &ast.ScalarDefault{Value: raw}}
	}
}

// inferTypedDefault best-effort infers the shape of a raw value when the
// declared type cannot be resolved (e.g. behind a Ref).
func inferTypedDefault(raw any) ast.TypedDefault {
	switch v := raw.(type) {
	case map[string]any:
		fields := make(map[string]ast.TypedDefault, len(v))
		for k, val := range v {
			fields[k] = inferTypedDefault(val)
		}
		return ast.TypedDefault{Kind: ast.KindStruct, Struct: &ast.StructDefault{Fields: fields}}
	case []any:
		items := make([]ast.TypedDefault, 0, len(v))
		for _, val := range v {
			items = append(items, inferTypedDefault(val))
		}
		return ast.TypedDefault{Kind: ast.KindArray, Array: &ast.ArrayDefault{Items: items}}
	default:
		return ast.TypedDefault{Kind: ast.KindScalar, Scalar: &ast.ScalarDefault{Value: raw}}
	}
}

