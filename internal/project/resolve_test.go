package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveUsesExplicitProject(t *testing.T) {
	dir := t.TempDir()

	resolved, err := Resolve(Options{Project: "Campaign Builder", WorkDir: dir})
	if err != nil {
		t.Fatal(err)
	}

	if resolved.Name != "campaign-builder" {
		t.Fatalf("expected sanitized project, got %q", resolved.Name)
	}
}

func TestResolveUsesConfigProject(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".diary"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".diary", "config.yml"), []byte("project: Configured App\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	resolved, err := Resolve(Options{WorkDir: dir})
	if err != nil {
		t.Fatal(err)
	}

	if resolved.Name != "configured-app" {
		t.Fatalf("expected configured project, got %q", resolved.Name)
	}
}

func TestResolveUsesGitRootFromSubdirectory(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "internal", "thing")
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}

	resolved, err := Resolve(Options{WorkDir: subdir})
	if err != nil {
		t.Fatal(err)
	}

	if resolved.Root != dir {
		t.Fatalf("expected git root %q, got %q", dir, resolved.Root)
	}
	if resolved.Name != filepath.Base(dir) {
		t.Fatalf("expected root basename, got %q", resolved.Name)
	}
}
