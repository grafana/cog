package yaml

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/ir/transforms"
)

type Transform struct {
	EntrypointIdentification *EntrypointIdentification `yaml:"entrypoint_identification"`
	DataqueryIdentification  *DataqueryIdentification  `yaml:"dataquery_identification"`
	Unspec                   *Unspec                   `yaml:"unspec"`
	ReplaceReference         *ReplaceReference         `yaml:"replace_reference"`
	FieldsSetDefault         *FieldsSetDefault         `yaml:"fields_set_default"`
	FieldsSetRequired        *FieldsSetRequired        `yaml:"fields_set_required"`
	FieldsSetNotRequired     *FieldsSetNotRequired     `yaml:"fields_set_not_required"`
	Omit                     *Omit                     `yaml:"omit"`
	AddFields                *AddFields                `yaml:"add_fields"`
	NameAnonymousStruct      *NameAnonymousStruct      `yaml:"name_anonymous_struct"`
	AddObject                *AddObject                `yaml:"add_object"`
	RenameObject             *RenameObject             `yaml:"rename_object"`
	RetypeObject             *RetypeObject             `yaml:"retype_object"`
	HintObject               *HintObject               `yaml:"hint_object"`
	DeprecateObject          *DeprecateObject          `yaml:"deprecate_object"`
	RetypeField              *RetypeField              `yaml:"retype_field"`
	OmitFields               *OmitFields               `yaml:"omit_fields"`
	SchemaSetIdentifier      *SchemaSetIdentifier      `yaml:"schema_set_identifier"`
	SchemaSetEntryPoint      *SchemaSetEntryPoint      `yaml:"schema_set_entry_point"`
	DuplicateObject          *DuplicateObject          `yaml:"duplicate_object"`
	TrimEnumValues           *TrimEnumValues           `yaml:"trim_enum_values"`
	ConstantToEnum           *ConstantToEnum           `yaml:"constant_to_enum"`
	ExtractK8ResourceNames   *CleanupK8ResourceNames   `yaml:"cleanup_k8_resource_names"`
	TrimObjectNamePrefix     *TrimObjectNamePrefix     `yaml:"trim_object_name_prefix"`
	SanitizeEnumMemberNames  *SanitizeEnumMemberNames  `yaml:"sanitize_enum_member_names"`

	AnonymousStructsToNamed     *AnonymousStructsToNamed `yaml:"anonymous_structs_to_named"`
	AnonymousEnumToExplicitType *AnonymousEnumsToNamed   `yaml:"anonymous_enum_to_named"`

	DisjunctionToType                       *DisjunctionToType                       `yaml:"disjunction_to_type"`
	DisjunctionOfAnonymousStructsToExplicit *DisjunctionOfAnonymousStructsToExplicit `yaml:"disjunction_of_anonymous_structs_to_explicit"`
	DisjunctionInferMapping                 *DisjunctionInferMapping                 `yaml:"disjunction_infer_mapping"`
	UndiscriminatedDisjunctionToAny         *UndiscriminatedDisjunctionToAny         `yaml:"undiscriminated_disjunction_to_any"`
	DisjunctionWithConstantToDefault        *DisjunctionWithConstantToDefault        `yaml:"disjunction_with_constant_to_default"`
}

func (pass Transform) AsTransform() (transforms.Transform, error) {
	if pass.EntrypointIdentification != nil {
		return pass.EntrypointIdentification.AsTransform(), nil
	}
	if pass.DataqueryIdentification != nil {
		return pass.DataqueryIdentification.AsTransform(), nil
	}
	if pass.Unspec != nil {
		return pass.Unspec.AsTransform(), nil
	}
	if pass.ReplaceReference != nil {
		return pass.ReplaceReference.AsTransform()
	}
	if pass.FieldsSetDefault != nil {
		return pass.FieldsSetDefault.AsTransform()
	}
	if pass.FieldsSetRequired != nil {
		return pass.FieldsSetRequired.AsTransform()
	}
	if pass.FieldsSetNotRequired != nil {
		return pass.FieldsSetNotRequired.AsTransform()
	}
	if pass.Omit != nil {
		return pass.Omit.AsTransform()
	}
	if pass.AddFields != nil {
		return pass.AddFields.AsTransform()
	}
	if pass.NameAnonymousStruct != nil {
		return pass.NameAnonymousStruct.AsTransform()
	}
	if pass.RetypeObject != nil {
		return pass.RetypeObject.AsTransform()
	}
	if pass.HintObject != nil {
		return pass.HintObject.AsTransform()
	}
	if pass.DeprecateObject != nil {
		return pass.DeprecateObject.AsTransform()
	}
	if pass.AddObject != nil {
		return pass.AddObject.AsTransform()
	}
	if pass.RenameObject != nil {
		return pass.RenameObject.AsTransform()
	}
	if pass.RetypeField != nil {
		return pass.RetypeField.AsTransform()
	}
	if pass.OmitFields != nil {
		return pass.OmitFields.AsTransform()
	}
	if pass.SchemaSetIdentifier != nil {
		return pass.SchemaSetIdentifier.AsTransform()
	}
	if pass.SchemaSetEntryPoint != nil {
		return pass.SchemaSetEntryPoint.AsTransform()
	}
	if pass.DuplicateObject != nil {
		return pass.DuplicateObject.AsTransform()
	}
	if pass.TrimEnumValues != nil {
		return pass.TrimEnumValues.AsTransform()
	}
	if pass.ConstantToEnum != nil {
		return pass.ConstantToEnum.AsTransform()
	}
	if pass.ExtractK8ResourceNames != nil {
		return pass.ExtractK8ResourceNames.AsTransform()
	}
	if pass.TrimObjectNamePrefix != nil {
		return pass.TrimObjectNamePrefix.AsTransform()
	}
	if pass.SanitizeEnumMemberNames != nil {
		return pass.SanitizeEnumMemberNames.AsTransform()
	}

	if pass.AnonymousStructsToNamed != nil {
		return pass.AnonymousStructsToNamed.AsTransform()
	}
	if pass.AnonymousEnumToExplicitType != nil {
		return pass.AnonymousEnumToExplicitType.AsTransform()
	}

	if pass.DisjunctionToType != nil {
		return pass.DisjunctionToType.AsTransform()
	}
	if pass.DisjunctionOfAnonymousStructsToExplicit != nil {
		return pass.DisjunctionOfAnonymousStructsToExplicit.AsTransform()
	}
	if pass.DisjunctionInferMapping != nil {
		return pass.DisjunctionInferMapping.AsTransform()
	}
	if pass.UndiscriminatedDisjunctionToAny != nil {
		return pass.UndiscriminatedDisjunctionToAny.AsTransform()
	}
	if pass.DisjunctionWithConstantToDefault != nil {
		return pass.DisjunctionWithConstantToDefault.AsTransform()
	}

	return nil, fmt.Errorf("empty compiler pass")
}

type EntrypointIdentification struct {
}

func (pass EntrypointIdentification) AsTransform() *transforms.InferEntrypoint {
	return &transforms.InferEntrypoint{}
}

type DataqueryIdentification struct {
}

func (pass DataqueryIdentification) AsTransform() *transforms.DataqueryIdentification {
	return &transforms.DataqueryIdentification{}
}

type Unspec struct {
}

func (pass Unspec) AsTransform() *transforms.Unspec {
	return &transforms.Unspec{}
}

type ReplaceReference struct {
	From string // Expected format: [package].[object]
	To   string // Expected format: [package].[object]
}

func (pass ReplaceReference) AsTransform() (*transforms.ReplaceReference, error) {
	fromRef, err := transforms.ObjectReferenceFromString(pass.From)
	if err != nil {
		return nil, err
	}

	toRef, err := transforms.ObjectReferenceFromString(pass.To)
	if err != nil {
		return nil, err
	}

	return &transforms.ReplaceReference{
		From: fromRef,
		To:   toRef,
	}, nil
}

type FieldsSetDefault struct {
	Defaults map[string]any // Expected format: [package].[object].[field] → value
}

func (pass FieldsSetDefault) AsTransform() (*transforms.FieldsSetDefault, error) {
	defaults := make(map[transforms.FieldReference]any, len(pass.Defaults))

	for ref, value := range pass.Defaults {
		fieldRef, err := transforms.FieldReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		defaults[fieldRef] = value
	}

	return &transforms.FieldsSetDefault{DefaultValues: defaults}, nil
}

type FieldsSetRequired struct {
	Fields []string // Expected format: [package].[object].[field]
}

func (pass FieldsSetRequired) AsTransform() (*transforms.FieldsSetRequired, error) {
	fieldRefs := make([]transforms.FieldReference, 0, len(pass.Fields))

	for _, ref := range pass.Fields {
		fieldRef, err := transforms.FieldReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		fieldRefs = append(fieldRefs, fieldRef)
	}

	return &transforms.FieldsSetRequired{Fields: fieldRefs}, nil
}

type FieldsSetNotRequired struct {
	Fields []string // Expected format: [package].[object].[field]
}

func (pass FieldsSetNotRequired) AsTransform() (*transforms.FieldsSetNotRequired, error) {
	fieldRefs := make([]transforms.FieldReference, 0, len(pass.Fields))

	for _, ref := range pass.Fields {
		fieldRef, err := transforms.FieldReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		fieldRefs = append(fieldRefs, fieldRef)
	}

	return &transforms.FieldsSetNotRequired{Fields: fieldRefs}, nil
}

type Omit struct {
	Objects []string // Expected format: [package].[object]
}

func (pass Omit) AsTransform() (*transforms.Omit, error) {
	objectRefs := make([]transforms.ObjectReference, 0, len(pass.Objects))

	for _, ref := range pass.Objects {
		objectRef, err := transforms.ObjectReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		objectRefs = append(objectRefs, objectRef)
	}

	return &transforms.Omit{Objects: objectRefs}, nil
}

type AddFields struct {
	// Expected format: [package].[object]
	To     string
	Fields []ir.StructField
}

func (pass AddFields) AsTransform() (*transforms.AddFields, error) {
	objectRef, err := transforms.ObjectReferenceFromString(pass.To)
	if err != nil {
		return nil, err
	}

	return &transforms.AddFields{
		Object: objectRef,
		Fields: pass.Fields,
	}, nil
}

type NameAnonymousStruct struct {
	Field string // Expected format: [package].[object].[field]
	As    string
}

func (pass NameAnonymousStruct) AsTransform() (*transforms.NameAnonymousStruct, error) {
	fieldRef, err := transforms.FieldReferenceFromString(pass.Field)
	if err != nil {
		return nil, err
	}

	return &transforms.NameAnonymousStruct{
		Field: fieldRef,
		As:    pass.As,
	}, nil
}

type RetypeObject struct {
	Object   string // Expected format: [package].[object]
	As       ir.Type
	Comments []string
}

func (pass RetypeObject) AsTransform() (*transforms.RetypeObject, error) {
	objectRef, err := transforms.ObjectReferenceFromString(pass.Object)
	if err != nil {
		return nil, err
	}

	return &transforms.RetypeObject{
		Object:   objectRef,
		As:       pass.As,
		Comments: pass.Comments,
	}, nil
}

type HintObject struct {
	Object string // Expected format: [package].[object]
	Hints  ir.JenniesHints
}

func (pass HintObject) AsTransform() (*transforms.HintObject, error) {
	objectRef, err := transforms.ObjectReferenceFromString(pass.Object)
	if err != nil {
		return nil, err
	}

	return &transforms.HintObject{
		Object: objectRef,
		Hints:  pass.Hints,
	}, nil
}

type DuplicateObject struct {
	Object     string // Expected format: [package].[object]
	As         string
	OmitFields []string `yaml:"omit_fields"`
}

func (pass DuplicateObject) AsTransform() (*transforms.DuplicateObject, error) {
	objectRef, err := transforms.ObjectReferenceFromString(pass.Object)
	if err != nil {
		return nil, err
	}

	destinationRef, err := transforms.ObjectReferenceFromString(pass.As)
	if err != nil {
		return nil, err
	}

	return &transforms.DuplicateObject{
		Object:     objectRef,
		As:         destinationRef,
		OmitFields: pass.OmitFields,
	}, nil
}

type AddObject struct {
	Object   string // Expected format: [package].[object]
	As       ir.Type
	Comments []string
}

func (pass AddObject) AsTransform() (*transforms.AddObject, error) {
	objectRef, err := transforms.ObjectReferenceFromString(pass.Object)
	if err != nil {
		return nil, err
	}

	return &transforms.AddObject{
		Object:   objectRef,
		As:       pass.As,
		Comments: pass.Comments,
	}, nil
}

type RenameObject struct {
	From string // Expected format: [package].[object]
	To   string
}

func (pass RenameObject) AsTransform() (*transforms.RenameObject, error) {
	objectRef, err := transforms.ObjectReferenceFromString(pass.From)
	if err != nil {
		return nil, err
	}

	return &transforms.RenameObject{
		From: objectRef,
		To:   pass.To,
	}, nil
}

type RetypeField struct {
	Field    string // Expected format: [package].[object].[field]
	As       ir.Type
	Comments []string
}

func (pass RetypeField) AsTransform() (*transforms.RetypeField, error) {
	fieldRef, err := transforms.FieldReferenceFromString(pass.Field)
	if err != nil {
		return nil, err
	}

	return &transforms.RetypeField{
		Field:    fieldRef,
		As:       pass.As,
		Comments: pass.Comments,
	}, nil
}

type OmitFields struct {
	Fields []string // Expected format: [package].[object].[field]
}

func (pass OmitFields) AsTransform() (*transforms.OmitFields, error) {
	fieldRefs := make([]transforms.FieldReference, 0, len(pass.Fields))
	for _, field := range pass.Fields {
		fieldRef, err := transforms.FieldReferenceFromString(field)
		if err != nil {
			return nil, err
		}
		fieldRefs = append(fieldRefs, fieldRef)
	}

	return &transforms.OmitFields{
		Fields: fieldRefs,
	}, nil
}

type SchemaSetIdentifier struct {
	Package    string
	Identifier string
}

func (pass SchemaSetIdentifier) AsTransform() (*transforms.SchemaSetIdentifier, error) {
	return &transforms.SchemaSetIdentifier{
		Package:    pass.Package,
		Identifier: pass.Identifier,
	}, nil
}

type SchemaSetEntryPoint struct {
	Package    string
	EntryPoint string `yaml:"entry_point"`
}

func (pass SchemaSetEntryPoint) AsTransform() (*transforms.SchemaSetEntrypoint, error) {
	return &transforms.SchemaSetEntrypoint{
		Package:    pass.Package,
		EntryPoint: pass.EntryPoint,
	}, nil
}

type AnonymousStructsToNamed struct {
}

func (pass AnonymousStructsToNamed) AsTransform() (*transforms.AnonymousStructsToNamed, error) {
	return &transforms.AnonymousStructsToNamed{}, nil
}

type AnonymousEnumsToNamed struct {
}

func (pass AnonymousEnumsToNamed) AsTransform() (*transforms.AnonymousEnumToExplicitType, error) {
	return &transforms.AnonymousEnumToExplicitType{}, nil
}

type DisjunctionToType struct {
}

func (pass DisjunctionToType) AsTransform() (*transforms.DisjunctionToType, error) {
	return &transforms.DisjunctionToType{}, nil
}

type DisjunctionOfAnonymousStructsToExplicit struct {
}

func (pass DisjunctionOfAnonymousStructsToExplicit) AsTransform() (*transforms.DisjunctionOfAnonymousStructsToExplicit, error) {
	return &transforms.DisjunctionOfAnonymousStructsToExplicit{}, nil
}

type DisjunctionInferMapping struct {
}

func (pass DisjunctionInferMapping) AsTransform() (*transforms.DisjunctionInferMapping, error) {
	return &transforms.DisjunctionInferMapping{}, nil
}

type UndiscriminatedDisjunctionToAny struct {
}

func (pass UndiscriminatedDisjunctionToAny) AsTransform() (*transforms.UndiscriminatedDisjunctionToAny, error) {
	return &transforms.UndiscriminatedDisjunctionToAny{}, nil
}

type DisjunctionWithConstantToDefault struct {
}

func (pass DisjunctionWithConstantToDefault) AsTransform() (*transforms.DisjunctionWithConstantToDefault, error) {
	return &transforms.DisjunctionWithConstantToDefault{}, nil
}

type TrimEnumValues struct{}

func (pass TrimEnumValues) AsTransform() (*transforms.TrimEnumValues, error) {
	return &transforms.TrimEnumValues{}, nil
}

type ConstantToEnum struct {
	Objects []string // Expected format: [package].[object]
}

func (pass ConstantToEnum) AsTransform() (*transforms.ConstantToEnum, error) {
	objectRefs := make([]transforms.ObjectReference, 0, len(pass.Objects))

	for _, ref := range pass.Objects {
		objectRef, err := transforms.ObjectReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		objectRefs = append(objectRefs, objectRef)
	}

	return &transforms.ConstantToEnum{Objects: objectRefs}, nil
}

type CleanupK8ResourceNames struct {
	PrefixToRemove string `yaml:"prefix_to_remove"`
}

func (pass CleanupK8ResourceNames) AsTransform() (*transforms.CleanupK8ResourceNames, error) {
	return &transforms.CleanupK8ResourceNames{
		PrefixToRemove: pass.PrefixToRemove,
	}, nil
}

type TrimObjectNamePrefix struct {
	Prefix string `yaml:"prefix"`
}

func (pass TrimObjectNamePrefix) AsTransform() (*transforms.TrimObjectNamePrefix, error) {
	return &transforms.TrimObjectNamePrefix{
		Prefix: pass.Prefix,
	}, nil
}

type SanitizeEnumMemberNames struct {
}

func (pass SanitizeEnumMemberNames) AsTransform() (*transforms.SanitizeEnumMemberNames, error) {
	return &transforms.SanitizeEnumMemberNames{}, nil
}

type DeprecateObject struct {
	Object  string // Expected format: [package].[object]
	Message string
}

func (pass DeprecateObject) AsTransform() (*transforms.DeprecateObject, error) {
	ref, err := transforms.ObjectReferenceFromString(pass.Object)
	if err != nil {
		return nil, err
	}

	return &transforms.DeprecateObject{
		Object:  ref,
		Message: pass.Message,
	}, nil
}
