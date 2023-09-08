package tools

import (
	"github.com/grafana/codejen"
)

type foreach[Input any] struct {
	inner codejen.OneToMany[Input]
}

func (jenny foreach[Input]) JennyName() string {
	return "ForeachFile"
}

func (jenny foreach[Input]) Generate(inputs ...[]Input) (codejen.Files, error) {
	outputs := make([]codejen.File, 0, len(inputs))

	for _, input := range inputs {
		for _, item := range input {
			out, err := jenny.inner.Generate(item)
			if err != nil {
				return nil, err
			}

			outputs = append(outputs, out...)

		}
	}

	return outputs, nil
}

func Foreach[InputInner any](decoratedJenny codejen.OneToMany[InputInner]) codejen.ManyToMany[[]InputInner] {
	return foreach[InputInner]{
		inner: decoratedJenny,
	}
}
