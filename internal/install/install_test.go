package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallDryRunDoesNotWrite(t *testing.T) {
	dir := t.TempDir()

	results, err := Install(Options{
		Target: TargetCodex,
		Path:   dir,
		DryRun: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 3 {
		t.Fatalf("expected three results, got %d", len(results))
	}
	if !results[0].DryRun || results[0].Installed {
		t.Fatalf("expected dry-run result, got %#v", results[0])
	}
	if _, err := os.Stat(filepath.Join(dir, "diary-record", "SKILL.md")); !os.IsNotExist(err) {
		t.Fatalf("expected no file to be written, stat err: %v", err)
	}
}

func TestInstallWritesSkill(t *testing.T) {
	dir := t.TempDir()

	results, err := Install(Options{
		Target: TargetClaude,
		Path:   dir,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 3 || !results[0].Installed {
		t.Fatalf("expected installed result, got %#v", results)
	}

	data, err := os.ReadFile(filepath.Join(dir, "diary-get", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "diary get") {
		t.Fatalf("expected skill content to mention diary get")
	}
}

func TestInstallRefusesExistingFileWithoutForce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "diary-record", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Install(Options{
		Target: TargetCodex,
		Path:   dir,
	})
	if err == nil || !strings.Contains(err.Error(), "--force") {
		t.Fatalf("expected force error, got %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "existing" {
		t.Fatalf("expected existing file to be preserved")
	}
}

func TestInstallForceOverwritesExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "diary-record", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := Install(Options{
		Target: TargetCodex,
		Path:   dir,
		Force:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("expected three results, got %#v", results)
	}
	var foundOverwrite bool
	for _, result := range results {
		if result.Skill == "diary-record" && result.Overwritten {
			foundOverwrite = true
		}
	}
	if !foundOverwrite {
		t.Fatalf("expected overwritten result, got %#v", results)
	}
}

func TestInstallRejectsCustomPathWithAllTargets(t *testing.T) {
	_, err := Install(Options{
		Target: TargetAll,
		Path:   t.TempDir(),
	})
	if err == nil || !strings.Contains(err.Error(), "--path") {
		t.Fatalf("expected custom path error, got %v", err)
	}
}
