package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitDryRun(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	var out bytes.Buffer
	cmd := New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"init", "--target", "codex", "--dry-run"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "would install codex global instruction") {
		t.Fatalf("expected dry-run output, got %q", out.String())
	}
}

func TestInitRequiresTarget(t *testing.T) {
	cmd := New(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	cmd.SetArgs([]string{"init", "--dry-run"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "required flag") {
		t.Fatalf("expected required target error, got %v", err)
	}
}

func TestInitWithInstallSkillsOutput(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var out bytes.Buffer
	cmd := New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"init", "--target", "codex", "--install-skills"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := out.String()
	if !strings.Contains(output, filepath.Join(home, ".codex", "AGENTS.md")) {
		t.Fatalf("expected instruction path in output, got %q", output)
	}
	if !strings.Contains(output, "installed codex diary-record skill") {
		t.Fatalf("expected skill install output, got %q", output)
	}
}
