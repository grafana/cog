package codegen

import (
	"testing"

	"github.com/grafana/cog/internal/jennies/rust"
)

func TestOutputLanguages_Rust(t *testing.T) {
	pipeline := &Pipeline{
		Output: Output{
			Languages: []*OutputLanguage{
				{Rust: &rust.Config{}},
			},
		},
	}

	outputs, err := pipeline.OutputLanguages()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	language, ok := outputs[rust.LanguageRef]
	if !ok {
		t.Fatalf("expected output languages to contain key %q", rust.LanguageRef)
	}
	if language == nil {
		t.Fatalf("expected non-nil language for key %q", rust.LanguageRef)
	}
	if language.Name() != rust.LanguageRef {
		t.Fatalf("expected language under key %q to have Name() == %q, got %q", rust.LanguageRef, rust.LanguageRef, language.Name())
	}
}

func TestOutputLanguage_InterpolateRust(t *testing.T) {
	outputLanguage := &OutputLanguage{
		Rust: &rust.Config{PathPrefix: "{{ .x }}"},
	}

	interpolator := func(input string) string {
		if input == "{{ .x }}" {
			return "interpolated-value"
		}
		return input
	}

	outputLanguage.interpolateParameters(&Output{}, interpolator)

	if outputLanguage.Rust.PathPrefix != "interpolated-value" {
		t.Fatalf("expected PathPrefix to be interpolated to %q, got %q", "interpolated-value", outputLanguage.Rust.PathPrefix)
	}
}
