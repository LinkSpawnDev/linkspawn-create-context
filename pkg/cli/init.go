package cli

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/LinkSpawnDev/linkspawn-create-context/pkg/scaffold"
)

func newInitCmd(cfg Config) *cobra.Command {
	var (
		presetName string
		name       string
		mission    string
		author     string
		force      bool
		yes        bool
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Scaffold .context/ and root pointer files into the current directory.",
		Long: `init renders the selected preset into the current working directory.

Existing files are never overwritten unless --force is given; an existing
.gitignore is appended to, not replaced. Run with --yes to skip the
interactive prompts and accept defaults (suitable for scripting).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to determine working directory: %w", err)
			}

			if name == "" {
				name = filepath.Base(cwd)
			}
			if author == "" {
				author = gitUserName()
			}
			if mission == "" {
				mission = "<one line: what this project is, for whom, and why it exists>"
			}

			interactive := !yes && term.IsTerminal(int(os.Stdin.Fd()))
			if interactive {
				if err := runForm(cfg, &name, &mission, &author, &presetName); err != nil {
					return err
				}
			}

			preset, err := cfg.preset(presetName)
			if err != nil {
				return err
			}
			templates, err := fs.Sub(cfg.Templates, preset.Dir)
			if err != nil {
				return fmt.Errorf("internal error: preset %q has no template dir %q: %w", preset.Name, preset.Dir, err)
			}

			written, err := scaffold.Run(scaffold.Options{
				TargetDir: cwd,
				Templates: templates,
				Data: scaffold.Data{
					ProjectName: name,
					Mission:     mission,
					Author:      author,
					Date:        time.Now().Format("2006-01-02"),
					CLIVersion:  cfg.Version,
				},
				Force: force,
			})
			if err != nil {
				return err
			}

			agentsPath := "AGENTS.md"
			for _, rel := range written {
				fmt.Printf("  created %s\n", rel)
				if rel == ".context/AGENTS.md" {
					agentsPath = rel
				}
			}
			fmt.Printf("\n✓ %s preset scaffolded for %q.\n", preset.Name, name)
			fmt.Printf("  Next: open %s and fill in the marked sections —\n", agentsPath)
			fmt.Println("  the role line, Definition of Done commands, and architecture mandates.")
			return nil
		},
	}

	cmd.Flags().StringVar(&presetName, "preset", cfg.DefaultPreset, "template preset to scaffold "+presetUsage(cfg))
	cmd.Flags().StringVar(&name, "name", "", "project name (default: current directory name)")
	cmd.Flags().StringVar(&mission, "mission", "", "one-line project mission")
	cmd.Flags().StringVar(&author, "author", "", "project owner (default: git config user.name)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing files")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "accept defaults; skip interactive prompts")
	return cmd
}

func runForm(cfg Config, name, mission, author, presetName *string) error {
	presetOptions := make([]huh.Option[string], 0, len(cfg.Presets))
	for _, p := range cfg.Presets {
		presetOptions = append(presetOptions, huh.NewOption(p.Name+" — "+p.Description, p.Name))
	}

	form := huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("Project name").Value(name),
		huh.NewInput().Title("Mission (one line — what, for whom, why)").Value(mission),
		huh.NewInput().Title("Owner").Value(author),
		huh.NewSelect[string]().Title("Preset").Options(presetOptions...).Value(presetName),
	))
	if err := form.Run(); err != nil {
		return fmt.Errorf("aborted: %w", err)
	}
	return nil
}

func presetUsage(cfg Config) string {
	names := make([]string, 0, len(cfg.Presets))
	for _, p := range cfg.Presets {
		names = append(names, p.Name)
	}
	return "(" + strings.Join(names, "|") + ")"
}

func gitUserName() string {
	out, err := exec.Command("git", "config", "user.name").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
