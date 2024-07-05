package java

import (
	"bytes"
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/languages"
)

type Gradle struct {
	config Config
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
	buf := new(bytes.Buffer)
	err := templates.ExecuteTemplate(buf, fmt.Sprintf("gradle/%s", tmpl), map[string]string{})
	return codejen.NewFile(tmpl, buf.Bytes(), jenny), err
}
