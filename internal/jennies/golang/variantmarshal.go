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

	files = append(files, *codejen.NewFile("variants.go", jenny.variantRegistries(), jenny))

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

	buffer.WriteString(fmt.Sprintf(`cog.RegisterDataqueryVariant("%[1]s", cog.DataqueryVariantMarshalConfig{
`, strings.ToLower(schema.Metadata.Identifier)))

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
		break
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

func (jenny VariantMarshalConfig) variantRegistries() []byte {
	return []byte(`package cog

import (
	"encoding/json"
	"fmt"
)

type PanelcfgVariantMarshalConfig struct {
	OptionsUnmarshaler     func(raw []byte) (any, error)
	FieldConfigUnmarshaler func(raw []byte) (any, error)
}

var panelcfgVariantsRegistry map[string]PanelcfgVariantMarshalConfig

func RegisterPanelcfgVariant(variantName string, config PanelcfgVariantMarshalConfig) {
	if panelcfgVariantsRegistry == nil {
		panelcfgVariantsRegistry = make(map[string]PanelcfgVariantMarshalConfig)
	}
	panelcfgVariantsRegistry[variantName] = config
}

func MarshalConfigForPanelcfgVariant(variantName string) (PanelcfgVariantMarshalConfig, bool) {
	config, found := panelcfgVariantsRegistry[variantName]

	return config, found
}

type DataqueryVariantMarshalConfig struct {
	DataqueryUnmarshaler func(raw []byte) (interface {
	ImplementsDataqueryVariant()
}, error)
}

var dataqueryVariantsRegistry map[string]DataqueryVariantMarshalConfig

func RegisterDataqueryVariant(variantName string, config DataqueryVariantMarshalConfig) {
	if dataqueryVariantsRegistry == nil {
		dataqueryVariantsRegistry = make(map[string]DataqueryVariantMarshalConfig)
	}
	dataqueryVariantsRegistry[variantName] = config
}

func UnmarshalDataqueryArray[VariantT interface {
	ImplementsDataqueryVariant()
}](raw []byte) ([]VariantT, error) {
	rawDataqueries := []json.RawMessage{}
	if err := json.Unmarshal(raw, &rawDataqueries); err != nil {
		return nil, err
	}

	dataqueries := make([]VariantT, 0, len(rawDataqueries))
	for _, rawDataquery := range rawDataqueries {
		dataquery, err := UnmarshalDataquery[VariantT](rawDataquery)
		if err != nil {
			return nil, err
		}

		dataqueries = append(dataqueries, dataquery)
	}

	return dataqueries, nil
}

func UnmarshalDataquery[VariantT any](raw []byte) (VariantT, error) {
	var empty VariantT

	for _, config := range dataqueryVariantsRegistry {
		if config.DataqueryUnmarshaler == nil {
			continue
		}

		dataquery, err := config.DataqueryUnmarshaler(raw)
		if err != nil {
			continue
		}

		return dataquery.(VariantT), nil
	}

	return empty, fmt.Errorf("could not unmarshal dataquery")
}

`)
}
