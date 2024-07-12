package java

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Serializers struct {
	config Config
}

func (jenny *Serializers) JennyName() string {
	return "JavaDeserializers"
}

func (jenny *Serializers) Generate(context languages.Context) (codejen.Files, error) {
	serializers := make(codejen.Files, 0)
	for _, schema := range context.Schemas {
		var hasErr error
		schema.Objects.Iterate(func(key string, obj ast.Object) {
			if obj.Type.HasHint(ast.HintDisjunctionOfScalars) {
				f, err := jenny.genSerializer(obj)
				if err != nil {
					hasErr = err
				} else {
					serializers = append(serializers, *f)
				}
			}
		})
		if hasErr != nil {
			return nil, hasErr
		}
	}

	return serializers, nil
}

func (jenny *Serializers) genSerializer(obj ast.Object) (*codejen.File, error) {
	buf := bytes.Buffer{}

	if err := templates.ExecuteTemplate(&buf, "marshalling/disjunctions_of_scalars.json_marshall.tmpl", Unmarshalling{
		Package: jenny.formatPackage(obj.SelfRef.ReferredPkg),
		Name:    tools.UpperCamelCase(obj.Name),
		Fields:  obj.Type.AsStruct().Fields,
	}); err != nil {
		return nil, err
	}

	path := filepath.Join(jenny.config.ProjectPath, obj.SelfRef.ReferredPkg, fmt.Sprintf("%sSerializer.java", tools.UpperCamelCase(obj.SelfRef.ReferredType)))
	return codejen.NewFile(path, buf.Bytes(), jenny), nil
}

func (jenny *Serializers) formatPackage(pkg string) string {
	if jenny.config.PackagePath != "" {
		return fmt.Sprintf("%s.%s", jenny.config.PackagePath, pkg)
	}

	return pkg
}
