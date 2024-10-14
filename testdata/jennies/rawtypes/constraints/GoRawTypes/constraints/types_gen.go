package constraints

import (
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"strconv"
)

type SomeStruct struct {
	Id uint64 `json:"id"`
	MaybeId *uint64 `json:"maybeId,omitempty"`
	Title string `json:"title"`
	RefStruct *RefStruct `json:"refStruct,omitempty"`
}

func (resource SomeStruct) Equals(other SomeStruct) bool {
		if resource.Id != other.Id {
			return false
		}
		if resource.MaybeId == nil && other.MaybeId != nil || resource.MaybeId != nil && other.MaybeId == nil {
			return false
		}

		if resource.MaybeId != nil {
		if *resource.MaybeId != *other.MaybeId {
			return false
		}
		}
		if resource.Title != other.Title {
			return false
		}
		if resource.RefStruct == nil && other.RefStruct != nil || resource.RefStruct != nil && other.RefStruct == nil {
			return false
		}

		if resource.RefStruct != nil {
		if !resource.RefStruct.Equals(*other.RefStruct) {
			return false
		}
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource SomeStruct) Validate() error {
	var errs cog.BuildErrors
		if !(resource.Id >= 5) {
			errs = append(errs, cog.MakeBuildErrors(
				"id",
				errors.New("must be >= 5"),
			)...)
		}
		if !(resource.Id < 10) {
			errs = append(errs, cog.MakeBuildErrors(
				"id",
				errors.New("must be < 10"),
			)...)
		}
		if resource.MaybeId != nil {
		if !(*resource.MaybeId >= 5) {
			errs = append(errs, cog.MakeBuildErrors(
				"maybeId",
				errors.New("must be >= 5"),
			)...)
		}
		if !(*resource.MaybeId < 10) {
			errs = append(errs, cog.MakeBuildErrors(
				"maybeId",
				errors.New("must be < 10"),
			)...)
		}
		}
		if !(len([]rune(resource.Title)) >= 1) {
			errs = append(errs, cog.MakeBuildErrors(
				"title",
				errors.New("must be >= 1"),
			)...)
		}
		if resource.RefStruct != nil {
		if err := resource.RefStruct.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("refStruct", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type RefStruct struct {
	Labels map[string]string `json:"labels"`
	Tags []string `json:"tags"`
}

func (resource RefStruct) Equals(other RefStruct) bool {

		if len(resource.Labels) != len(other.Labels) {
			return false
		}

		for key1 := range resource.Labels {
		if resource.Labels[key1] != other.Labels[key1] {
			return false
		}
		}

		if len(resource.Tags) != len(other.Tags) {
			return false
		}

		for i1 := range resource.Tags {
		if resource.Tags[i1] != other.Tags[i1] {
			return false
		}
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource RefStruct) Validate() error {
	var errs cog.BuildErrors

		for key1 := range resource.Labels {
		if !(len([]rune(resource.Labels[key1])) >= 1) {
			errs = append(errs, cog.MakeBuildErrors(
				"labels["+key1+"]",
				errors.New("must be >= 1"),
			)...)
		}
		}

		for i1 := range resource.Tags {
		if !(len([]rune(resource.Tags[i1])) >= 1) {
			errs = append(errs, cog.MakeBuildErrors(
				"tags["+strconv.Itoa(i1)+"]",
				errors.New("must be >= 1"),
			)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

