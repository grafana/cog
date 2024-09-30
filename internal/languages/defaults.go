package languages

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

type Defaults struct {
	Name string
}

func GenerateDefaults(language Language, context Context) (Context, error) {
	defaultConfig := DefaultConfig{}

	if languageConfig, ok := language.(DefaultKindsProvider); ok {
		defaultConfig = languageConfig.DefaultKinds()
	}

	defaultVisitor := ast.BuilderVisitor{
		OnOption: func(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, option ast.Option) (ast.Option, error) {
			if option.Default == nil {
				return option, nil
			}

			if len(option.Args) == 0 {
				return option, nil
			}

			for _, arg := range option.Args {
				switch arg.Type.Kind {
				case ast.KindScalar:
					if defaultConfig.FormatScalarFunc != nil {
						option.Default.ArgsValues = []any{defaultConfig.FormatScalarFunc(arg.Type, arg.Type.Default)}
					}
				case ast.KindArray:
					if defaultConfig.FormatListFunc != nil {
						option.Default.ArgsValues = []any{defaultConfig.FormatListFunc(arg.Type, arg.Type.Default)}
					}
				case ast.KindRef:
					ref := arg.Type.AsRef()
					obj, ok := schemas.LocateObject(ref.ReferredPkg, ref.ReferredType)
					if !ok {
						return ast.Option{}, fmt.Errorf("cannot locate %s.%s object", ref.ReferredPkg, ref.ReferredType)
					}

					if obj.Type.Kind == ast.KindEnum {
						if defaultConfig.FormatEnumFunc != nil {
							option.Default.ArgsValues = []any{defaultConfig.FormatEnumFunc(obj.Name, obj.Type, arg.Type.Default)}
						}
					} else if defaultConfig.FormatStructFunc != nil {
						option.Default.ArgsValues = []any{defaultConfig.FormatStructFunc(obj.Name, obj.Type, arg.Type.Default)}
					}
				case ast.KindStruct:
					if defaultConfig.FormatStructFunc != nil {
						option.Default.ArgsValues = []any{defaultConfig.FormatStructFunc(builder.Name, arg.Type, arg.Type.Default)}
					}
				case ast.KindMap:
					if defaultConfig.FormatMapFunc != nil {
						option.Default.ArgsValues = []any{defaultConfig.FormatMapFunc(arg.Type, arg.Type.Default)}
					}
				}
			}

			return option, nil
		},
	}

	var err error
	context.Builders, err = defaultVisitor.Visit(context.Schemas, context.Builders)
	if err != nil {
		return Context{}, err
	}

	return context, nil
}
