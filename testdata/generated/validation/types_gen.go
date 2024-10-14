// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Using jennies:
//     GoRawTypes

package validation

import (
	"errors"
	"strconv"

	cog "github.com/grafana/cog/testdata/generated/cog"
)

type Dashboard struct {
	// Modified by compiler pass 'NotRequiredFieldAsNullableType[nullable=true]'
	Uid *string `json:"uid,omitempty"`
	// Modified by compiler pass 'NotRequiredFieldAsNullableType[nullable=true]'
	Id     *int64            `json:"id,omitempty"`
	Title  string            `json:"title"`
	Tags   []string          `json:"tags"`
	Labels map[string]string `json:"labels"`
	Panels []Panel           `json:"panels"`
}

func (resource Dashboard) Equals(other Dashboard) bool {
	if resource.Uid == nil && other.Uid != nil || resource.Uid != nil && other.Uid == nil {
		return false
	}

	if resource.Uid != nil {
		if *resource.Uid != *other.Uid {
			return false
		}
	}
	if resource.Id == nil && other.Id != nil || resource.Id != nil && other.Id == nil {
		return false
	}

	if resource.Id != nil {
		if *resource.Id != *other.Id {
			return false
		}
	}
	if resource.Title != other.Title {
		return false
	}

	if len(resource.Tags) != len(other.Tags) {
		return false
	}

	for i1 := range resource.Tags {
		if resource.Tags[i1] != other.Tags[i1] {
			return false
		}
	}

	if len(resource.Labels) != len(other.Labels) {
		return false
	}

	for key1 := range resource.Labels {
		if resource.Labels[key1] != other.Labels[key1] {
			return false
		}
	}

	if len(resource.Panels) != len(other.Panels) {
		return false
	}

	for i1 := range resource.Panels {
		if !resource.Panels[i1].Equals(other.Panels[i1]) {
			return false
		}
	}

	return true
}

// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource Dashboard) Validate() error {
	var errs cog.BuildErrors
	if resource.Uid != nil {
		if !(len([]rune(*resource.Uid)) >= 1) {
			errs = append(errs, cog.MakeBuildErrors(
				"uid",
				errors.New("must be >= 1"),
			)...)
		}
	}
	if resource.Id != nil {
		if !(*resource.Id > 0) {
			errs = append(errs, cog.MakeBuildErrors(
				"id",
				errors.New("must be > 0"),
			)...)
		}
	}
	if !(len([]rune(resource.Title)) >= 1) {
		errs = append(errs, cog.MakeBuildErrors(
			"title",
			errors.New("must be >= 1"),
		)...)
	}

	for i1 := range resource.Tags {
		if !(len([]rune(resource.Tags[i1])) >= 1) {
			errs = append(errs, cog.MakeBuildErrors(
				"tags["+strconv.Itoa(i1)+"]",
				errors.New("must be >= 1"),
			)...)
		}
	}

	for key1 := range resource.Labels {
		if !(len([]rune(resource.Labels[key1])) >= 1) {
			errs = append(errs, cog.MakeBuildErrors(
				"labels["+key1+"]",
				errors.New("must be >= 1"),
			)...)
		}
	}

	for i1 := range resource.Panels {
		if err := resource.Panels[i1].Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("panels["+strconv.Itoa(i1)+"]", err)...)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

type Panel struct {
	Title string `json:"title"`
}

func (resource Panel) Equals(other Panel) bool {
	if resource.Title != other.Title {
		return false
	}

	return true
}

// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource Panel) Validate() error {
	var errs cog.BuildErrors
	if !(len([]rune(resource.Title)) >= 1) {
		errs = append(errs, cog.MakeBuildErrors(
			"title",
			errors.New("must be >= 1"),
		)...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
