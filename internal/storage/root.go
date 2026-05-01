package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"diary/internal/project"
)

const EnvRoot = "DIARY_ROOT"

type StoreOptions struct {
	Resolution   project.Resolution
	RootOverride string
	Now          time.Time
}

type Store struct {
	Paths    Paths
	Location string
	Entry    ProjectEntry
}

type ProjectMap struct {
	Version  int            `json:"version"`
	Projects []ProjectEntry `json:"projects"`
}

type ProjectEntry struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Root       string `json:"root"`
	CreatedAt  string `json:"created_at"`
	LastSeenAt string `json:"last_seen_at"`
}

type RenameProjectResult struct {
	Old ProjectEntry
	New ProjectEntry
}

func ResolveStore(opts StoreOptions) (Store, error) {
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}

	rootOverride := firstNonEmpty(opts.RootOverride, os.Getenv(EnvRoot))
	if rootOverride != "" {
		return resolveUserStore(rootOverride, opts)
	}

	localPaths := NewPaths(opts.Resolution.Root, opts.Resolution.Name)
	if existsDir(localPaths.ProjectDir) {
		return Store{
			Paths:    localPaths,
			Location: "project",
		}, nil
	}

	diaryRoot, err := DefaultRoot()
	if err != nil {
		return Store{}, err
	}
	return resolveUserStore(diaryRoot, opts)
}

func ResolveStoreForRoot(opts StoreOptions) (Store, error) {
	rootOverride := firstNonEmpty(opts.RootOverride, os.Getenv(EnvRoot))

	localPaths := NewPaths(opts.Resolution.Root, opts.Resolution.Name)
	if rootOverride == "" && existsDir(localPaths.ProjectDir) {
		return Store{
			Paths:    localPaths,
			Location: "project",
		}, nil
	}

	diaryRoot := rootOverride
	if diaryRoot == "" {
		var err error
		diaryRoot, err = DefaultRoot()
		if err != nil {
			return Store{}, err
		}
	}

	return resolveUserStoreForRoot(diaryRoot, opts)
}

func ResolveNamedStore(name string, resolution project.Resolution, now time.Time) (Store, error) {
	switch name {
	case "", "user":
		root, err := DefaultRoot()
		if err != nil {
			return Store{}, err
		}
		return resolveUserStore(root, StoreOptions{Resolution: resolution, Now: now})
	case "project":
		return Store{
			Paths:    NewPaths(resolution.Root, resolution.Name),
			Location: "project",
		}, nil
	default:
		return resolveUserStore(name, StoreOptions{Resolution: resolution, Now: now})
	}
}

func DefaultRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".diary"), nil
}

func ProjectsInRoot(diaryRoot string) ([]string, error) {
	projectMap, err := ReadProjectMap(diaryRoot)
	if err != nil {
		return nil, err
	}
	if len(projectMap.Projects) > 0 {
		projects := make([]string, 0, len(projectMap.Projects))
		for _, entry := range projectMap.Projects {
			projects = append(projects, entry.ID)
		}
		sort.Strings(projects)
		return projects, nil
	}

	entries, err := os.ReadDir(filepath.Join(diaryRoot, "projects"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	projects := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			projects = append(projects, entry.Name())
		}
	}
	sort.Strings(projects)
	return projects, nil
}

func ReadProjectMap(diaryRoot string) (ProjectMap, error) {
	path := filepath.Join(diaryRoot, "projects.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ProjectMap{Version: 1}, nil
		}
		return ProjectMap{}, err
	}

	var projectMap ProjectMap
	if err := json.Unmarshal(data, &projectMap); err != nil {
		return ProjectMap{}, err
	}
	if projectMap.Version == 0 {
		projectMap.Version = 1
	}
	return projectMap, nil
}

func RenameProject(diaryRoot string, resolution project.Resolution, now time.Time) (RenameProjectResult, error) {
	if now.IsZero() {
		now = time.Now().UTC()
	}

	diaryRoot, err := filepath.Abs(diaryRoot)
	if err != nil {
		return RenameProjectResult{}, err
	}
	root, err := filepath.Abs(resolution.Root)
	if err != nil {
		return RenameProjectResult{}, err
	}

	projectMap, err := ReadProjectMap(diaryRoot)
	if err != nil {
		return RenameProjectResult{}, err
	}

	for i, entry := range projectMap.Projects {
		if entry.Root != root {
			continue
		}

		oldEntry := entry
		newEntry := entry
		newEntry.ID = projectID(resolution.Name, root)
		newEntry.Name = resolution.Name
		newEntry.LastSeenAt = now.UTC().Format(time.RFC3339)

		if oldEntry.ID != newEntry.ID {
			oldPath := NewDiaryRootPaths(diaryRoot, oldEntry.ID).ProjectDir
			newPath := NewDiaryRootPaths(diaryRoot, newEntry.ID).ProjectDir
			if existsDir(oldPath) {
				if existsDir(newPath) {
					return RenameProjectResult{}, fmt.Errorf("target project already exists: %s", newEntry.ID)
				}
				if err := os.Rename(oldPath, newPath); err != nil {
					return RenameProjectResult{}, err
				}
			}
		}

		projectMap.Projects[i] = newEntry
		sort.Slice(projectMap.Projects, func(i, j int) bool {
			return projectMap.Projects[i].ID < projectMap.Projects[j].ID
		})
		if err := writeProjectMap(diaryRoot, projectMap); err != nil {
			return RenameProjectResult{}, err
		}
		return RenameProjectResult{Old: oldEntry, New: newEntry}, nil
	}

	return RenameProjectResult{}, fmt.Errorf("no diary project found for root: %s", root)
}

func resolveUserStore(diaryRoot string, opts StoreOptions) (Store, error) {
	diaryRoot, err := filepath.Abs(diaryRoot)
	if err != nil {
		return Store{}, err
	}

	entry, err := upsertProjectMapEntry(diaryRoot, opts)
	if err != nil {
		return Store{}, err
	}

	return Store{
		Paths:    NewDiaryRootPaths(diaryRoot, entry.ID),
		Location: "user",
		Entry:    entry,
	}, nil
}

func resolveUserStoreForRoot(diaryRoot string, opts StoreOptions) (Store, error) {
	diaryRoot, err := filepath.Abs(diaryRoot)
	if err != nil {
		return Store{}, err
	}

	root, err := filepath.Abs(opts.Resolution.Root)
	if err != nil {
		return Store{}, err
	}

	projectMap, err := ReadProjectMap(diaryRoot)
	if err != nil {
		return Store{}, err
	}

	for _, entry := range projectMap.Projects {
		if entry.Root == root {
			return Store{
				Paths:    NewDiaryRootPaths(diaryRoot, entry.ID),
				Location: "user",
				Entry:    entry,
			}, nil
		}
	}

	return Store{
		Paths:    NewDiaryRootPaths(diaryRoot, projectID(opts.Resolution.Name, root)),
		Location: "user",
	}, nil
}

func upsertProjectMapEntry(diaryRoot string, opts StoreOptions) (ProjectEntry, error) {
	root, err := filepath.Abs(opts.Resolution.Root)
	if err != nil {
		return ProjectEntry{}, err
	}

	projectMap, err := ReadProjectMap(diaryRoot)
	if err != nil {
		return ProjectEntry{}, err
	}

	now := opts.Now.UTC().Format(time.RFC3339)
	for i, entry := range projectMap.Projects {
		if entry.Root == root {
			projectMap.Projects[i].Name = opts.Resolution.Name
			projectMap.Projects[i].LastSeenAt = now
			if err := writeProjectMap(diaryRoot, projectMap); err != nil {
				return ProjectEntry{}, err
			}
			return projectMap.Projects[i], nil
		}
	}

	entry := ProjectEntry{
		ID:         projectID(opts.Resolution.Name, root),
		Name:       opts.Resolution.Name,
		Root:       root,
		CreatedAt:  now,
		LastSeenAt: now,
	}
	projectMap.Projects = append(projectMap.Projects, entry)
	sort.Slice(projectMap.Projects, func(i, j int) bool {
		return projectMap.Projects[i].ID < projectMap.Projects[j].ID
	})

	if err := writeProjectMap(diaryRoot, projectMap); err != nil {
		return ProjectEntry{}, err
	}
	return entry, nil
}

func writeProjectMap(diaryRoot string, projectMap ProjectMap) error {
	if projectMap.Version == 0 {
		projectMap.Version = 1
	}
	if err := os.MkdirAll(diaryRoot, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(projectMap, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(diaryRoot, "projects.json"), data, 0o644)
}

func projectID(name, root string) string {
	sum := sha256.Sum256([]byte(root))
	return fmt.Sprintf("%s-%s", project.Sanitize(name), hex.EncodeToString(sum[:])[:8])
}

func existsDir(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
