package plugin

import (
	"testing"

	"github.com/golangci/plugin-module-register/register"
)

func TestPluginRegistersAndConfiguresAnalyzer(t *testing.T) {
	newPlugin, err := register.GetPlugin(pluginName)
	if err != nil {
		t.Fatalf("expected plugin %q to be registered: %v", pluginName, err)
	}

	p, err := newPlugin(map[string]any{"target": "example.Target"})
	if err != nil {
		t.Fatalf("unexpected error constructing plugin: %v", err)
	}

	analyzers, err := p.BuildAnalyzers()
	if err != nil {
		t.Fatalf("unexpected error building analyzers: %v", err)
	}

	if len(analyzers) != 1 {
		t.Fatalf("expected a single analyzer, got %d", len(analyzers))
	}

	if got := analyzers[0].Name; got != "goroutinedeferguard" {
		t.Fatalf("unexpected analyzer name: %s", got)
	}

	if got := analyzers[0].Flags.Lookup("target").Value.String(); got != "example.Target" {
		t.Fatalf("target flag not propagated to analyzer: %s", got)
	}

	if got := p.GetLoadMode(); got != register.LoadModeTypesInfo {
		t.Fatalf("unexpected load mode: %s", got)
	}
}
