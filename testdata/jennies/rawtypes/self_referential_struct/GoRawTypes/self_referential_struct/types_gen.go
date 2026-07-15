package self_referential_struct

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

// Node represents a node in a singly-linked list.
// The next field points to the following node, or is absent if this is the last node.
type Node struct {
    Value string `json:"value"`
    Next *Node `json:"next,omitempty"`
}

// NewNode creates a new Node object.
func NewNode() *Node {
	return &Node{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Node` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Node) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "value"
	if fields["value"] != nil {
		if string(fields["value"]) != "null" {
			if err := json.Unmarshal(fields["value"], &resource.Value); err != nil {
				errs = append(errs, cog.MakeBuildErrors("value", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("value", errors.New("required field is null"))...)
		
		}
		delete(fields, "value")
	} else {errs = append(errs, cog.MakeBuildErrors("value", errors.New("required field is missing from input"))...)
	}
	// Field "next"
	if fields["next"] != nil {
		if string(fields["next"]) != "null" {
			
			resource.Next = &Node{}
			if err := resource.Next.UnmarshalJSONStrict(fields["next"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("next", err)...)
			}
		
		}
		delete(fields, "next")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Node", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `Node` objects.
func (resource Node) Equals(other Node) bool {
		if resource.Value != other.Value {
			return false
		}
		if resource.Next == nil && other.Next != nil || resource.Next != nil && other.Next == nil {
			return false
		}

		if resource.Next != nil {
		if !resource.Next.Equals(*other.Next) {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `Node` fields for violations and returns them.
func (resource Node) Validate() error {
	var errs cog.BuildErrors
		if resource.Next != nil {
		if err := resource.Next.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("next", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


