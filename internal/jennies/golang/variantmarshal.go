package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/tools"
)

type VariantMarshalConfig struct {
}

func (jenny VariantMarshalConfig) JennyName() string {
	return "GoVariantMarshalConfig"
}

func (jenny VariantMarshalConfig) Generate(context context.Builders) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		if schema.Metadata.Kind != ast.SchemaKindComposable {
			continue
		}

		output := jenny.generateMarshalConfigRegistration(schema)
		if len(output) == 0 {
			continue
		}

		filename := filepath.Join(
			strings.ToLower(schema.Package),
			"init.go",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny VariantMarshalConfig) generateMarshalConfigRegistration(schema *ast.Schema) []byte {
	switch schema.Metadata.Variant {
	case ast.SchemaVariantPanel:
		return jenny.panelcfgVariantInit(schema)
	case ast.SchemaVariantDataQuery:
		return jenny.dataqueryVariantInit(schema)
	}

	return nil
}

func (jenny VariantMarshalConfig) dataqueryVariantInit(schema *ast.Schema) []byte {
	var buffer strings.Builder
	imports := newImportMap()
	imports.Add("cog", "github.com/grafana/cog/generated")

	buffer.WriteString("func init() {\n")

	buffer.WriteString(`cog.RegisterDataqueryVariant(cog.DataqueryVariantMarshalConfig{
`)

	variantObjFound := false
	for _, obj := range schema.Objects {
		if !obj.Type.ImplementsVariant() || obj.Type.ImplementedVariant() != string(ast.SchemaVariantDataQuery) {
			continue
		}

		variantObjFound = true
		buffer.WriteString(fmt.Sprintf(`DataqueryUnmarshaler: func(raw []byte) (interface {
	ImplementsDataqueryVariant()
}, error) {
	dataquery := %s{}

	if err := json.Unmarshal(raw, &dataquery); err != nil {
		return nil, err
	}

	return dataquery, nil
},
`, tools.UpperCamelCase(obj.Name)))
	}

	if !variantObjFound {
		return nil
	}

	buffer.WriteString("})\n")
	buffer.WriteString("}\n")

	importStatements := imports.Format()
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(fmt.Sprintf(`package %[1]s

%[2]s%[3]s`, strings.ToLower(schema.Package), importStatements, buffer.String()))
}

func (jenny VariantMarshalConfig) panelcfgVariantInit(schema *ast.Schema) []byte {
	var buffer strings.Builder
	imports := newImportMap()
	imports.Add("cog", "github.com/grafana/cog/generated")

	buffer.WriteString("func init() {\n")

	buffer.WriteString(fmt.Sprintf(`cog.RegisterPanelcfgVariant("%[1]s", cog.PanelcfgVariantMarshalConfig{
`, strings.ToLower(schema.Metadata.Identifier)))

	optionsObj := schema.LocateObject("Options")
	if optionsObj.Name != "" {
		buffer.WriteString(`OptionsUnmarshaler: func(raw []byte) (any, error) {
	options := Options{}

	if err := json.Unmarshal(raw, &options); err != nil {
		return nil, err
	}

	return options, nil
},
`)
	}

	fieldConfigObj := schema.LocateObject("FieldConfig")
	if fieldConfigObj.Name != "" {
		buffer.WriteString(`FieldConfigUnmarshaler: func(raw []byte) (any, error) {
	fieldConfig := FieldConfig{}

	if err := json.Unmarshal(raw, &fieldConfig); err != nil {
		return nil, err
	}

	return fieldConfig, nil
},
`)
	}

	buffer.WriteString("})\n")
	buffer.WriteString("}\n")

	importStatements := imports.Format()
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(fmt.Sprintf(`package %[1]s

%[2]s%[3]s`, strings.ToLower(schema.Package), importStatements, buffer.String()))
}
