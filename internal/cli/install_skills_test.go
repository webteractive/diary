package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestInstallSkillsDryRun(t *testing.T) {
	var out bytes.Buffer
	cmd := New(strings.NewReader(""), &out, &bytes.Buffer{})
	cmd.SetArgs([]string{"install-skills", "--target", "codex", "--dry-run"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "would install codex diary-record skill") {
		t.Fatalf("expected dry-run output, got %q", out.String())
	}
}

func TestInstallSkillsRequiresTarget(t *testing.T) {
	cmd := New(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	cmd.SetArgs([]string{"install-skills", "--dry-run"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "required flag") {
		t.Fatalf("expected required target error, got %v", err)
	}
}
