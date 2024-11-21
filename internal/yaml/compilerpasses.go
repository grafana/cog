package yaml

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
)

type CompilerPass struct {
	EntrypointIdentification *EntrypointIdentification `yaml:"entrypoint_identification"`
	DataqueryIdentification  *DataqueryIdentification  `yaml:"dataquery_identification"`
	Unspec                   *Unspec                   `yaml:"unspec"`
	SetDatasourceToDataquery *SetDatasourceToDataquery `yaml:"set_datasource_to_dataquery"`
	FieldsSetDefault         *FieldsSetDefault         `yaml:"fields_set_default"`
	FieldsSetRequired        *FieldsSetRequired        `yaml:"fields_set_required"`
	FieldsSetNotRequired     *FieldsSetNotRequired     `yaml:"fields_set_not_required"`
	Omit                     *Omit                     `yaml:"omit"`
	AddFields                *AddFields                `yaml:"add_fields"`
	NameAnonymousStruct      *NameAnonymousStruct      `yaml:"name_anonymous_struct"`
	RenameObject             *RenameObject             `yaml:"rename_object"`
	RetypeObject             *RetypeObject             `yaml:"retype_object"`
	HintObject               *HintObject               `yaml:"hint_object"`
	RetypeField              *RetypeField              `yaml:"retype_field"`
	SchemaSetIdentifier      *SchemaSetIdentifier      `yaml:"schema_set_identifier"`
	DuplicateObject          *DuplicateObject          `yaml:"duplicate_object"`

	AnonymousStructsToNamed *AnonymousStructsToNamed `yaml:"anonymous_structs_to_named"`

	DisjunctionToType                       *DisjunctionToType                       `yaml:"disjunction_to_type"`
	DisjunctionOfAnonymousStructsToExplicit *DisjunctionOfAnonymousStructsToExplicit `yaml:"disjunction_of_anonymous_structs_to_explicit"`
	DisjunctionInferMapping                 *DisjunctionInferMapping                 `yaml:"disjunction_infer_mapping"`
	DisjunctionWithConstantToDefault        *DisjunctionWithConstantToDefault        `yaml:"disjunction_with_constant_to_default"`

	Cloudwatch *Cloudwatch `yaml:"cloudwatch"`
}

func (pass CompilerPass) AsCompilerPass() (compiler.Pass, error) {
	if pass.EntrypointIdentification != nil {
		return pass.EntrypointIdentification.AsCompilerPass(), nil
	}
	if pass.DataqueryIdentification != nil {
		return pass.DataqueryIdentification.AsCompilerPass(), nil
	}
	if pass.Unspec != nil {
		return pass.Unspec.AsCompilerPass(), nil
	}
	if pass.SetDatasourceToDataquery != nil {
		return pass.SetDatasourceToDataquery.AsCompilerPass(), nil
	}
	if pass.FieldsSetDefault != nil {
		return pass.FieldsSetDefault.AsCompilerPass()
	}
	if pass.FieldsSetRequired != nil {
		return pass.FieldsSetRequired.AsCompilerPass()
	}
	if pass.FieldsSetNotRequired != nil {
		return pass.FieldsSetNotRequired.AsCompilerPass()
	}
	if pass.Omit != nil {
		return pass.Omit.AsCompilerPass()
	}
	if pass.AddFields != nil {
		return pass.AddFields.AsCompilerPass()
	}
	if pass.NameAnonymousStruct != nil {
		return pass.NameAnonymousStruct.AsCompilerPass()
	}
	if pass.RetypeObject != nil {
		return pass.RetypeObject.AsCompilerPass()
	}
	if pass.HintObject != nil {
		return pass.HintObject.AsCompilerPass()
	}
	if pass.RenameObject != nil {
		return pass.RenameObject.AsCompilerPass()
	}
	if pass.RetypeField != nil {
		return pass.RetypeField.AsCompilerPass()
	}
	if pass.SchemaSetIdentifier != nil {
		return pass.SchemaSetIdentifier.AsCompilerPass()
	}
	if pass.DuplicateObject != nil {
		return pass.DuplicateObject.AsCompilerPass()
	}

	if pass.AnonymousStructsToNamed != nil {
		return pass.AnonymousStructsToNamed.AsCompilerPass()
	}

	if pass.DisjunctionToType != nil {
		return pass.DisjunctionToType.AsCompilerPass()
	}
	if pass.DisjunctionOfAnonymousStructsToExplicit != nil {
		return pass.DisjunctionOfAnonymousStructsToExplicit.AsCompilerPass()
	}
	if pass.DisjunctionInferMapping != nil {
		return pass.DisjunctionInferMapping.AsCompilerPass()
	}
	if pass.DisjunctionWithConstantToDefault != nil {
		return pass.DisjunctionWithConstantToDefault.AsCompilerPass()
	}

	if pass.Cloudwatch != nil {
		return pass.Cloudwatch.AsCompilerPass(), nil
	}

	return nil, fmt.Errorf("empty compiler pass")
}

type EntrypointIdentification struct {
}

func (pass EntrypointIdentification) AsCompilerPass() *compiler.InferEntrypoint {
	return &compiler.InferEntrypoint{}
}

type DataqueryIdentification struct {
}

func (pass DataqueryIdentification) AsCompilerPass() *compiler.DataqueryIdentification {
	return &compiler.DataqueryIdentification{}
}

type Unspec struct {
}

func (pass Unspec) AsCompilerPass() *compiler.Unspec {
	return &compiler.Unspec{}
}

type SetDatasourceToDataquery struct{}

func (pass SetDatasourceToDataquery) AsCompilerPass() *compiler.SetDatasourceToDataquery {
	return &compiler.SetDatasourceToDataquery{}
}

type FieldsSetDefault struct {
	Defaults map[string]any // Expected format: [package].[object].[field] → value
}

func (pass FieldsSetDefault) AsCompilerPass() (*compiler.FieldsSetDefault, error) {
	defaults := make(map[compiler.FieldReference]any, len(pass.Defaults))

	for ref, value := range pass.Defaults {
		fieldRef, err := compiler.FieldReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		defaults[fieldRef] = value
	}

	return &compiler.FieldsSetDefault{DefaultValues: defaults}, nil
}

type FieldsSetRequired struct {
	Fields []string // Expected format: [package].[object].[field]
}

func (pass FieldsSetRequired) AsCompilerPass() (*compiler.FieldsSetRequired, error) {
	fieldRefs := make([]compiler.FieldReference, 0, len(pass.Fields))

	for _, ref := range pass.Fields {
		fieldRef, err := compiler.FieldReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		fieldRefs = append(fieldRefs, fieldRef)
	}

	return &compiler.FieldsSetRequired{Fields: fieldRefs}, nil
}

type FieldsSetNotRequired struct {
	Fields []string // Expected format: [package].[object].[field]
}

func (pass FieldsSetNotRequired) AsCompilerPass() (*compiler.FieldsSetNotRequired, error) {
	fieldRefs := make([]compiler.FieldReference, 0, len(pass.Fields))

	for _, ref := range pass.Fields {
		fieldRef, err := compiler.FieldReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		fieldRefs = append(fieldRefs, fieldRef)
	}

	return &compiler.FieldsSetNotRequired{Fields: fieldRefs}, nil
}

type Omit struct {
	Objects []string // Expected format: [package].[object]
}

func (pass Omit) AsCompilerPass() (*compiler.Omit, error) {
	objectRefs := make([]compiler.ObjectReference, 0, len(pass.Objects))

	for _, ref := range pass.Objects {
		objectRef, err := compiler.ObjectReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		objectRefs = append(objectRefs, objectRef)
	}

	return &compiler.Omit{Objects: objectRefs}, nil
}

type AddFields struct {
	// Expected format: [package].[object]
	To     string
	Fields []ast.StructField
}

func (pass AddFields) AsCompilerPass() (*compiler.AddFields, error) {
	objectRef, err := compiler.ObjectReferenceFromString(pass.To)
	if err != nil {
		return nil, err
	}

	return &compiler.AddFields{
		Object: objectRef,
		Fields: pass.Fields,
	}, nil
}

type NameAnonymousStruct struct {
	Field string // Expected format: [package].[object].[field]
	As    string
}

func (pass NameAnonymousStruct) AsCompilerPass() (*compiler.NameAnonymousStruct, error) {
	fieldRef, err := compiler.FieldReferenceFromString(pass.Field)
	if err != nil {
		return nil, err
	}

	return &compiler.NameAnonymousStruct{
		Field: fieldRef,
		As:    pass.As,
	}, nil
}

type RetypeObject struct {
	Object   string // Expected format: [package].[object]
	As       ast.Type
	Comments []string
}

func (pass RetypeObject) AsCompilerPass() (*compiler.RetypeObject, error) {
	objectRef, err := compiler.ObjectReferenceFromString(pass.Object)
	if err != nil {
		return nil, err
	}

	return &compiler.RetypeObject{
		Object:   objectRef,
		As:       pass.As,
		Comments: pass.Comments,
	}, nil
}

type HintObject struct {
	Object string // Expected format: [package].[object]
	Hints  ast.JenniesHints
}

func (pass HintObject) AsCompilerPass() (*compiler.HintObject, error) {
	objectRef, err := compiler.ObjectReferenceFromString(pass.Object)
	if err != nil {
		return nil, err
	}

	return &compiler.HintObject{
		Object: objectRef,
		Hints:  pass.Hints,
	}, nil
}

type DuplicateObject struct {
	Object     string // Expected format: [package].[object]
	As         string
	OmitFields []string `yaml:"omit_fields"`
}

func (pass DuplicateObject) AsCompilerPass() (*compiler.DuplicateObject, error) {
	objectRef, err := compiler.ObjectReferenceFromString(pass.Object)
	if err != nil {
		return nil, err
	}

	destinationRef, err := compiler.ObjectReferenceFromString(pass.As)
	if err != nil {
		return nil, err
	}

	return &compiler.DuplicateObject{
		Object:     objectRef,
		As:         destinationRef,
		OmitFields: pass.OmitFields,
	}, nil
}

type RenameObject struct {
	From string // Expected format: [package].[object]
	To   string
}

func (pass RenameObject) AsCompilerPass() (*compiler.RenameObject, error) {
	objectRef, err := compiler.ObjectReferenceFromString(pass.From)
	if err != nil {
		return nil, err
	}

	return &compiler.RenameObject{
		From: objectRef,
		To:   pass.To,
	}, nil
}

type RetypeField struct {
	Field    string // Expected format: [package].[object].[field]
	As       ast.Type
	Comments []string
}

func (pass RetypeField) AsCompilerPass() (*compiler.RetypeField, error) {
	fieldRef, err := compiler.FieldReferenceFromString(pass.Field)
	if err != nil {
		return nil, err
	}

	return &compiler.RetypeField{
		Field:    fieldRef,
		As:       pass.As,
		Comments: pass.Comments,
	}, nil
}

type SchemaSetIdentifier struct {
	Package    string
	Identifier string
}

func (pass SchemaSetIdentifier) AsCompilerPass() (*compiler.SchemaSetIdentifier, error) {
	return &compiler.SchemaSetIdentifier{
		Package:    pass.Package,
		Identifier: pass.Identifier,
	}, nil
}

type AnonymousStructsToNamed struct {
}

func (pass AnonymousStructsToNamed) AsCompilerPass() (*compiler.AnonymousStructsToNamed, error) {
	return &compiler.AnonymousStructsToNamed{}, nil
}

type DisjunctionToType struct {
}

func (pass DisjunctionToType) AsCompilerPass() (*compiler.DisjunctionToType, error) {
	return &compiler.DisjunctionToType{}, nil
}

type DisjunctionOfAnonymousStructsToExplicit struct {
}

func (pass DisjunctionOfAnonymousStructsToExplicit) AsCompilerPass() (*compiler.DisjunctionOfAnonymousStructsToExplicit, error) {
	return &compiler.DisjunctionOfAnonymousStructsToExplicit{}, nil
}

type DisjunctionInferMapping struct {
}

func (pass DisjunctionInferMapping) AsCompilerPass() (*compiler.DisjunctionInferMapping, error) {
	return &compiler.DisjunctionInferMapping{}, nil
}

type DisjunctionWithConstantToDefault struct {
}

func (pass DisjunctionWithConstantToDefault) AsCompilerPass() (*compiler.DisjunctionWithConstantToDefault, error) {
	return &compiler.DisjunctionWithConstantToDefault{}, nil
}

type Cloudwatch struct {
}

func (pass Cloudwatch) AsCompilerPass() *compiler.Cloudwatch {
	return &compiler.Cloudwatch{}
}
