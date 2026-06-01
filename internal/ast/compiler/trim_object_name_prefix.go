package compiler

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*TrimObjectNamePrefix)(nil)

// TrimObjectNamePrefix removes the given prefix from every object's name.
type TrimObjectNamePrefix struct {
	Prefix string
}

func (pass *TrimObjectNamePrefix) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	if pass.Prefix == "" {
		return schemas, nil
	}

	for _, schema := range schemas {
		schema.EntryPoint = strings.TrimPrefix(schema.EntryPoint, pass.Prefix)
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

func (pass *TrimObjectNamePrefix) processObject(visitor *Visitor, schema *ast.Schema, object ast.Object) (ast.Object, error) {
	var err error

	if strings.HasPrefix(object.Name, pass.Prefix) {
		originalName := object.Name
		object.Name = strings.TrimPrefix(originalName, pass.Prefix)
		object.SelfRef.ReferredType = object.Name
		object.AddToPassesTrail(fmt.Sprintf("TrimObjectNamePrefix[%s → %s]", originalName, object.Name))
	}

	object.Type, err = visitor.VisitType(schema, object.Type)
	if err != nil {
		return ast.Object{}, err
	}

	return object, nil
}

func (pass *TrimObjectNamePrefix) processStruct(visitor *Visitor, schema *ast.Schema, structDef ast.Type) (ast.Type, error) {
	var err error
	for i, field := range structDef.Struct.Fields {
		structDef.Struct.Fields[i], err = visitor.VisitStructField(schema, field)
		if err != nil {
			return ast.Type{}, err
		}
	}

	if structDef.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) {
		var modified bool
		disjunction := structDef.Hints[ast.HintDiscriminatedDisjunctionOfRefs].(ast.DisjunctionType)
		disjunction.DiscriminatorMapping, modified = pass.processDisjunctionMapping(disjunction.DiscriminatorMapping)
		structDef.Hints[ast.HintDiscriminatedDisjunctionOfRefs] = disjunction

		if modified {
			structDef.AddToPassesTrail(fmt.Sprintf("TrimObjectNamePrefix[prefix=%s]", pass.Prefix))
		}
	}

	return structDef, nil
}

func (pass *TrimObjectNamePrefix) processDisjunction(visitor *Visitor, schema *ast.Schema, disjunction ast.Type) (ast.Type, error) {
	var modified bool

	disjunction.Disjunction.DiscriminatorMapping, modified = pass.processDisjunctionMapping(disjunction.Disjunction.DiscriminatorMapping)
	if modified {
		disjunction.AddToPassesTrail(fmt.Sprintf("TrimObjectNamePrefix[prefix=%s]", pass.Prefix))
	}

	var err error
	for i, branch := range disjunction.Disjunction.Branches {
		disjunction.Disjunction.Branches[i], err = visitor.VisitType(schema, branch)
		if err != nil {
			return ast.Type{}, err
		}
	}

	return disjunction, nil
}

func (pass *TrimObjectNamePrefix) processDisjunctionMapping(discriminatorMapping map[string]string) (map[string]string, bool) {
	modified := false
	newMapping := make(map[string]string, len(discriminatorMapping))
	for discriminator, typeName := range discriminatorMapping {
		if strings.HasPrefix(typeName, pass.Prefix) {
			newMapping[discriminator] = strings.TrimPrefix(typeName, pass.Prefix)
			modified = true
		} else {
			newMapping[discriminator] = typeName
		}
	}

	return newMapping, modified
}

func (pass *TrimObjectNamePrefix) processRef(_ *Visitor, _ *ast.Schema, ref ast.Type) (ast.Type, error) {
	if !strings.HasPrefix(ref.Ref.ReferredType, pass.Prefix) {
		return ref, nil
	}

	originalName := ref.Ref.ReferredType
	ref.Ref.ReferredType = strings.TrimPrefix(originalName, pass.Prefix)
	ref.AddToPassesTrail(fmt.Sprintf("TrimObjectNamePrefix[%s → %s]", originalName, ref.Ref.ReferredType))

	return ref, nil
}

func (pass *TrimObjectNamePrefix) processEnum(_ *Visitor, _ *ast.Schema, enum ast.Type) (ast.Type, error) {
	values := make([]ast.EnumValue, 0, len(enum.AsEnum().Values))
	for _, val := range enum.AsEnum().Values {
		name := tools.UpperCamelCase(strings.TrimPrefix(val.Name, tools.UpperCamelCase(pass.Prefix)))

		values = append(values, ast.EnumValue{
			Type:  val.Type,
			Name:  name,
			Value: val.Value,
		})
	}

	enum.Enum.Values = values

	return enum, nil
}

func (pass *TrimObjectNamePrefix) processConstantRef(_ *Visitor, _ *ast.Schema, ref ast.Type) (ast.Type, error) {
	if !strings.HasPrefix(ref.ConstantReference.ReferredType, pass.Prefix) {
		return ref, nil
	}

	originalName := ref.ConstantReference.ReferredType
	ref.ConstantReference.ReferredType = strings.TrimPrefix(originalName, pass.Prefix)
	ref.AddToPassesTrail(fmt.Sprintf("TrimObjectNamePrefix[%s → %s]", originalName, ref.ConstantReference.ReferredType))

	return ref, nil
}
