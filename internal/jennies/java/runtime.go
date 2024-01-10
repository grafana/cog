package java

import (
	"bytes"
	"fmt"
	
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
)

type Runtime struct {
}

func (jenny Runtime) JennyName() string {
	return "JavaRuntime"
}

func (jenny Runtime) Generate(_ common.Context) (codejen.Files, error) {
	variants, err := jenny.renderDataQueryVariant("Dataquery")
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile("cog/variants/Dataquery.java", variants, jenny),
	}, nil
}

func (jenny Runtime) renderDataQueryVariant(variant string) ([]byte, error) {
	buf := bytes.Buffer{}
	if err := templates.ExecuteTemplate(&buf, "runtime/variants.tmpl", map[string]any{
		"Variant": variant,
	}); err != nil {
		return nil, fmt.Errorf("failed executing template: %w", err)
	}

	return buf.Bytes(), nil
}
