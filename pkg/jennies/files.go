package jennies

import (
	"encoding/gob"

	"github.com/grafana/codejen"
)

func init() { //nolint:gochecknoinits
	gob.Register(unnamed)
}

var unnamed = &unNamedJenny{} //nolint:gochecknoglobals

type unNamedJenny struct {
}

func (jenny *unNamedJenny) JennyName() string {
	return ""
}

// NewFile creates a [codejen.File] so that it can be used in the context of
// language plugins.
func NewFile(path string, data []byte) *codejen.File {
	return codejen.NewFile(path, data, unnamed)
}
