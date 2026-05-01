package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"diary/internal/install"
)

func TestInitGlobalInstructionDryRunDoesNotWrite(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	results, err := Init(Options{
		Target: install.TargetCodex,
		Scope:  ScopeGlobal,
		DryRun: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 || !results[0].DryRun || !results[0].Installed {
		t.Fatalf("expected dry-run instruction result, got %#v", results)
	}
	if _, err := os.Stat(filepath.Join(home, ".codex", "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("expected no instruction file to be written, got %v", err)
	}
}

func TestInitGlobalInstructionWritesCodexFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	results, err := Init(Options{
		Target: install.TargetCodex,
		Scope:  ScopeGlobal,
	})
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(home, ".codex", "AGENTS.md")
	if len(results) != 1 || results[0].Path != path || !results[0].Installed {
		t.Fatalf("expected installed instruction result, got %#v", results)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "<diary>") || !strings.Contains(string(data), "diary-record") {
		t.Fatalf("expected Diary instruction block, got %q", string(data))
	}
}

func TestInitProjectInstructionWritesHarnessFileAtProjectRoot(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	results, err := Init(Options{
		Target:  install.TargetClaude,
		Scope:   ScopeProject,
		WorkDir: filepath.Join(dir, "nested"),
	})
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "CLAUDE.md")
	if len(results) != 1 || results[0].Path != path {
		t.Fatalf("expected project instruction at %s, got %#v", path, results)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "</diary>") {
		t.Fatalf("expected Diary instruction block, got %q", string(data))
	}
}

func TestInitInstructionIsIdempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if _, err := Init(Options{Target: install.TargetCodex, Scope: ScopeGlobal}); err != nil {
		t.Fatal(err)
	}
	results, err := Init(Options{Target: install.TargetCodex, Scope: ScopeGlobal})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 || !results[0].Unchanged {
		t.Fatalf("expected unchanged result, got %#v", results)
	}
}

func TestInitWithInstallSkillsInstallsInstructionAndSkills(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	results, err := Init(Options{
		Target:        install.TargetCodex,
		Scope:         ScopeGlobal,
		InstallSkills: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 5 {
		t.Fatalf("expected one instruction and four skills, got %#v", results)
	}
	if _, err := os.Stat(filepath.Join(home, ".codex", "skills", "diary-record", "SKILL.md")); err != nil {
		t.Fatalf("expected diary-record skill to be installed: %v", err)
	}
}
