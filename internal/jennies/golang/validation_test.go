package golang_test

import (
	"testing"

	"github.com/grafana/cog/testdata/generated/validation"
	"github.com/stretchr/testify/require"
)

func TestValidation_withValidInputs(t *testing.T) {
	req := require.New(t)

	testCases := []struct {
		input validation.Dashboard
	}{
		{
			input: validation.Dashboard{Title: "Hey there"},
		},
		{
			input: validation.Dashboard{
				Uid:    toPtr("uid"),
				Id:     toPtr[int64](42),
				Title:  "Hey there",
				Tags:   []string{"tag1", "tag2"},
				Labels: map[string]string{"foo": "bar"},
				Panels: []validation.Panel{
					{Title: "Some panel"},
				},
			},
		},
	}

	for _, testCase := range testCases {
		req.NoError(testCase.input.Validate())
	}
}

func TestValidation_withInvalidInputs(t *testing.T) {
	req := require.New(t)

	testCases := []struct {
		input             validation.Dashboard
		expectedErrorsMsg string
	}{
		{
			input:             validation.Dashboard{},
			expectedErrorsMsg: "title: must be >= 1",
		},
		{
			input: validation.Dashboard{
				Uid: toPtr(""),
				Id:  toPtr[int64](-1),
			},
			expectedErrorsMsg: `uid: must be >= 1
id: must be > 0
title: must be >= 1`,
		},
		{
			input: validation.Dashboard{
				Title:  "not empty",
				Tags:   []string{"foo", "", "bar"},
				Labels: map[string]string{"foo": "bar", "empty": "", "lala": "lolo"},
			},
			expectedErrorsMsg: `tags[1]: must be >= 1
labels[empty]: must be >= 1`,
		},
		{
			input: validation.Dashboard{
				Title: "not empty",
				Panels: []validation.Panel{
					{Title: ""},
				},
			},
			expectedErrorsMsg: `panels[0].title: must be >= 1`,
		},
	}

	for _, testCase := range testCases {
		err := testCase.input.Validate()
		req.Error(err)
		req.Equal(testCase.expectedErrorsMsg, err.Error())
	}
}
