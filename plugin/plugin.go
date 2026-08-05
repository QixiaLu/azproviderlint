// Package plugin registers azproviderlint's analyzers as a golangci-lint module plugin.
package plugin

import (
	"fmt"
	"slices"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/katbyte/azproviderlint/checks"
)

func init() {
	register.Plugin("azproviderlint", New)
}

// Settings allows rules to be enabled/disabled per-rule from .golangci.yml via
// linters.settings.custom.azproviderlint.settings. An empty enable list means all rules.
type Settings struct {
	Enable  []string `json:"enable"`
	Disable []string `json:"disable"`
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}

	return &Plugin{settings: s}, nil
}

type Plugin struct {
	settings Settings
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	known := make(map[string]bool, len(checks.All))
	for _, a := range checks.All {
		known[a.Name] = true
	}
	for _, name := range slices.Concat(p.settings.Enable, p.settings.Disable) {
		if !known[name] {
			return nil, fmt.Errorf("unknown azproviderlint rule %q in settings", name)
		}
	}

	analyzers := make([]*analysis.Analyzer, 0, len(checks.All))
	for _, a := range checks.All {
		if len(p.settings.Enable) > 0 && !slices.Contains(p.settings.Enable, a.Name) {
			continue
		}
		if slices.Contains(p.settings.Disable, a.Name) {
			continue
		}
		analyzers = append(analyzers, a)
	}

	return analyzers, nil
}

func (p *Plugin) GetLoadMode() string {
	// AZS001 resolves named types and aliases via the type checker
	return register.LoadModeTypesInfo
}
