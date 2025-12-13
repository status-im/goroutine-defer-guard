package plugin

import (
	"fmt"
	"log"
	"os"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/status-im/goroutine-defer-guard/pkg/analyzer"
)

const (
	pluginName = "goroutine-defer-guard"
)

// Settings configures the golangci-lint plugin.
type Settings struct {
	// Target fully qualified handler identifier in the form full/pkg/path.Func.
	Target string `json:"target"`
}

type Plugin struct {
	settings Settings
}

func init() {
	register.Plugin(pluginName, New)
}

var _ register.LinterPlugin = (*Plugin)(nil)

// New constructs the plugin instance used by golangci-lint.
func New(conf any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[Settings](conf)
	if err != nil {
		return nil, fmt.Errorf("decode settings: %w", err)
	}

	return &Plugin{
		settings: settings,
	}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	gdg := analyzer.New(logger)

	if p.settings.Target != "" {
		if err := gdg.Flags.Set("target", p.settings.Target); err != nil {
			return nil, fmt.Errorf("set target flag: %w", err)
		}
	}

	return []*analysis.Analyzer{gdg}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
