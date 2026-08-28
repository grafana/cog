package veneers

import (
	"fmt"

	"github.com/grafana/cog/pkg/ir"
)

type Option struct {
	Name        string        `yaml:"name"`
	Comments    []string      `yaml:"comments"`
	Arguments   []ir.Argument `yaml:"arguments"`
	Assignments []Assignment  `yaml:"assignments"`
}

func (opt Option) AsIR(schemas ir.Schemas, root ir.Builder) (ir.Option, error) {
	assignments := make([]ir.Assignment, 0, len(opt.Assignments))
	for _, assignment := range opt.Assignments {
		irAssignment, err := assignment.AsIR(schemas, root)
		if err != nil {
			return ir.Option{}, err
		}

		assignments = append(assignments, irAssignment)
	}

	return ir.Option{
		Name:        opt.Name,
		Comments:    opt.Comments,
		Args:        opt.Arguments,
		Assignments: assignments,
	}, nil
}

type Assignment struct {
	Path    string              `yaml:"path"`
	PathObj ir.Path             `yaml:"path_obj"`
	Method  ir.AssignmentMethod `yaml:"method"`
	Value   AssignmentValue     `yaml:"value"`
}

func (assignment Assignment) AsIR(schemas ir.Schemas, root ir.Builder) (ir.Assignment, error) {
	var path ir.Path
	var err error

	if assignment.PathObj != nil {
		path = assignment.PathObj
	} else {
		path, err = root.MakePath(schemas, assignment.Path)
		if err != nil {
			return ir.Assignment{}, err
		}
	}

	value, err := assignment.Value.AsIR(schemas, path)
	if err != nil {
		return ir.Assignment{}, err
	}

	return ir.Assignment{
		Path:   path,
		Value:  value,
		Method: assignment.Method,
	}, nil
}

type AssignmentValue struct {
	Argument *ir.Argument `json:",omitempty"`
	Constant any
	Envelope *AssignmentEnvelope `json:",omitempty"`
}

func (value AssignmentValue) AsIR(schemas ir.Schemas, assignmentPath ir.Path) (ir.AssignmentValue, error) {
	if value.Argument != nil {
		return ir.AssignmentValue{Argument: value.Argument}, nil
	}
	if value.Envelope != nil {
		envelopeType := assignmentPath.Last().Type
		if envelopeType.IsArray() {
			envelopeType = envelopeType.Array.ValueType
		}
		if envelopeType.IsMap() {
			envelopeType = envelopeType.Map.ValueType
		}

		envelope, err := value.Envelope.AsIR(schemas, envelopeType)
		if err != nil {
			return ir.AssignmentValue{}, err
		}

		return ir.AssignmentValue{Envelope: &envelope}, nil
	}

	return ir.AssignmentValue{Constant: value.Constant}, nil
}

type AssignmentEnvelope struct {
	Values []EnvelopeFieldValue
}

func (envelope AssignmentEnvelope) AsIR(schemas ir.Schemas, envelopeType ir.Type) (ir.AssignmentEnvelope, error) {
	var err error
	values := make([]ir.EnvelopeFieldValue, len(envelope.Values))
	for i, val := range envelope.Values {
		values[i], err = val.AsIR(schemas, envelopeType)
		if err != nil {
			return ir.AssignmentEnvelope{}, err
		}
	}

	return ir.AssignmentEnvelope{
		Type:   envelopeType,
		Values: values,
	}, nil
}

type EnvelopeFieldValue struct {
	Field string          // where to assign within the struct/ref
	Value AssignmentValue // what to assign
}

func (envelopeField EnvelopeFieldValue) AsIR(schemas ir.Schemas, envelopeType ir.Type) (ir.EnvelopeFieldValue, error) {
	resolvedEnvelope := schemas.Resolve(envelopeType)
	field, found := resolvedEnvelope.Struct.FieldByName(envelopeField.Field)
	if !found {
		return ir.EnvelopeFieldValue{}, fmt.Errorf("envelope field %s not found", envelopeField.Field)
	}

	path := ir.PathFromStructField(field)

	value, err := envelopeField.Value.AsIR(schemas, path)
	if err != nil {
		return ir.EnvelopeFieldValue{}, err
	}

	return ir.EnvelopeFieldValue{
		Path:  path,
		Value: value,
	}, nil
}
