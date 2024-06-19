package terraform

import (
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
)

type GoMod struct {
	Config Config
}

func (jenny GoMod) JennyName() string {
	return "GoMod"
}

func (jenny GoMod) Generate(_ common.Context) (codejen.Files, error) {
	return codejen.Files{
		*codejen.NewFile("go.mod", []byte(jenny.generateGoMod()), jenny),
	}, nil
}

func (jenny GoMod) generateGoMod() string {
	return fmt.Sprintf(`module %s

go 1.21

`, jenny.Config.PackageRoot)
}
