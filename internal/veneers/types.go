package veneers

import (
	"github.com/grafana/cog/internal/ast"
)

type Option struct {
	Name        string         `yaml:"name"`
	Comments    []string       `yaml:"comments"`
	Arguments   []ast.Argument `yaml:"arguments"`
	Assignments []Assignment   `yaml:"assignments"`
}

func (opt Option) AsIR(builders ast.Builders, root ast.Builder) (ast.Option, error) {
	assignments := make([]ast.Assignment, 0, len(opt.Assignments))
	for _, assignment := range opt.Assignments {
		irAssignment, err := assignment.AsIR(builders, root)
		if err != nil {
			return ast.Option{}, err
		}

		assignments = append(assignments, irAssignment)
	}

	return ast.Option{
		Name:        opt.Name,
		Comments:    opt.Comments,
		Args:        opt.Arguments,
		Assignments: assignments,
	}, nil
}

type Assignment struct {
	Path   string               `yaml:"path"`
	Method ast.AssignmentMethod `yaml:"method"`
	Value  ast.AssignmentValue  `yaml:"value"`
}

func (assignment Assignment) AsIR(builders ast.Builders, root ast.Builder) (ast.Assignment, error) {
	// TODO
	path, err := root.MakePath(builders, assignment.Path)
	if err != nil {
		return ast.Assignment{}, err
	}

	return ast.Assignment{
		Path:   path,
		Value:  assignment.Value,
		Method: assignment.Method,
	}, nil
}
