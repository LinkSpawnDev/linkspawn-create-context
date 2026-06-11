// Package cli builds the spawn command tree around a template FS supplied by
// the binary. The public binary embeds the core/minimal presets; a private
// build can embed its own template overlay and preset list and reuse this
// package unchanged.
package cli

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

// Preset names one selectable template tree inside Config.Templates.
type Preset struct {
	// Name is the value passed to --preset.
	Name string
	// Description is shown in help text and the interactive picker.
	Description string
	// Dir is the subdirectory of Config.Templates holding this preset.
	Dir string
}

// Config wires a binary's embedded templates into the shared command tree.
type Config struct {
	BinaryName    string
	Version       string
	Templates     fs.FS
	Presets       []Preset
	DefaultPreset string
}

func (c Config) preset(name string) (Preset, error) {
	for _, p := range c.Presets {
		if p.Name == name {
			return p, nil
		}
	}
	names := make([]string, 0, len(c.Presets))
	for _, p := range c.Presets {
		names = append(names, p.Name)
	}
	return Preset{}, fmt.Errorf("unknown preset %q (available: %v)", name, names)
}

// NewRoot assembles the root command with all subcommands attached.
func NewRoot(cfg Config) *cobra.Command {
	root := &cobra.Command{
		Use:   cfg.BinaryName,
		Short: "Scaffold a .context/ operating system for AI-agent-driven development.",
		Long: cfg.BinaryName + ` scaffolds a .context/ directory — agent rules, project state,
architecture notes, lessons and decision logs — plus root pointer files so
any AI coding agent (Claude, Gemini, or the next one) lands in the same
context cascade.

Plain markdown in, plain markdown out: the generated files belong to you and
outlive any tool, including this one.`,
		SilenceUsage: true,
		Version:      cfg.Version,
	}
	root.AddCommand(newInitCmd(cfg))
	return root
}

// Execute runs the root command and exits non-zero on error.
func Execute(cfg Config) {
	if err := NewRoot(cfg).Execute(); err != nil {
		os.Exit(1)
	}
}
