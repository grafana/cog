package jennies

import (
	"github.com/grafana/codejen"
)

type noopOneToManyJenny[Input any] struct {
}

func (jenny noopOneToManyJenny[Input]) JennyName() string {
	return "noopOneToManyJenny"
}

func (jenny noopOneToManyJenny[Input]) Generate(_ Input) (codejen.Files, error) {
	return nil, nil
}

// If runs the inner jenny when condition is true.
func If[Input any](condition bool, inner codejen.OneToMany[Input]) codejen.OneToMany[Input] {
	if !condition {
		return noopOneToManyJenny[Input]{}
	}

	return inner
}
