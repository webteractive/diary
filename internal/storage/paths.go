package storage

import "path/filepath"

type Paths struct {
	Root       string
	Diary      string
	Project    string
	ProjectDir string
	RecordsDir string
	Latest     string
	Index      string
}

func NewPaths(root, project string) Paths {
	diary := filepath.Join(root, ".diary")
	return NewDiaryRootPaths(diary, project)
}

func NewDiaryRootPaths(diaryRoot, project string) Paths {
	diaryRoot = filepath.Clean(diaryRoot)
	projectDir := filepath.Join(diaryRoot, "projects", project)
	return Paths{
		Root:       filepath.Dir(diaryRoot),
		Diary:      diaryRoot,
		Project:    project,
		ProjectDir: projectDir,
		RecordsDir: filepath.Join(projectDir, "records"),
		Latest:     filepath.Join(projectDir, "latest.md"),
		Index:      filepath.Join(projectDir, "index.json"),
	}
}
