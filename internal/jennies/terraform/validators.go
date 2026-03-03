package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/terraform/validators"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type validatorGenerator struct {
	Config Config

	context       languages.Context
	typeFormatter *typeFormatter
}

func newValidatorGenerator(config Config, context languages.Context, typeFormatter *typeFormatter) *validatorGenerator {
	return &validatorGenerator{
		Config:        config,
		context:       context,
		typeFormatter: typeFormatter,
	}
}

func (v *validatorGenerator) generateValidatorIfNeeded(name string, t ast.Type) (string, error) {
	switch t.Kind {
	case ast.KindEnum:
		return v.generateEnumValidator(name, t.AsEnum().Values)
	case ast.KindArray:
		return v.generateListValidator(name, t.AsArray())
	case ast.KindMap:
		return v.generateMapValidator(name, t.AsMap())
	case ast.KindRef:
		return v.generateReferenceValidator(name, t.AsRef())
	case ast.KindStruct:
		return v.generateStructValidator(name, t.AsStruct())
	case ast.KindScalar:
		return v.generateScalarValidator(name, t.AsScalar())
	default:
		fmt.Println(name, t.Kind)
		return "", nil
	}
}

func (v *validatorGenerator) generateEnumValidator(name string, enumValues []ast.EnumValue) (string, error) {
	if enumValues[0].Type.AsScalar().IsNumeric() {
		return v.generateTFScalarValidator(name, validators.NewInt64Validator(parseEnumValues(enumValues)))
	}
	return v.generateTFScalarValidator(name, validators.NewStringValidator(parseEnumValues(enumValues)))
}

func (v *validatorGenerator) generateListValidator(name string, a ast.ArrayType) (string, error) {
	if a.IsArrayOf(ast.KindRef, ast.KindStruct) {
		return "", nil
	}
	return "", nil
}

func (v *validatorGenerator) generateMapValidator(name string, m ast.MapType) (string, error) {
	if m.ValueType.IsAny() {
		return "", nil
	}

	return "", nil
}

func (v *validatorGenerator) generateReferenceValidator(name string, r ast.RefType) (string, error) {
	obj, ok := v.context.LocateObject(r.ReferredPkg, r.ReferredType)
	if !ok {
		return "func UnknownValidator() {}", nil
	}

	return v.generateValidatorIfNeeded(name, obj.Type)
}

func (v *validatorGenerator) generateStructValidator(name string, asStruct ast.StructType) (string, error) {
	var buffer strings.Builder
	for _, field := range asStruct.Fields {
		validator, err := v.generateValidatorIfNeeded(fmt.Sprintf("%s%s", name, tools.UpperCamelCase(field.Name)), field.Type)
		if err != nil {
			return "", err
		}
		if validator != "" {
			buffer.WriteString(validator)
		}
	}
	return buffer.String(), nil
}

func (v *validatorGenerator) generateScalarValidator(name string, scalar ast.ScalarType) (string, error) {
	if scalar.IsConcrete() {
		switch scalar.ScalarKind {
		case ast.KindString, ast.KindBytes, ast.KindNull:
			return v.generateTFScalarValidator(name, validators.NewStringValidator([]any{scalar.Value}))
		case ast.KindBool:
			return v.generateTFScalarValidator(name, validators.NewBoolValidator([]any{scalar.Value}))
		case ast.KindInt32, ast.KindUint32:
			return v.generateTFScalarValidator(name, validators.NewInt32Validator([]any{scalar.Value}))
		case ast.KindInt64, ast.KindUint64:
			return v.generateTFScalarValidator(name, validators.NewInt64Validator([]any{scalar.Value}))
		case ast.KindFloat32:
			return v.generateTFScalarValidator(name, validators.NewFloat32Validator([]any{scalar.Value}))
		case ast.KindFloat64:
			return v.generateTFScalarValidator(name, validators.NewFloat64Validator([]any{scalar.Value}))
		case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16:
			return v.generateTFScalarValidator(name, validators.NewNumberValidator([]any{scalar.Value}))
		}
	}
	if len(scalar.Constraints) == 0 {
		// TODO: Implement constraints
	}

	return "", nil
}

func (v *validatorGenerator) generateTFScalarValidator(name string, validateFunc validators.ValidatorFuncs) (string, error) {
	v.typeFormatter.packageMapper("context")
	v.typeFormatter.packageMapper("slices")
	v.typeFormatter.packageMapper("github.com/hashicorp/terraform-plugin-framework/schema/validator")
	var buffer strings.Builder
	buffer.WriteString(fmt.Sprintf("type %sValidator struct {}\n", tools.UpperCamelCase(name)))
	buffer.WriteString(fmt.Sprintf("func (v %sValidator) Description(_ context.Context) string {\n", tools.UpperCamelCase(name)))
	buffer.WriteString(fmt.Sprintf("return \"%s item must be one of '%s'\"\n", name, validateFunc.Values()))
	buffer.WriteString("}\n")

	buffer.WriteString(fmt.Sprintf("func (v %sValidator) MarkdownDescription(ctx context.Context) string {\nreturn v.Description(ctx)\n", tools.UpperCamelCase(name)))
	buffer.WriteString("}\n")

	buffer.WriteString(fmt.Sprintf("func (v %sValidator) %s(ctx context.Context, req validator.%s, resp *validator.%s){\n", tools.UpperCamelCase(name), validateFunc.ValidatorFunc(), validateFunc.ValidatorReq(), validateFunc.ValidatorResp()))
	buffer.WriteString(fmt.Sprintf("if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {\nreturn\n}\n"))
	buffer.WriteString(fmt.Sprintf("value := req.ConfigValue.%s()\n", validateFunc.ValidatorValue()))
	buffer.WriteString(fmt.Sprintf("if !slices.Contains(%s, value) {\n", validateFunc.ValuesMatcher()))
	buffer.WriteString(fmt.Sprintf("resp.Diagnostics.AddAttributeError(req.Path, v.Description(ctx), value)\n"))
	buffer.WriteString(fmt.Sprintf("}\n"))
	buffer.WriteString("}\n")

	return buffer.String(), nil
}

func parseEnumValues(enumValues []ast.EnumValue) []any {
	var values []any
	for _, v := range enumValues {
		values = append(values, v.Value)
	}

	return values
}
