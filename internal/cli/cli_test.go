package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecordListAndGet(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previous)
	})

	var out bytes.Buffer
	cmd := New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"record", "--project", "Diary", "--harness", "codex", "Implemented first slice"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "sha256:") {
		t.Fatalf("expected record output to include hash, got %q", out.String())
	}

	out.Reset()
	cmd = New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"list", "--project", "Diary"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Implemented first slice") {
		t.Fatalf("expected list preview, got %q", out.String())
	}

	out.Reset()
	cmd = New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"get", "--project", "Diary"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Implemented first slice") {
		t.Fatalf("expected get context, got %q", out.String())
	}
}

func TestListAndGetDefaultToCurrentProjectRoot(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	first := filepath.Join(t.TempDir(), "same")
	second := filepath.Join(t.TempDir(), "same")
	firstSubdir := filepath.Join(first, "internal")
	for _, dir := range []string{filepath.Join(first, ".git"), filepath.Join(second, ".git"), firstSubdir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previous)
	})

	if err := os.Chdir(first); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	cmd := New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"record", "--harness", "codex", "Alpha root context"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(second); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	cmd = New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"record", "--harness", "codex", "Beta root context"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(firstSubdir); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	cmd = New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if output := out.String(); !strings.Contains(output, "Alpha root context") || strings.Contains(output, "Beta root context") {
		t.Fatalf("expected list to stay scoped to first root, got %q", output)
	}

	out.Reset()
	cmd = New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"get"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if output := out.String(); !strings.Contains(output, "Alpha root context") || strings.Contains(output, "Beta root context") {
		t.Fatalf("expected get to stay scoped to first root, got %q", output)
	}
}

func TestListWithoutProjectDoesNotCreateProjectMapEntry(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previous)
	})

	var out bytes.Buffer
	cmd := New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(home, ".diary", "projects.json")); !os.IsNotExist(err) {
		t.Fatalf("expected read-only list not to create projects.json, got %v", err)
	}

	cmd = New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"get"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected get without records to fail")
	}

	if _, err := os.Stat(filepath.Join(home, ".diary", "projects.json")); !os.IsNotExist(err) {
		t.Fatalf("expected read-only get not to create projects.json, got %v", err)
	}
}

func TestRootHelpIncludesBannerAndVersionFlag(t *testing.T) {
	var out bytes.Buffer
	cmd := New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := out.String()
	for _, expected := range []string{"diary", "local AI context", "-v, --version"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected help to include %q, got %q", expected, output)
		}
	}
}

func TestVersionFlags(t *testing.T) {
	for _, args := range [][]string{{"--version"}, {"-v"}} {
		var out bytes.Buffer
		cmd := New(strings.NewReader(""), &out, &bytes.Buffer{})
		cmd.SetArgs(args)
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}

		if got, want := out.String(), "diary version dev\n"; got != want {
			t.Fatalf("expected %q for args %v, got %q", want, args, got)
		}
	}
}
