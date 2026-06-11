package main

import (
	"github.com/LinkSpawnDev/linkspawn-create-context/pkg/cli"
	"github.com/LinkSpawnDev/linkspawn-create-context/templates"
)

// version is stamped at release time via:
//
//	go build -ldflags "-X main.version=v0.1.0" ./cmd/spawn
var version = "dev"

func main() {
	cli.Execute(cli.Config{
		BinaryName: "spawn",
		Version:    version,
		Templates:  templates.FS,
		Presets: []cli.Preset{
			{Name: "core", Description: "full .context/ core tier + root pointers (recommended)", Dir: "core"},
			{Name: "minimal", Description: "single root AGENTS.md + CLAUDE.md pointer, no .context/", Dir: "minimal"},
		},
		DefaultPreset: "core",
	})
}
