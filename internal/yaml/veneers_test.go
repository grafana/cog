package yaml

import (
	"os"
	"testing"

	"github.com/grafana/cog/internal/veneers/rewrite"
	"github.com/stretchr/testify/require"
)

func TestLoader_Load_withValidInputs(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		check func(req *require.Assertions, rules rewrite.RuleSet)
	}{
		{
			desc: "no rules",
			input: `languages: [all]
package: accesspolicy
builders: ~
options: ~`,
			check: func(req *require.Assertions, rules rewrite.RuleSet) {
				req.ElementsMatch([]string{rewrite.AllLanguages}, rules.Languages)
				req.Empty(rules.BuilderRules)
				req.Empty(rules.OptionRules)
			},
		},
		{
			desc: "single builder rule",
			input: `languages: [go]
package: dashboard
builders: 
  - omit: { by_object: GridPos }
options: ~`,
			check: func(req *require.Assertions, rules rewrite.RuleSet) {
				req.ElementsMatch([]string{"go"}, rules.Languages)
				req.Len(rules.BuilderRules, 1)
				req.Empty(rules.OptionRules)
			},
		},
		{
			desc: "single option rule",
			input: `languages: [go]
package: dashboard
builders: ~
options: 
  - omit: { by_name: Dashboard.schemaVersion }`,
			check: func(req *require.Assertions, rules rewrite.RuleSet) {
				req.ElementsMatch([]string{"go"}, rules.Languages)
				req.Empty(rules.BuilderRules, 1)
				req.Len(rules.OptionRules, 1)
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.desc, func(t *testing.T) {
			req := require.New(t)

			rules, err := NewVeneersLoader().load(tmpFile(t, []byte(tc.input)))
			req.NoError(err)

			tc.check(req, rules)
		})
	}
}

func TestLoader_Load_withNoPackage(t *testing.T) {
	req := require.New(t)
	input := `languages: [all]
builders: ~
options: ~`

	_, err := NewVeneersLoader().load(tmpFile(t, []byte(input)))
	req.Error(err)
	req.ErrorContains(err, "missing 'package'")
}

func tmpFile(t *testing.T, contents []byte) string {
	t.Helper()

	handle, err := os.CreateTemp(t.TempDir(), "cog-tests-veneers")
	require.NoError(t, err)

	_, err = handle.Write(contents)
	require.NoError(t, err)

	require.NoError(t, handle.Close())

	return handle.Name()
}
