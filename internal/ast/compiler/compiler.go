package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

type Pass interface {
	Process(schemas []*ast.Schema) ([]*ast.Schema, error)
}

func CommonPasses() []Pass {
	return []Pass{
		&Unspec{},
		&DashboardPanelsRewrite{},
		&DashboardTargetsRewrite{},
		&DashboardTimePicker{},
		&Cloudwatch{},
	}
}
