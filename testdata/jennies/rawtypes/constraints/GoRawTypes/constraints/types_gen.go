package constraints

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"strconv"
)

type SomeStruct struct {
    Id uint64 `json:"id"`
    MaybeId *uint64 `json:"maybeId,omitempty"`
    Title string `json:"title"`
    RefStruct *RefStruct `json:"refStruct,omitempty"`
}

// NewSomeStruct creates a new SomeStruct object.
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `SomeStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *SomeStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "id"
	if fields["id"] != nil {
		if string(fields["id"]) != "null" {
			if err := json.Unmarshal(fields["id"], &resource.Id); err != nil {
				errs = append(errs, cog.MakeBuildErrors("id", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("id", errors.New("required field is null"))...)
		
		}
		delete(fields, "id")
	} else {errs = append(errs, cog.MakeBuildErrors("id", errors.New("required field is missing from input"))...)
	}
	// Field "maybeId"
	if fields["maybeId"] != nil {
		if string(fields["maybeId"]) != "null" {
			if err := json.Unmarshal(fields["maybeId"], &resource.MaybeId); err != nil {
				errs = append(errs, cog.MakeBuildErrors("maybeId", err)...)
			}
		
		}
		delete(fields, "maybeId")
	
	}
	// Field "title"
	if fields["title"] != nil {
		if string(fields["title"]) != "null" {
			if err := json.Unmarshal(fields["title"], &resource.Title); err != nil {
				errs = append(errs, cog.MakeBuildErrors("title", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("title", errors.New("required field is null"))...)
		
		}
		delete(fields, "title")
	} else {errs = append(errs, cog.MakeBuildErrors("title", errors.New("required field is missing from input"))...)
	}
	// Field "refStruct"
	if fields["refStruct"] != nil {
		if string(fields["refStruct"]) != "null" {
			
			resource.RefStruct = &RefStruct{}
			if err := resource.RefStruct.UnmarshalJSONStrict(fields["refStruct"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("refStruct", err)...)
			}
		
		}
		delete(fields, "refStruct")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("SomeStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `SomeStruct` objects.
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


// Validate checks all the validation constraints that may be defined on `SomeStruct` fields for violations and returns them.
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

// NewRefStruct creates a new RefStruct object.
func NewRefStruct() *RefStruct {
	return &RefStruct{
		Labels: map[string]string{},
		Tags: []string{},
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `RefStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *RefStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "labels"
	if fields["labels"] != nil {
		if string(fields["labels"]) != "null" {
			
			if err := json.Unmarshal(fields["labels"], &resource.Labels); err != nil {
				errs = append(errs, cog.MakeBuildErrors("labels", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("labels", errors.New("required field is null"))...)
		
		}
		delete(fields, "labels")
	} else {errs = append(errs, cog.MakeBuildErrors("labels", errors.New("required field is missing from input"))...)
	}
	// Field "tags"
	if fields["tags"] != nil {
		if string(fields["tags"]) != "null" {
			
			if err := json.Unmarshal(fields["tags"], &resource.Tags); err != nil {
				errs = append(errs, cog.MakeBuildErrors("tags", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("tags", errors.New("required field is null"))...)
		
		}
		delete(fields, "tags")
	} else {errs = append(errs, cog.MakeBuildErrors("tags", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("RefStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `RefStruct` objects.
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


// Validate checks all the validation constraints that may be defined on `RefStruct` fields for violations and returns them.
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


