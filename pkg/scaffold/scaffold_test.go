package scaffold_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/LinkSpawnDev/linkspawn-create-context/pkg/scaffold"
)

func testData() scaffold.Data {
	return scaffold.Data{
		ProjectName: "demo-project",
		Mission:     "Prove the scaffolder works.",
		Author:      "Test Author",
		Date:        "2026-06-11",
		CLIVersion:  "0.0.0-test",
	}
}

func TestRun_RendersVariablesAndStripsTmplSuffix(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	templates := fstest.MapFS{
		"CLAUDE.md.tmpl": {Data: []byte("# {{.ProjectName}} by {{.Author}} on {{.Date}}\n")},
	}

	written, err := scaffold.Run(scaffold.Options{TargetDir: tmpDir, Templates: templates, Data: testData()})
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if len(written) != 1 || written[0] != "CLAUDE.md" {
		t.Fatalf("expected written [CLAUDE.md], got %v", written)
	}

	got, err := os.ReadFile(filepath.Join(tmpDir, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	const want = "# demo-project by Test Author on 2026-06-11\n"
	if string(got) != want {
		t.Errorf("content mismatch:\n  got:  %q\n  want: %q", string(got), want)
	}
}

func TestRun_CreatesNestedDirectories(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	templates := fstest.MapFS{
		".context/AGENTS.md.tmpl":  {Data: []byte("# {{.ProjectName}} field manual\n")},
		".context/LESSONS.md.tmpl": {Data: []byte("# Lessons\n")},
	}

	if _, err := scaffold.Run(scaffold.Options{TargetDir: tmpDir, Templates: templates, Data: testData()}); err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}

	for _, rel := range []string{".context/AGENTS.md", ".context/LESSONS.md"} {
		if _, err := os.Stat(filepath.Join(tmpDir, rel)); err != nil {
			t.Errorf("expected %s to exist: %v", rel, err)
		}
	}
}

func TestRun_RefusesToOverwriteAndReportsAllConflicts(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	for _, rel := range []string{"CLAUDE.md", "GEMINI.md"} {
		if err := os.WriteFile(filepath.Join(tmpDir, rel), []byte("precious\n"), 0o644); err != nil {
			t.Fatalf("test setup: %v", err)
		}
	}
	templates := fstest.MapFS{
		"CLAUDE.md.tmpl": {Data: []byte("new\n")},
		"GEMINI.md.tmpl": {Data: []byte("new\n")},
	}

	_, err := scaffold.Run(scaffold.Options{TargetDir: tmpDir, Templates: templates, Data: testData()})
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	for _, rel := range []string{"CLAUDE.md", "GEMINI.md"} {
		if !strings.Contains(err.Error(), rel) {
			t.Errorf("conflict error should mention %s, got: %v", rel, err)
		}
		got, _ := os.ReadFile(filepath.Join(tmpDir, rel))
		if string(got) != "precious\n" {
			t.Errorf("%s was modified despite conflict abort", rel)
		}
	}
}

func TestRun_ForceOverwrites(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "CLAUDE.md"), []byte("old\n"), 0o644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	templates := fstest.MapFS{
		"CLAUDE.md.tmpl": {Data: []byte("new\n")},
	}

	if _, err := scaffold.Run(scaffold.Options{TargetDir: tmpDir, Templates: templates, Data: testData(), Force: true}); err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(tmpDir, "CLAUDE.md"))
	if string(got) != "new\n" {
		t.Errorf("expected forced overwrite, got %q", string(got))
	}
}

func TestRun_GitignoreAppendsAndStaysIdempotent(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte("node_modules/\n"), 0o644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	templates := fstest.MapFS{
		".gitignore.tmpl": {Data: []byte(scaffold.GitignoreMarker + "\n.DS_Store\n")},
	}
	opts := scaffold.Options{TargetDir: tmpDir, Templates: templates, Data: testData()}

	if _, err := scaffold.Run(opts); err != nil {
		t.Fatalf("first Run returned unexpected error: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(tmpDir, ".gitignore"))
	if !strings.HasPrefix(string(got), "node_modules/\n") {
		t.Errorf("existing .gitignore content was lost: %q", string(got))
	}
	if !strings.Contains(string(got), scaffold.GitignoreMarker) {
		t.Errorf("managed block was not appended: %q", string(got))
	}

	// Second run must not duplicate the managed block.
	if _, err := scaffold.Run(opts); err != nil {
		t.Fatalf("second Run returned unexpected error: %v", err)
	}
	got2, _ := os.ReadFile(filepath.Join(tmpDir, ".gitignore"))
	if strings.Count(string(got2), scaffold.GitignoreMarker) != 1 {
		t.Errorf("managed block duplicated on re-run:\n%s", string(got2))
	}
}

func TestRun_InvalidTemplateFailsBeforeWriting(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	templates := fstest.MapFS{
		"GOOD.md.tmpl": {Data: []byte("fine\n")},
		"BAD.md.tmpl":  {Data: []byte("{{.NoSuchField}}\n")},
	}

	_, err := scaffold.Run(scaffold.Options{TargetDir: tmpDir, Templates: templates, Data: testData()})
	if err == nil {
		t.Fatal("expected render error, got nil")
	}
	if _, statErr := os.Stat(filepath.Join(tmpDir, "GOOD.md")); !os.IsNotExist(statErr) {
		t.Error("no files should be written when any template fails to render")
	}
}
