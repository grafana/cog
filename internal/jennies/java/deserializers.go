package java

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Deserializers struct {
	config  Config
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
			if jenny.objectNeedsCustomDeserialiser(context, obj) {
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
	buf := bytes.Buffer{}

	jenny.imports = jenny.genImports(obj)
	dataqueryCode, err := jenny.genDataqueryCode(context, obj)
	if err != nil {
		return nil, err
	}

	if err := templates.ExecuteTemplate(&buf, "marshalling/deserialiser.tmpl", Unmarshalling{
		Package:                   jenny.formatPackage(obj.SelfRef.ReferredPkg),
		Name:                      obj.Name,
		ShouldUnmarshallingPanels: obj.SelfRef.ReferredPkg == "dashboard" && obj.Name == "Panel",
		Imports:                   jenny.imports,
		DataqueryUnmarshalling:    dataqueryCode,
	}); err != nil {
		return nil, fmt.Errorf("failed executing template: %w", err)
	}

	path := filepath.Join(jenny.config.ProjectPath, obj.SelfRef.ReferredPkg, fmt.Sprintf("%sDeserializer.java", obj.SelfRef.ReferredType))
	return codejen.NewFile(path, buf.Bytes(), jenny), nil
}

func (jenny *Deserializers) objectNeedsCustomDeserialiser(context languages.Context, obj ast.Object) bool {
	// an object needs a custom unmarshal if:
	// - it is a struct that was generated from a disjunction by the `DisjunctionToType` compiler pass.
	// - it is a struct and one or more of its fields is a KindComposableSlot, or an array of KindComposableSlot

	if !obj.Type.IsStruct() {
		return false
	}

	if obj.Type.HasHint(ast.HintDisjunctionOfScalars) || obj.Type.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) {
		return false
	}

	// is it a struct generated from a disjunction?
	if obj.Type.IsStructGeneratedFromDisjunction() {
		return true
	}

	// is there a KindComposableSlot field somewhere?
	for _, field := range obj.Type.AsStruct().Fields {
		if _, ok := context.ResolveToComposableSlot(field.Type); ok {
			return true
		}
	}

	return false
}

func (jenny *Deserializers) genDataqueryCode(context languages.Context, obj ast.Object) (DataqueryUnmarshalling, error) {
	for _, field := range obj.Type.AsStruct().Fields {
		composableSlotType, resolved := context.ResolveToComposableSlot(field.Type)
		if !resolved {
			continue
		}

		if composableSlotType.AsComposableSlot().Variant == ast.SchemaVariantDataQuery {
			return jenny.renderUnmarshalDataqueryField(obj, field), nil
		}
	}

	return DataqueryUnmarshalling{}, errors.New("cannot find dataquery reference")
}

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
			jenny.imports = append(jenny.imports, jenny.formatPackage(fmt.Sprintf("%s.%s", f.Type.AsRef().ReferredPkg, "DataSourceRef")))
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

func (jenny *Deserializers) genImports(obj ast.Object) []string {
	imports := []string{
		jenny.formatPackage("cog.variants.Dataquery"),
		jenny.formatPackage("cog.variants.Registry"),
	}

	if obj.SelfRef.ReferredPkg == "dashboard" && obj.Name == "Panel" {
		imports = append(imports, jenny.formatPackage("cog.variants.PanelConfig"))
	}

	return imports
}

func (jenny *Deserializers) formatPackage(pkg string) string {
	if jenny.config.PackagePath != "" {
		return fmt.Sprintf("%s.%s", jenny.config.PackagePath, pkg)
	}

	return pkg
}
