package simplecue

import (
	"fmt"
	"path"
	"strings"

	"cuelang.org/go/cue"
	cueast "cuelang.org/go/cue/ast"
)

type referenceResolverConfig struct {
	SchemaPackage string
	SourcePackage string
	Libraries     []LibraryInclude
}

type referenceResolver struct {
	config referenceResolverConfig

	rootFileName    string
	importsAliasMap map[string]string // map[alias]packageName
	librariesMap    map[string]string // map[libraryImportPath]packageName
}

func newReferenceResolver(root cue.Value, config referenceResolverConfig) *referenceResolver {
	resolver := &referenceResolver{
		config:       config,
		rootFileName: getFilename(root),
	}

	resolver.importsAliasMap = resolver.buildImportsAliasMap(root)
	resolver.librariesMap = resolver.buildLibrariesMap(config.Libraries)

	if config.SourcePackage != config.SchemaPackage {
		resolver.importsAliasMap[config.SourcePackage] = config.SchemaPackage
	}

	// map the package originally declared in the CUE file to the one configured by the user
	if fileAst, ok := root.Source().(*cueast.File); ok {
		for _, decl := range fileAst.Decls {
			if pkgDecl, ok := decl.(*cueast.Package); ok {
				resolver.importsAliasMap[pkgDecl.Name.String()] = config.SchemaPackage
			}
		}
	}

	return resolver
}

// FIXME: this is probably very brittle and not always correct :|
func (resolver *referenceResolver) PackageForNode(source cueast.Node) (string, error) {
	switch source.(type) { //nolint: gocritic
	case *cueast.SelectorExpr:
		selector := source.(*cueast.SelectorExpr)

		x := selector.X.(*cueast.Ident)

		return resolver.resolveImportAlias(x.Name), nil
	case *cueast.Field:
		field := source.(*cueast.Field)

		if _, ok := field.Value.(*cueast.SelectorExpr); ok {
			return resolver.PackageForNode(field.Value)
		}

		ident := field.Value.(*cueast.Ident)

		if ident.Scope == nil {
			return resolver.config.SchemaPackage, nil
		}

		if _, ok := ident.Scope.(*cueast.File); !ok {
			return resolver.config.SchemaPackage, nil
		}

		scope := ident.Scope.(*cueast.File)
		if path.Dir(scope.Filename) == path.Dir(resolver.rootFileName) {
			return resolver.config.SchemaPackage, nil
		}

		if len(scope.Decls) == 0 {
			return resolver.config.SchemaPackage, nil
		}

		if _, ok := scope.Decls[0].(*cueast.Package); !ok {
			return resolver.config.SchemaPackage, nil
		}

		referredTypePkg := scope.Decls[0].(*cueast.Package).Name

		return resolver.resolveImportAlias(referredTypePkg.Name), nil
	case *cueast.Ident:
		ident := source.(*cueast.Ident)
		if ident.Scope == nil {
			return resolver.config.SchemaPackage, nil
		}

		if _, ok := ident.Scope.(*cueast.File); !ok {
			return resolver.config.SchemaPackage, nil
		}

		scope := ident.Scope.(*cueast.File)
		if len(scope.Decls) == 0 {
			return resolver.config.SchemaPackage, nil
		}

		if _, ok := ident.Scope.(*cueast.File).Decls[0].(*cueast.Package); !ok {
			return resolver.config.SchemaPackage, nil
		}

		referredTypePkg := ident.Scope.(*cueast.File).Decls[0].(*cueast.Package).Name

		return resolver.resolveImportAlias(referredTypePkg.Name), nil
	case *cueast.Ellipsis:
		if source.(*cueast.Ellipsis).Type == nil {
			return resolver.config.SchemaPackage, nil
		}

		if _, ok := source.(*cueast.Ellipsis).Type.(*cueast.SelectorExpr); ok {
			return resolver.PackageForNode(source.(*cueast.Ellipsis).Type)
		}

		return resolver.packageForToken(source), nil
	default:
		return "", fmt.Errorf("can't extract reference package")
	}
}

func (resolver *referenceResolver) packageForToken(source cueast.Node) string {
	filename := source.Pos().Filename()
	if filename == "" {
		return resolver.config.SchemaPackage
	}

	for importPath, pkg := range resolver.librariesMap {
		if strings.Contains(filename, importPath) {
			return pkg
		}
	}

	return resolver.config.SchemaPackage
}

func (resolver *referenceResolver) resolveImportAlias(alias string) string {
	if resolved, found := resolver.importsAliasMap[alias]; found {
		return resolved
	}

	return alias
}

func (resolver *referenceResolver) buildImportsAliasMap(v cue.Value) map[string]string {
	aliasMap := make(map[string]string)
	syntax := v.Syntax()

	if _, ok := syntax.(*cueast.File); !ok {
		return aliasMap
	}

	file := syntax.(*cueast.File)

	for _, decl := range file.Decls {
		if _, ok := decl.(*cueast.ImportDecl); !ok {
			continue
		}

		importDecl := decl.(*cueast.ImportDecl)

		for _, spec := range importDecl.Specs {
			pkgName := resolver.packageFromImportPath(spec.Path.Value)

			alias := pkgName
			if spec.Name != nil {
				alias = spec.Name.Name
			}

			aliasMap[alias] = pkgName
		}
	}

	return aliasMap
}

func (resolver *referenceResolver) buildLibrariesMap(libraries []LibraryInclude) map[string]string {
	librariesMap := make(map[string]string)

	for _, library := range libraries {
		path := "cue.mod/pkg/" + library.ImportPath
		pkg := resolver.packageFromImportPath(path)

		librariesMap[path] = pkg
	}

	return librariesMap
}

func (resolver *referenceResolver) packageFromImportPath(importPath string) string {
	importPath = strings.Trim(importPath, "\"")
	parts := strings.Split(importPath, "/")

	return parts[len(parts)-1]
}
