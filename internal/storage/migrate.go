package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"diary/internal/project"
)

type MigrateOptions struct {
	Resolution   project.Resolution
	From         string
	To           string
	Force        bool
	DeleteSource bool
	DryRun       bool
	Now          time.Time
}

type MigrateResult struct {
	From         string `json:"from"`
	To           string `json:"to"`
	FromPath     string `json:"from_path"`
	ToPath       string `json:"to_path"`
	Records      int    `json:"records"`
	Copied       int    `json:"copied"`
	Skipped      int    `json:"skipped"`
	DryRun       bool   `json:"dry_run"`
	DeleteSource bool   `json:"delete_source"`
}

func Migrate(opts MigrateOptions) (MigrateResult, error) {
	if opts.From == "" {
		return MigrateResult{}, fmt.Errorf("--from is required")
	}
	if opts.To == "" {
		return MigrateResult{}, fmt.Errorf("--to is required")
	}
	if opts.From == opts.To {
		return MigrateResult{}, fmt.Errorf("--from and --to must be different")
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}

	fromStore, err := ResolveNamedStore(opts.From, opts.Resolution, opts.Now)
	if err != nil {
		return MigrateResult{}, err
	}
	toStore, err := ResolveNamedStore(opts.To, opts.Resolution, opts.Now)
	if err != nil {
		return MigrateResult{}, err
	}

	records, err := ReadRecords(fromStore.Paths)
	if err != nil {
		return MigrateResult{}, err
	}

	result := MigrateResult{
		From:         opts.From,
		To:           opts.To,
		FromPath:     fromStore.Paths.ProjectDir,
		ToPath:       toStore.Paths.ProjectDir,
		Records:      len(records),
		DryRun:       opts.DryRun,
		DeleteSource: opts.DeleteSource,
	}

	for _, record := range records {
		destPath := filepath.Join(toStore.Paths.RecordsDir, filepath.Base(record.Path))
		if _, err := os.Stat(destPath); err == nil && !opts.Force {
			return MigrateResult{}, fmt.Errorf("%s already exists; pass --force to overwrite", destPath)
		} else if err != nil && !os.IsNotExist(err) {
			return MigrateResult{}, err
		}

		if opts.DryRun {
			result.Copied++
			continue
		}

		if err := os.MkdirAll(toStore.Paths.RecordsDir, 0o755); err != nil {
			return MigrateResult{}, err
		}
		data, err := os.ReadFile(record.Path)
		if err != nil {
			return MigrateResult{}, err
		}
		if err := os.WriteFile(destPath, data, 0o644); err != nil {
			return MigrateResult{}, err
		}
		result.Copied++
	}

	if opts.DryRun {
		return result, nil
	}

	if err := RebuildIndex(toStore.Paths); err != nil {
		return MigrateResult{}, err
	}
	if err := copyLatest(fromStore.Paths, toStore.Paths); err != nil {
		return MigrateResult{}, err
	}

	if opts.DeleteSource {
		if err := os.RemoveAll(fromStore.Paths.ProjectDir); err != nil {
			return MigrateResult{}, err
		}
	}

	return result, nil
}

func copyLatest(from Paths, to Paths) error {
	data, err := os.ReadFile(from.Latest)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if err := os.MkdirAll(filepath.Dir(to.Latest), 0o755); err != nil {
		return err
	}
	return os.WriteFile(to.Latest, data, 0o644)
}
