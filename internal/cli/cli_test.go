package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRecordListAndGet(t *testing.T) {
	dir := t.TempDir()
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
