package java

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Deserializers struct {
	config  Config
	tmpl    *template.Template
	imports []string
}

func (jenny *Deserializers) JennyName() string {
	return "JavaDeserializers"
}

func (jenny *Deserializers) Generate(context languages.Context) (codejen.Files, error) {
	deserialisers := make(codejen.Files, 0)
	for _, schema := range context.Schemas {
		var hasErr error
		schema.Objects.Iterate(func(key string, obj ast.Object) {
			if objectNeedsCustomDeserialiser(context, obj) {
				f, err := jenny.genCustomDeserialiser(context, obj)
				if err != nil {
					hasErr = err
				} else {
					deserialisers = append(deserialisers, *f)
				}
			}
		})
		if hasErr != nil {
			return nil, hasErr
		}
	}

	return deserialisers, nil
}

func (jenny *Deserializers) genCustomDeserialiser(context languages.Context, obj ast.Object) (*codejen.File, error) {
	if obj.Type.IsStruct() && obj.Type.HasHint(ast.HintDisjunctionOfScalars) {
		return jenny.genDisjunctionsDeserialiser(obj, "disjunctions_of_scalars")
	}

	if obj.Type.IsStruct() && obj.Type.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) {
		return jenny.genDisjunctionsDeserialiser(obj, "disjunctions_of_refs")
	}

	// TODO(kgz): this shouldn't be done by cog
	return jenny.genDataqueryDeserialiser(context, obj)
}

// TODO(kgz): this shouldn't be done by cog
func (jenny *Deserializers) genDataqueryDeserialiser(context languages.Context, obj ast.Object) (*codejen.File, error) {
	jenny.imports = jenny.genImports(obj)

	rendered, err := jenny.tmpl.Render("marshalling/unmarshalling.tmpl", Unmarshalling{
		Package:                   jenny.config.formatPackage(obj.SelfRef.ReferredPkg),
		Name:                      obj.Name,
		ShouldUnmarshallingPanels: obj.SelfRef.ReferredPkg == "dashboard" && obj.Name == "Panel",
		Imports:                   jenny.imports,
		Fields:                    obj.Type.AsStruct().Fields,
		DataqueryUnmarshalling:    jenny.genDataqueryCode(context, obj),
	})
	if err != nil {
		return nil, err
	}

	path := filepath.Join(jenny.config.ProjectPath, obj.SelfRef.ReferredPkg, fmt.Sprintf("%sDeserializer.java", obj.SelfRef.ReferredType))
	return codejen.NewFile(path, []byte(rendered), jenny), nil
}

// TODO(kgz): this shouldn't be done by cog
func (jenny *Deserializers) genDataqueryCode(context languages.Context, obj ast.Object) []DataqueryUnmarshalling {
	dataqueryUnmarshalling := make([]DataqueryUnmarshalling, 0)
	for _, field := range obj.Type.AsStruct().Fields {
		composableSlotType, resolved := context.ResolveToComposableSlot(field.Type)
		if !resolved {
			continue
		}

		if composableSlotType.AsComposableSlot().Variant == ast.SchemaVariantDataQuery {
			dataqueryUnmarshalling = append(dataqueryUnmarshalling, jenny.renderUnmarshalDataqueryField(obj, field))
		}
	}

	return dataqueryUnmarshalling
}

// TODO(kgz): this shouldn't be done by cog
func (jenny *Deserializers) renderUnmarshalDataqueryField(obj ast.Object, field ast.StructField) DataqueryUnmarshalling {
	var hintField *ast.StructField
	for i, f := range obj.Type.AsStruct().Fields {
		if !f.Type.IsRef() {
			continue
		}

		if f.Type.AsRef().ReferredType != "DataSourceRef" {
			continue
		}

		hintField = &obj.Type.AsStruct().Fields[i]
		if obj.SelfRef.ReferredPkg != f.Type.AsRef().ReferredPkg {
			jenny.imports = append(jenny.imports, jenny.config.formatPackage(fmt.Sprintf("%s.%s", f.Type.AsRef().ReferredPkg, "DataSourceRef")))
		}
	}

	dataqueryHint := `""`
	hintFieldName := ""
	if hintField != nil {
		hintFieldName = hintField.Name
		dataqueryHint = fmt.Sprintf("%s.datasource.type", tools.LowerCamelCase(obj.Name))
	}

	return DataqueryUnmarshalling{
		DataqueryHint:   dataqueryHint,
		IsArray:         field.Type.IsArray(),
		FieldName:       field.Name,
		DatasourceField: hintFieldName,
	}
}

// TODO(kgz): this shouldn't be done by cog
func (jenny *Deserializers) genImports(obj ast.Object) []string {
	imports := []string{
		jenny.config.formatPackage("cog.variants.Dataquery"),
		jenny.config.formatPackage("cog.variants.Registry"),
	}

	if obj.SelfRef.ReferredPkg == "dashboard" && obj.Name == "Panel" {
		imports = append(imports, jenny.config.formatPackage("cog.variants.PanelConfig"))
	}

	return imports
}

func (jenny *Deserializers) genDisjunctionsDeserialiser(obj ast.Object, tmpl string) (*codejen.File, error) {
	rendered, err := jenny.tmpl.Render(fmt.Sprintf("marshalling/%s.json_unmarshall.tmpl", tmpl), Unmarshalling{
		Package: jenny.config.formatPackage(obj.SelfRef.ReferredPkg),
		Name:    tools.UpperCamelCase(obj.Name),
		Fields:  obj.Type.AsStruct().Fields,
		Hint:    obj.Type.Hints[ast.HintDiscriminatedDisjunctionOfRefs],
	})
	if err != nil {
		return nil, err
	}

	path := filepath.Join(jenny.config.ProjectPath, obj.SelfRef.ReferredPkg, fmt.Sprintf("%sDeserializer.java", tools.UpperCamelCase(obj.SelfRef.ReferredType)))
	return codejen.NewFile(path, []byte(rendered), jenny), nil
}
