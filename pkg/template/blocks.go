package template

import (
	"fmt"

	"github.com/grafana/cog/pkg/ir"
)

// CustomObjectMethodAllBlock returns the name of the template block to use to
// define custom methods on all objects.
func CustomObjectMethodAllBlock() string {
	return "object_all_custom_methods"
}

// CustomObjectUnmarshalBlock returns the name of the template block to use to
// define a custom unmarshal function for given object.
func CustomObjectUnmarshalBlock(obj ir.Object) string {
	return fmt.Sprintf("object_%s_%s_custom_unmarshal", obj.SelfRef.ReferredPkg, obj.SelfRef.ReferredType)
}

// CustomObjectStrictUnmarshalBlock returns the name of the template block to
// use to define a custom — strict — unmarshal logic for a field.
func CustomObjectStrictUnmarshalBlock(obj ir.Object) string {
	return fmt.Sprintf("object_%s_%s_custom_strict_unmarshal", obj.SelfRef.ReferredPkg, obj.SelfRef.ReferredType)
}

// ExtraPackageDocsBlock returns the name of the template block to use to
// render extra content on the documentation page of the given package.
func ExtraPackageDocsBlock(schema *ir.Schema) string {
	return fmt.Sprintf("api_reference_package_%s_extra", schema.Package)
}

// ExtraObjectDocsBlock returns the name of the template block to use to
// render extra content on the documentation page of the given object.
func ExtraObjectDocsBlock(obj ir.Object) string {
	return fmt.Sprintf("api_reference_object_%s_%s_extra", obj.SelfRef.ReferredPkg, obj.SelfRef.ReferredType)
}

// ExtraBuilderDocsBlock returns the name of the template block to use to
// render additional content on the documentation page for the given builder.
func ExtraBuilderDocsBlock(builder ir.Builder) string {
	return fmt.Sprintf("api_reference_builder_%s_%s_extra", builder.Package, builder.Name)
}

// CustomObjectVariantBlock returns the name of the template block to use to
// add custom logic for all objects implementing the variant implemented by the
// current object.
func CustomObjectVariantBlock(object ir.Object) string {
	return fmt.Sprintf("object_variant_%s", object.Type.ImplementedVariant())
}

// CustomObjectMethodsBlock returns the name of the template block to use to
// add custom methods to the given object.
func CustomObjectMethodsBlock(obj ir.Object) string {
	return fmt.Sprintf("object_%s_%s_custom_methods", obj.SelfRef.ReferredPkg, obj.SelfRef.ReferredType)
}

// CustomSchemaVariantBlock returns the name of the template block to use to
// add custom logic for all schemas implementing the variant exposed by the
// given schema.
func CustomSchemaVariantBlock(schema *ir.Schema) string {
	return fmt.Sprintf("schema_variant_%s", schema.Metadata.Variant)
}

// VariantFieldUnmarshalBlock returns the name of the template block to use to
// define how to unmarshal fields of the given variant.
func VariantFieldUnmarshalBlock(variant string) string {
	return fmt.Sprintf("variant_%s_field_unmarshal", variant)
}

func DynamicFilesBlock() string {
	return "dynamic_files"
}
