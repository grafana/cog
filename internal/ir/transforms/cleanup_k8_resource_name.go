package transforms

import (
	"strings"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*CleanupK8ResourceNames)(nil)

type CleanupK8ResourceNames struct {
	PrefixToRemove string
}

func (pass *CleanupK8ResourceNames) Process(schemas ir.Schemas) (ir.Schemas, error) {
	visitor := Visitor{
		OnObject:      pass.parseObject,
		OnRef:         pass.parseReference,
		OnConstantRef: pass.parseConstantReference,
		OnStructField: pass.parseField,
		OnDisjunction: pass.parseDisjunction,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *CleanupK8ResourceNames) parseObject(visitor *Visitor, schema *ir.Schema, object ir.Object) (ir.Object, error) {
	newObject := object
	newObject.Name = pass.cleanupName(object.Name)
	newObject.SelfRef = ir.NewRef(newObject.SelfRef.ReferredPkg, pass.cleanupName(newObject.SelfRef.ReferredType)).AsRef()

	if !newObject.Type.IsStruct() {
		return newObject, nil
	}

	for i, f := range object.Type.AsStruct().Fields {
		t, err := visitor.VisitType(schema, f.Type)
		if err != nil {
			return ir.Object{}, err
		}
		newObject.Type.AsStruct().Fields[i].Type = t
	}

	return newObject, nil
}

func (pass *CleanupK8ResourceNames) parseReference(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
	refType := pass.cleanupName(def.AsRef().ReferredType)
	return ir.NewRef(def.AsRef().ReferredPkg, refType), nil
}

func (pass *CleanupK8ResourceNames) parseConstantReference(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
	refType := pass.cleanupName(def.AsConstantRef().ReferredType)
	return ir.NewConstantReferenceType(def.AsConstantRef().ReferredPkg, refType, def.AsConstantRef().ReferenceValue), nil
}

func (pass *CleanupK8ResourceNames) parseField(_ *Visitor, _ *ir.Schema, field ir.StructField) (ir.StructField, error) {
	field.Name = pass.cleanupName(field.Name)
	return field, nil
}

func (pass *CleanupK8ResourceNames) parseDisjunction(visitor *Visitor, schema *ir.Schema, def ir.Type) (ir.Type, error) {
	for i, b := range def.AsDisjunction().Branches {
		t, err := visitor.VisitType(schema, b)
		if err != nil {
			return ir.Type{}, err
		}
		def.AsDisjunction().Branches[i] = t
	}

	return def, nil
}

func (pass *CleanupK8ResourceNames) cleanupName(s string) string {
	elements := strings.Split(s, ".")
	lastElement := elements[len(elements)-1]
	lastElement = strings.TrimPrefix(lastElement, pass.PrefixToRemove)
	return lastElement
}
