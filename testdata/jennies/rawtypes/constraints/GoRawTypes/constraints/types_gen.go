package constraints

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"strconv"
	"regexp"
)

type SomeStruct struct {
    Id uint64 `json:"id"`
    MaybeId *uint64 `json:"maybeId,omitempty"`
    GreaterThanZero uint64 `json:"greaterThanZero"`
    Negative int64 `json:"negative"`
    Title string `json:"title"`
    Labels map[string]string `json:"labels"`
    Tags []string `json:"tags"`
    Regex string `json:"regex"`
    NegativeRegex string `json:"negativeRegex"`
    MinMaxList []string `json:"minMaxList"`
    UniqueList []string `json:"uniqueList"`
    FullConstraintList []int64 `json:"fullConstraintList"`
}

// NewSomeStruct creates a new SomeStruct object.
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{
		Labels: map[string]string{},
		Tags: []string{},
		MinMaxList: []string{},
		UniqueList: []string{},
		FullConstraintList: []int64{},
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
	// Field "greaterThanZero"
	if fields["greaterThanZero"] != nil {
		if string(fields["greaterThanZero"]) != "null" {
			if err := json.Unmarshal(fields["greaterThanZero"], &resource.GreaterThanZero); err != nil {
				errs = append(errs, cog.MakeBuildErrors("greaterThanZero", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("greaterThanZero", errors.New("required field is null"))...)
		
		}
		delete(fields, "greaterThanZero")
	} else {errs = append(errs, cog.MakeBuildErrors("greaterThanZero", errors.New("required field is missing from input"))...)
	}
	// Field "negative"
	if fields["negative"] != nil {
		if string(fields["negative"]) != "null" {
			if err := json.Unmarshal(fields["negative"], &resource.Negative); err != nil {
				errs = append(errs, cog.MakeBuildErrors("negative", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("negative", errors.New("required field is null"))...)
		
		}
		delete(fields, "negative")
	} else {errs = append(errs, cog.MakeBuildErrors("negative", errors.New("required field is missing from input"))...)
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
	// Field "regex"
	if fields["regex"] != nil {
		if string(fields["regex"]) != "null" {
			if err := json.Unmarshal(fields["regex"], &resource.Regex); err != nil {
				errs = append(errs, cog.MakeBuildErrors("regex", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("regex", errors.New("required field is null"))...)
		
		}
		delete(fields, "regex")
	} else {errs = append(errs, cog.MakeBuildErrors("regex", errors.New("required field is missing from input"))...)
	}
	// Field "negativeRegex"
	if fields["negativeRegex"] != nil {
		if string(fields["negativeRegex"]) != "null" {
			if err := json.Unmarshal(fields["negativeRegex"], &resource.NegativeRegex); err != nil {
				errs = append(errs, cog.MakeBuildErrors("negativeRegex", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("negativeRegex", errors.New("required field is null"))...)
		
		}
		delete(fields, "negativeRegex")
	} else {errs = append(errs, cog.MakeBuildErrors("negativeRegex", errors.New("required field is missing from input"))...)
	}
	// Field "minMaxList"
	if fields["minMaxList"] != nil {
		if string(fields["minMaxList"]) != "null" {
			
			if err := json.Unmarshal(fields["minMaxList"], &resource.MinMaxList); err != nil {
				errs = append(errs, cog.MakeBuildErrors("minMaxList", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("minMaxList", errors.New("required field is null"))...)
		
		}
		delete(fields, "minMaxList")
	} else {errs = append(errs, cog.MakeBuildErrors("minMaxList", errors.New("required field is missing from input"))...)
	}
	// Field "uniqueList"
	if fields["uniqueList"] != nil {
		if string(fields["uniqueList"]) != "null" {
			
			if err := json.Unmarshal(fields["uniqueList"], &resource.UniqueList); err != nil {
				errs = append(errs, cog.MakeBuildErrors("uniqueList", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("uniqueList", errors.New("required field is null"))...)
		
		}
		delete(fields, "uniqueList")
	} else {errs = append(errs, cog.MakeBuildErrors("uniqueList", errors.New("required field is missing from input"))...)
	}
	// Field "fullConstraintList"
	if fields["fullConstraintList"] != nil {
		if string(fields["fullConstraintList"]) != "null" {
			
			if err := json.Unmarshal(fields["fullConstraintList"], &resource.FullConstraintList); err != nil {
				errs = append(errs, cog.MakeBuildErrors("fullConstraintList", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("fullConstraintList", errors.New("required field is null"))...)
		
		}
		delete(fields, "fullConstraintList")
	} else {errs = append(errs, cog.MakeBuildErrors("fullConstraintList", errors.New("required field is missing from input"))...)
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
		if resource.GreaterThanZero != other.GreaterThanZero {
			return false
		}
		if resource.Negative != other.Negative {
			return false
		}
		if resource.Title != other.Title {
			return false
		}

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
		if resource.Regex != other.Regex {
			return false
		}
		if resource.NegativeRegex != other.NegativeRegex {
			return false
		}

		if len(resource.MinMaxList) != len(other.MinMaxList) {
			return false
		}

		for i1 := range resource.MinMaxList {
		if resource.MinMaxList[i1] != other.MinMaxList[i1] {
			return false
		}
		}

		if len(resource.UniqueList) != len(other.UniqueList) {
			return false
		}

		for i1 := range resource.UniqueList {
		if resource.UniqueList[i1] != other.UniqueList[i1] {
			return false
		}
		}

		if len(resource.FullConstraintList) != len(other.FullConstraintList) {
			return false
		}

		for i1 := range resource.FullConstraintList {
		if resource.FullConstraintList[i1] != other.FullConstraintList[i1] {
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
		if !(resource.GreaterThanZero < 3) {
			errs = append(errs, cog.MakeBuildErrors(
				"greaterThanZero",
				errors.New("must be < 3"),
			)...)
		}
		if !(resource.Negative >= -19) {
			errs = append(errs, cog.MakeBuildErrors(
				"negative",
				errors.New("must be >= -19"),
			)...)
		}
		if !(resource.Negative < 10) {
			errs = append(errs, cog.MakeBuildErrors(
				"negative",
				errors.New("must be < 10"),
			)...)
		}
		if !(len([]rune(resource.Title)) >= 1) {
			errs = append(errs, cog.MakeBuildErrors(
				"title",
				errors.New("must be >= 1"),
			)...)
		}

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
		if !regexp.MustCompile("^[a-zA-Z0-9_-]+$").MatchString(resource.Regex) {
			errs = append(errs, cog.MakeBuildErrors(
				"regex",
				errors.New("must match regex ^[a-zA-Z0-9_-]+$"),
			)...)
		}
		if regexp.MustCompile("^[a-zA-Z0-9_-]+$").MatchString(resource.NegativeRegex) {
			errs = append(errs, cog.MakeBuildErrors(
				"negativeRegex",
				errors.New("must not match regex ^[a-zA-Z0-9_-]+$"),
			)...)
		}
		if len(resource.MinMaxList) < 1 {
			errs = append(errs, cog.MakeBuildErrors("minMaxList", errors.New("must have at least 1 items"))...)
		}
		if len(resource.MinMaxList) > 64 {
			errs = append(errs, cog.MakeBuildErrors("minMaxList", errors.New("must have at most 64 items"))...)
		}
		seenuniqueList := make(map[string]struct{})
		for _, item1 := range resource.UniqueList {
			key1 := fmt.Sprintf("%v", item1)
			if _, exists1 := seenuniqueList[key1]; exists1 {
				errs = append(errs, cog.MakeBuildErrors("uniqueList", errors.New("must have unique items"))...)
				break
			}
			seenuniqueList[key1] = struct{}{}
		}
		if len(resource.FullConstraintList) < 2 {
			errs = append(errs, cog.MakeBuildErrors("fullConstraintList", errors.New("must have at least 2 items"))...)
		}
		if len(resource.FullConstraintList) > 10 {
			errs = append(errs, cog.MakeBuildErrors("fullConstraintList", errors.New("must have at most 10 items"))...)
		}
		seenfullConstraintList := make(map[string]struct{})
		for _, item1 := range resource.FullConstraintList {
			key1 := fmt.Sprintf("%v", item1)
			if _, exists1 := seenfullConstraintList[key1]; exists1 {
				errs = append(errs, cog.MakeBuildErrors("fullConstraintList", errors.New("must have unique items"))...)
				break
			}
			seenfullConstraintList[key1] = struct{}{}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


