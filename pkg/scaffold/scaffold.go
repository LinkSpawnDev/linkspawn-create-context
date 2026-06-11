// Package scaffold renders an embedded template tree into a target directory.
//
// The template FS root mirrors the output tree: a file at
// ".context/AGENTS.md.tmpl" is rendered to "<target>/.context/AGENTS.md".
// Every file is parsed as a text/template against Data; a trailing ".tmpl"
// suffix is stripped from the output name.
package scaffold

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Data holds the variables available to every template file.
type Data struct {
	ProjectName string
	Mission     string
	Author      string
	Date        string // YYYY-MM-DD
	CLIVersion  string
}

// Options configures a single scaffold run.
type Options struct {
	// TargetDir is the directory the rendered tree is written into.
	TargetDir string
	// Templates is the template tree; its root mirrors TargetDir.
	Templates fs.FS
	// Data is passed to every template.
	Data Data
	// Force overwrites existing files instead of aborting.
	Force bool
}

// GitignoreMarker guards the managed block appended to an existing
// .gitignore so repeated runs stay idempotent.
const GitignoreMarker = "# --- spawn:context (managed) ---"

type renderedFile struct {
	rel      string // slash-separated output path, .tmpl stripped
	data     []byte
	appendTo bool // append to the existing file instead of writing fresh
	skip     bool // nothing to do (managed block already present)
}

// Run renders every file in opts.Templates into opts.TargetDir and returns
// the relative paths written, in walk order.
//
// Unless opts.Force is set, Run refuses to overwrite any existing file and
// reports every conflict before writing anything. A .gitignore template is
// the exception: its content is appended to an existing .gitignore, and it
// is skipped entirely when GitignoreMarker is already present.
func Run(opts Options) ([]string, error) {
	files, err := renderAll(opts.Templates, opts.Data)
	if err != nil {
		return nil, err
	}

	var conflicts []string
	for i := range files {
		f := &files[i]
		dest := filepath.Join(opts.TargetDir, filepath.FromSlash(f.rel))
		existing, err := os.ReadFile(dest)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to inspect %s: %w", dest, err)
		}
		if filepath.Base(f.rel) == ".gitignore" {
			if bytes.Contains(existing, []byte(GitignoreMarker)) {
				f.skip = true
			} else {
				f.appendTo = true
			}
			continue
		}
		if !opts.Force {
			conflicts = append(conflicts, f.rel)
		}
	}
	if len(conflicts) > 0 {
		return nil, fmt.Errorf(
			"refusing to overwrite existing files (use --force to override):\n  %s",
			strings.Join(conflicts, "\n  "),
		)
	}

	var written []string
	for _, f := range files {
		if f.skip {
			continue
		}
		dest := filepath.Join(opts.TargetDir, filepath.FromSlash(f.rel))
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return nil, fmt.Errorf("failed to create directory for %s: %w", f.rel, err)
		}
		if f.appendTo {
			if err := appendFile(dest, f.data); err != nil {
				return nil, err
			}
		} else if err := os.WriteFile(dest, f.data, 0o644); err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", f.rel, err)
		}
		written = append(written, f.rel)
	}
	return written, nil
}

func renderAll(templates fs.FS, data Data) ([]renderedFile, error) {
	var files []renderedFile
	err := fs.WalkDir(templates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		raw, err := fs.ReadFile(templates, path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}
		tmpl, err := template.New(path).Parse(string(raw))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("failed to render template %s: %w", path, err)
		}
		files = append(files, renderedFile{
			rel:  strings.TrimSuffix(path, ".tmpl"),
			data: buf.Bytes(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func appendFile(dest string, data []byte) error {
	f, err := os.OpenFile(dest, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open %s for append: %w", dest, err)
	}
	defer f.Close()
	if _, err := f.Write(append([]byte("\n"), data...)); err != nil {
		return fmt.Errorf("failed to append to %s: %w", dest, err)
	}
	return nil
}
