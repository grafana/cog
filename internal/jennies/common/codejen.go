package common

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

func If[Input any](condition bool, innerJenny codejen.OneToMany[Input]) codejen.OneToMany[Input] {
	if !condition {
		return noopOneToManyJenny[Input]{}
	}

	return innerJenny
}
