package plugin

import (
	"testing"

	"github.com/katbyte/azproviderlint/checks"
)

func buildAnalyzers(t *testing.T, settings any) ([]string, error) {
	t.Helper()

	p, err := New(settings)
	if err != nil {
		return nil, err
	}

	analyzers, err := p.BuildAnalyzers()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(analyzers))
	for _, a := range analyzers {
		names = append(names, a.Name)
	}
	return names, nil
}

func TestBuildAnalyzersDefault(t *testing.T) {
	t.Parallel()

	names, err := buildAnalyzers(t, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != len(checks.All) {
		t.Fatalf("expected all %d analyzers with no settings, got %d", len(checks.All), len(names))
	}
}

func TestBuildAnalyzersDisable(t *testing.T) {
	t.Parallel()

	names, err := buildAnalyzers(t, map[string]any{"disable": []string{"AZR002"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != len(checks.All)-1 {
		t.Fatalf("expected %d analyzers, got %d", len(checks.All)-1, len(names))
	}
	for _, name := range names {
		if name == "AZR002" {
			t.Fatal("AZR002 should have been disabled")
		}
	}
}

func TestBuildAnalyzersEnable(t *testing.T) {
	t.Parallel()

	names, err := buildAnalyzers(t, map[string]any{"enable": []string{"AZS001", "AZT001"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 || names[0] != "AZS001" || names[1] != "AZT001" {
		t.Fatalf("expected exactly [AZS001 AZT001], got %v", names)
	}
}

func TestBuildAnalyzersUnknownRule(t *testing.T) {
	t.Parallel()

	if _, err := buildAnalyzers(t, map[string]any{"disable": []string{"AZX999"}}); err == nil {
		t.Fatal("expected an error for an unknown rule name")
	}
}

func TestBuildAnalyzersRuleFlags(t *testing.T) {
	t.Parallel()

	names, err := buildAnalyzers(t, map[string]any{
		"enable": []string{"AZS006"},
		"AZS006": map[string]any{"ignore-sensitive": true}, // unquoted YAML bool
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 1 || names[0] != "AZS006" {
		t.Fatalf("expected exactly [AZS006], got %v", names)
	}

	p, err := New(map[string]any{"AZS006": map[string]any{"ignore-sensitive": "true"}})
	if err != nil {
		t.Fatal(err)
	}
	pl, ok := p.(*Plugin)
	if !ok {
		t.Fatalf("expected *Plugin, got %T", p)
	}
	analyzers, err := pl.BuildAnalyzers()
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range analyzers {
		if a.Name != "AZS006" {
			continue
		}
		if f := a.Flags.Lookup("ignore-sensitive"); f == nil || f.Value.String() != "true" {
			t.Fatalf("expected AZS006 ignore-sensitive flag to be true, got %v", f)
		}
	}
}

func TestBuildAnalyzersRuleFlagErrors(t *testing.T) {
	t.Parallel()

	if _, err := buildAnalyzers(t, map[string]any{"AZX999": map[string]any{"some-flag": "true"}}); err == nil {
		t.Fatal("expected an error for flags on an unknown rule name")
	}
	if _, err := New(map[string]any{"AZS006": "not-a-map"}); err == nil {
		t.Fatal("expected an error for a non-map rule settings value")
	}
	if _, err := buildAnalyzers(t, map[string]any{"AZS006": map[string]any{"no-such-flag": "true"}}); err == nil {
		t.Fatal("expected an error for an unknown flag name")
	}
}
