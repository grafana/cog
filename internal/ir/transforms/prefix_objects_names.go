package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/tools"
)

var _ Transform = (*PrefixObjectNames)(nil)

// PrefixObjectNames adds the given prefix to every object's name.
type PrefixObjectNames struct {
	Prefix string
}

func (pass *PrefixObjectNames) Process(schemas ir.Schemas) (ir.Schemas, error) {
	if pass.Prefix == "" {
		return schemas, nil
	}

	visitor := &Visitor{
		OnObject:      pass.processObject,
		OnStruct:      pass.processStruct,
		OnRef:         pass.processRef,
		OnEnum:        pass.processEnum,
		OnDisjunction: pass.processDisjunction,
		OnConstantRef: pass.processConstantRef,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *PrefixObjectNames) processObject(visitor *Visitor, schema *ir.Schema, object ir.Object) (ir.Object, error) {
	var err error

	originalName := object.Name
	object.Name = pass.Prefix + originalName
	object.SelfRef.ReferredType = object.Name
	object.AddToPassesTrail(fmt.Sprintf("PrefixObjectNames[%s → %s]", originalName, object.Name))

	object.Type, err = visitor.VisitType(schema, object.Type)
	if err != nil {
		return ir.Object{}, err
	}

	return object, nil
}

func (pass *PrefixObjectNames) processStruct(visitor *Visitor, schema *ir.Schema, structDef ir.Type) (ir.Type, error) {
	var err error
	for i, field := range structDef.Struct.Fields {
		structDef.Struct.Fields[i], err = visitor.VisitStructField(schema, field)
		if err != nil {
			return ir.Type{}, err
		}
	}

	if structDef.HasHint(ir.HintDiscriminatedDisjunctionOfRefs) {
		disjunction := structDef.Hints[ir.HintDiscriminatedDisjunctionOfRefs].(ir.DisjunctionType)
		disjunction.DiscriminatorMapping = pass.processDisjunctionMapping(disjunction.DiscriminatorMapping)
		structDef.Hints[ir.HintDiscriminatedDisjunctionOfRefs] = disjunction
		structDef.AddToPassesTrail(fmt.Sprintf("PrefixObjectNames[prefix=%s]", pass.Prefix))
	}

	return structDef, nil
}

func (pass *PrefixObjectNames) processDisjunction(visitor *Visitor, schema *ir.Schema, disjunction ir.Type) (ir.Type, error) {
	disjunction.Disjunction.DiscriminatorMapping = pass.processDisjunctionMapping(disjunction.Disjunction.DiscriminatorMapping)
	disjunction.AddToPassesTrail(fmt.Sprintf("PrefixObjectNames[prefix=%s]", pass.Prefix))

	var err error
	for i, branch := range disjunction.Disjunction.Branches {
		disjunction.Disjunction.Branches[i], err = visitor.VisitType(schema, branch)
		if err != nil {
			return ir.Type{}, err
		}
	}

	return disjunction, nil
}

func (pass *PrefixObjectNames) processDisjunctionMapping(discriminatorMapping map[string]string) map[string]string {
	newMapping := make(map[string]string, len(discriminatorMapping))
	for discriminator, typeName := range discriminatorMapping {
		newMapping[discriminator] = pass.Prefix + typeName
	}

	return newMapping
}

func (pass *PrefixObjectNames) processRef(_ *Visitor, _ *ir.Schema, ref ir.Type) (ir.Type, error) {
	originalName := ref.Ref.ReferredType
	ref.Ref.ReferredType = pass.Prefix + originalName
	ref.AddToPassesTrail(fmt.Sprintf("PrefixObjectNames[%s → %s]", originalName, ref.Ref.ReferredType))

	return ref, nil
}

func (pass *PrefixObjectNames) processEnum(_ *Visitor, _ *ir.Schema, enum ir.Type) (ir.Type, error) {
	values := make([]ir.EnumValue, 0, len(enum.AsEnum().Values))
	for _, val := range enum.AsEnum().Values {
		values = append(values, ir.EnumValue{
			Type:  val.Type,
			Name:  tools.UpperCamelCase(pass.Prefix) + tools.UpperCamelCase(val.Name),
			Value: val.Value,
		})
	}

	enum.Enum.Values = values

	return enum, nil
}

func (pass *PrefixObjectNames) processConstantRef(_ *Visitor, _ *ir.Schema, ref ir.Type) (ir.Type, error) {
	originalName := ref.ConstantReference.ReferredType
	ref.ConstantReference.ReferredType = pass.Prefix + originalName
	ref.AddToPassesTrail(fmt.Sprintf("PrefixObjectNames[%s → %s]", originalName, ref.ConstantReference.ReferredType))

	return ref, nil
}
