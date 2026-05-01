package cli

import (
	"bytes"
	"os"
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
