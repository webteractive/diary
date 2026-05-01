package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"diary/internal/storage"
)

func TestMigrateCommandDryRun(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := storage.CreateRecord(storage.CreateRecordOptions{
		Paths:   storage.NewPaths(dir, filepath.Base(dir)),
		Project: filepath.Base(dir),
		Message: "Project-local record",
		Now:     time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC),
	}); err != nil {
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
	cmd.SetArgs([]string{"migrate", "--from", "project", "--to", "user", "--dry-run"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "would migrate 1 records") {
		t.Fatalf("expected dry-run migration output, got %q", out.String())
	}
}

func TestMigrateCommandRequiresFromAndTo(t *testing.T) {
	cmd := New(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	cmd.SetArgs([]string{"migrate", "--from", "project"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "required flag") {
		t.Fatalf("expected required to error, got %v", err)
	}
}
