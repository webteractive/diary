package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateRecordWritesRecordLatestAndIndex(t *testing.T) {
	root := t.TempDir()

	record, err := CreateRecord(CreateRecordOptions{
		Root:    root,
		Project: "diary",
		Message: "Implemented storage layout",
		Type:    "context",
		Harness: "codex",
		Files:   []string{"internal/storage/record.go"},
		Now:     time.Date(2026, 5, 1, 10, 30, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}

	if record.ID == "" || record.Hash == "" {
		t.Fatalf("expected id and hash, got %#v", record)
	}
	if strings.Contains(record.ID, "codex") || strings.Contains(record.ID, "unknown") {
		t.Fatalf("expected record id to omit harness, got %q", record.ID)
	}
	if !strings.HasPrefix(record.ID, "2026-05-01T103000Z-") {
		t.Fatalf("expected timestamp-prefixed record id, got %q", record.ID)
	}
	if !strings.HasPrefix(record.Hash, "sha256:") {
		t.Fatalf("expected sha256 hash, got %q", record.Hash)
	}

	paths := NewPaths(root, "diary")
	for _, path := range []string{record.Path, paths.Latest, paths.Index} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}

	index, err := ReadIndex(paths)
	if err != nil {
		t.Fatal(err)
	}
	if len(index.Records) != 1 {
		t.Fatalf("expected one index record, got %d", len(index.Records))
	}
	if filepath.Base(index.Records[0].Path) != filepath.Base(record.Path) {
		t.Fatalf("expected index path to point to record")
	}
}

func TestFindByHashPrefixDetectsAmbiguity(t *testing.T) {
	root := t.TempDir()
	paths := NewPaths(root, "diary")
	if err := os.MkdirAll(paths.RecordsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	for _, record := range []Record{
		{ID: "one", Project: "diary", Hash: "sha256:abcdef", Type: "context", Timestamp: "2026-05-01T00:00:00Z", Body: "one"},
		{ID: "two", Project: "diary", Hash: "sha256:abc999", Type: "context", Timestamp: "2026-05-01T00:01:00Z", Body: "two"},
	} {
		data, err := RenderRecord(record)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(paths.RecordsDir, record.ID+".md"), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	_, err := FindByHashPrefix(paths, "abc")
	if err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("expected ambiguous hash error, got %v", err)
	}
}

func TestCreateRecordCanSeparateStorageIDFromProjectName(t *testing.T) {
	root := t.TempDir()

	record, err := CreateRecord(CreateRecordOptions{
		Paths:   NewDiaryRootPaths(root, "diary-12345678"),
		Project: "diary",
		Message: "Stored in mapped user project",
		Now:     time.Date(2026, 5, 1, 10, 30, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}

	if record.Project != "diary" {
		t.Fatalf("expected display project name, got %q", record.Project)
	}
	if !strings.Contains(record.Path, "diary-12345678") {
		t.Fatalf("expected record path to use storage id, got %q", record.Path)
	}
}
