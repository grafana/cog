package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"strings"
)

var _ Pass = (*ExtractK8sResourceName)(nil)

type ExtractK8sResourceName struct {
}

func (e *ExtractK8sResourceName) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := Visitor{
		OnObject:      e.parseObject,
		OnRef:         e.parseReference,
		OnConstantRef: e.parseConstantReference,
		OnStructField: e.parseField,
		OnDisjunction: e.parseDisjunction,
	}

	return visitor.VisitSchemas(schemas)
}

func (e *ExtractK8sResourceName) parseObject(visitor *Visitor, schema *ast.Schema, object ast.Object) (ast.Object, error) {
	newObject := object
	newObject.Name = getReferenceName(object.Name)
	newObject.SelfRef = ast.NewRef(newObject.SelfRef.ReferredPkg, getReferenceName(newObject.SelfRef.ReferredType)).AsRef()

	if !newObject.Type.IsStruct() {
		return newObject, nil
	}

	for i, f := range object.Type.AsStruct().Fields {
		t, err := visitor.VisitType(schema, f.Type)
		if err != nil {
			return ast.Object{}, err
		}
		newObject.Type.AsStruct().Fields[i].Type = t
	}

	return newObject, nil
}

func (*ExtractK8sResourceName) parseReference(visitor *Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
	refType := getReferenceName(def.AsRef().ReferredType)
	return ast.NewRef(def.AsRef().ReferredPkg, refType), nil
}

func (*ExtractK8sResourceName) parseConstantReference(visitor *Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
	refType := getReferenceName(def.AsConstantRef().ReferredType)
	return ast.NewConstantReferenceType(def.AsConstantRef().ReferredPkg, refType, def.AsConstantRef().ReferenceValue), nil
}

func (*ExtractK8sResourceName) parseField(visitor *Visitor, schema *ast.Schema, field ast.StructField) (ast.StructField, error) {
	field.Name = getReferenceName(field.Name)
	return field, nil
}

func (e *ExtractK8sResourceName) parseDisjunction(visitor *Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
	for i, b := range def.AsDisjunction().Branches {
		t, err := visitor.VisitType(schema, b)
		if err != nil {
			return ast.Type{}, err
		}
		def.AsDisjunction().Branches[i] = t
	}

	return def, nil
}

func getReferenceName(s string) string {
	elements := strings.Split(s, ".")
	return elements[len(elements)-1]
}
