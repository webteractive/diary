package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"diary/internal/project"
)

func TestMigrateCopiesProjectStoreToUserStore(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	projectRoot := t.TempDir()
	resolution := project.Resolution{Name: "diary", Root: projectRoot}

	record, err := CreateRecord(CreateRecordOptions{
		Paths:   NewPaths(projectRoot, "diary"),
		Project: "diary",
		Message: "Project-local record",
		Now:     time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := Migrate(MigrateOptions{
		Resolution: resolution,
		From:       "project",
		To:         "user",
		Now:        time.Date(2026, 5, 1, 11, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Copied != 1 || result.Records != 1 {
		t.Fatalf("expected one copied record, got %#v", result)
	}

	userStore, err := ResolveNamedStore("user", resolution, time.Date(2026, 5, 1, 11, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(userStore.Paths.RecordsDir, filepath.Base(record.Path))); err != nil {
		t.Fatalf("expected migrated record: %v", err)
	}
	if _, err := ReadIndex(userStore.Paths); err != nil {
		t.Fatalf("expected rebuilt destination index: %v", err)
	}
	if _, err := os.Stat(NewPaths(projectRoot, "diary").ProjectDir); err != nil {
		t.Fatalf("expected source to remain by default: %v", err)
	}
}

func TestMigrateDryRunDoesNotWrite(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	projectRoot := t.TempDir()
	resolution := project.Resolution{Name: "diary", Root: projectRoot}

	if _, err := CreateRecord(CreateRecordOptions{
		Paths:   NewPaths(projectRoot, "diary"),
		Project: "diary",
		Message: "Project-local record",
	}); err != nil {
		t.Fatal(err)
	}

	result, err := Migrate(MigrateOptions{
		Resolution: resolution,
		From:       "project",
		To:         "user",
		DryRun:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.DryRun || result.Copied != 1 {
		t.Fatalf("expected dry-run copied count, got %#v", result)
	}
	if _, err := os.Stat(filepath.Join(home, ".diary", "projects")); !os.IsNotExist(err) {
		t.Fatalf("expected dry-run not to write projects dir, got %v", err)
	}
}

func TestMigrateDeleteSourceRemovesProjectDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	projectRoot := t.TempDir()
	resolution := project.Resolution{Name: "diary", Root: projectRoot}
	sourcePaths := NewPaths(projectRoot, "diary")

	if _, err := CreateRecord(CreateRecordOptions{
		Paths:   sourcePaths,
		Project: "diary",
		Message: "Project-local record",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := Migrate(MigrateOptions{
		Resolution:   resolution,
		From:         "project",
		To:           "user",
		DeleteSource: true,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(sourcePaths.ProjectDir); !os.IsNotExist(err) {
		t.Fatalf("expected source project dir to be removed, got %v", err)
	}
}
