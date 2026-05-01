package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"diary/internal/project"
)

func TestResolveStoreDefaultsToUserRoot(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	projectRoot := t.TempDir()

	store, err := ResolveStore(StoreOptions{
		Resolution: project.Resolution{Name: "diary", Root: projectRoot},
		Now:        time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}

	if store.Location != "user" {
		t.Fatalf("expected user store, got %q", store.Location)
	}
	if !strings.HasPrefix(store.Paths.Diary, filepath.Join(home, ".diary")) {
		t.Fatalf("expected user Diary root, got %q", store.Paths.Diary)
	}
	if !strings.HasPrefix(store.Paths.Project, "diary-") {
		t.Fatalf("expected mapped project id, got %q", store.Paths.Project)
	}

	projectMap, err := ReadProjectMap(filepath.Join(home, ".diary"))
	if err != nil {
		t.Fatal(err)
	}
	if len(projectMap.Projects) != 1 || projectMap.Projects[0].Root != projectRoot {
		t.Fatalf("expected project map entry, got %#v", projectMap)
	}
}

func TestResolveStoreKeepsExistingProjectDiary(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	projectRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(projectRoot, ".diary", "projects", "diary"), 0o755); err != nil {
		t.Fatal(err)
	}

	store, err := ResolveStore(StoreOptions{
		Resolution: project.Resolution{Name: "diary", Root: projectRoot},
	})
	if err != nil {
		t.Fatal(err)
	}

	if store.Location != "project" {
		t.Fatalf("expected project store, got %q", store.Location)
	}
	if store.Paths.Diary != filepath.Join(projectRoot, ".diary") {
		t.Fatalf("expected project Diary root, got %q", store.Paths.Diary)
	}
	if _, err := os.Stat(filepath.Join(home, ".diary", "projects.json")); !os.IsNotExist(err) {
		t.Fatalf("expected no user project map when project .diary exists, got %v", err)
	}
}

func TestResolveStoreIgnoresEmptyProjectDiary(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	projectRoot := t.TempDir()
	if err := os.Mkdir(filepath.Join(projectRoot, ".diary"), 0o755); err != nil {
		t.Fatal(err)
	}

	store, err := ResolveStore(StoreOptions{
		Resolution: project.Resolution{Name: "diary", Root: projectRoot},
	})
	if err != nil {
		t.Fatal(err)
	}

	if store.Location != "user" {
		t.Fatalf("expected user store for empty project .diary, got %q", store.Location)
	}
}

func TestResolveStoreRootOverrideUsesCustomUserRoot(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	customRoot := filepath.Join(t.TempDir(), "work-diary")
	projectRoot := t.TempDir()
	if err := os.Mkdir(filepath.Join(projectRoot, ".diary"), 0o755); err != nil {
		t.Fatal(err)
	}

	store, err := ResolveStore(StoreOptions{
		Resolution:   project.Resolution{Name: "diary", Root: projectRoot},
		RootOverride: customRoot,
	})
	if err != nil {
		t.Fatal(err)
	}

	if store.Location != "user" {
		t.Fatalf("expected custom user store, got %q", store.Location)
	}
	if store.Paths.Diary != customRoot {
		t.Fatalf("expected custom Diary root, got %q", store.Paths.Diary)
	}
	if _, err := os.Stat(filepath.Join(customRoot, "projects.json")); err != nil {
		t.Fatalf("expected custom project map: %v", err)
	}
}

func TestRenameProjectMovesMappedProjectDirectory(t *testing.T) {
	diaryRoot := t.TempDir()
	projectRoot := t.TempDir()
	oldStore, err := resolveUserStore(diaryRoot, StoreOptions{
		Resolution: project.Resolution{Name: "old-name", Root: projectRoot},
		Now:        time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(oldStore.Paths.RecordsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(oldStore.Paths.RecordsDir, "record.md"), []byte("record"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := RenameProject(diaryRoot, project.Resolution{Name: "new-name", Root: projectRoot}, time.Date(2026, 5, 1, 11, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}

	if result.Old.ID == result.New.ID {
		t.Fatalf("expected project id to change, got %q", result.New.ID)
	}
	if !strings.HasPrefix(result.New.ID, "new-name-") {
		t.Fatalf("expected new project id, got %q", result.New.ID)
	}
	if _, err := os.Stat(oldStore.Paths.ProjectDir); !os.IsNotExist(err) {
		t.Fatalf("expected old project dir to move, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(diaryRoot, "projects", result.New.ID, "records", "record.md")); err != nil {
		t.Fatalf("expected record in renamed project dir: %v", err)
	}

	projectMap, err := ReadProjectMap(diaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(projectMap.Projects) != 1 || projectMap.Projects[0].ID != result.New.ID || projectMap.Projects[0].Name != "new-name" {
		t.Fatalf("expected renamed project map entry, got %#v", projectMap)
	}
}
