package java

import (
	"fmt"
	gotemplate "text/template"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type Gradle struct {
	config Config
	tmpl   *gotemplate.Template
}

func (jenny Gradle) JennyName() string {
	return "JavaGradle"
}

func (jenny Gradle) Generate(_ languages.Context) (codejen.Files, error) {
	bg, err := jenny.gen("build.gradle")
	if err != nil {
		return nil, err
	}

	gp, err := jenny.gen("gradle.properties")
	if err != nil {
		return nil, err
	}

	sg, err := jenny.gen("settings.gradle")
	if err != nil {
		return nil, err
	}

	return codejen.Files{*bg, *gp, *sg}, nil
}

func (jenny Gradle) gen(tmpl string) (*codejen.File, error) {
	rendered, err := template.Render(jenny.tmpl, fmt.Sprintf("gradle/%s", tmpl), map[string]string{})
	return codejen.NewFile(tmpl, []byte(rendered), jenny), err
}
