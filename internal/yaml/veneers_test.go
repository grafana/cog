package yaml

import (
	"strings"
	"testing"

	"github.com/grafana/cog/internal/veneers/rewrite"
	"github.com/stretchr/testify/require"
)

func TestLoader_Load_withValidInputs(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		check func(req *require.Assertions, rules rewrite.LanguageRules)
	}{
		{
			desc: "no rules",
			input: `language: all
package: accesspolicy
builders: ~
options: ~`,
			check: func(req *require.Assertions, rules rewrite.LanguageRules) {
				req.Equal(rewrite.AllLanguages, rules.Language)
				req.Empty(rules.BuilderRules)
				req.Empty(rules.OptionRules)
			},
		},
		{
			desc: "single builder rule",
			input: `language: go
package: dashboard
builders: 
  - omit: { by_object: GridPos }
options: ~`,
			check: func(req *require.Assertions, rules rewrite.LanguageRules) {
				req.Equal("go", rules.Language)
				req.Len(rules.BuilderRules, 1)
				req.Empty(rules.OptionRules)
			},
		},
		{
			desc: "single option rule",
			input: `language: go
package: dashboard
builders: ~
options: 
  - omit: { by_name: Dashboard.schemaVersion }`,
			check: func(req *require.Assertions, rules rewrite.LanguageRules) {
				req.Equal("go", rules.Language)
				req.Empty(rules.BuilderRules, 1)
				req.Len(rules.OptionRules, 1)
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.desc, func(t *testing.T) {
			req := require.New(t)

			rules, err := NewVeneersLoader().load(strings.NewReader(tc.input))
			req.NoError(err)

			tc.check(req, rules)
		})
	}
}

func TestLoader_Load_withNoPackage(t *testing.T) {
	req := require.New(t)
	input := `language: all
builders: ~
options: ~`

	_, err := NewVeneersLoader().load(strings.NewReader(input))
	req.Error(err)
	req.ErrorContains(err, "missing 'package'")
}
