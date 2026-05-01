package install

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Options struct {
	Target Target
	Path   string
	Force  bool
	DryRun bool
	Now    time.Time
}

type Result struct {
	Target      Target `json:"target"`
	Skill       string `json:"skill"`
	Path        string `json:"path"`
	Installed   bool   `json:"installed"`
	DryRun      bool   `json:"dry_run"`
	Overwritten bool   `json:"overwritten"`
}

func Install(opts Options) ([]Result, error) {
	targets, err := ExpandTargets(opts.Target)
	if err != nil {
		return nil, err
	}
	if opts.Path != "" && len(targets) > 1 {
		return nil, fmt.Errorf("--path can only be used with a single target")
	}

	results := make([]Result, 0, len(targets)*len(Templates()))
	for _, target := range targets {
		targetResults, err := installOne(target, opts)
		if err != nil {
			return nil, err
		}
		results = append(results, targetResults...)
	}

	return results, nil
}

func installOne(target Target, opts Options) ([]Result, error) {
	if err := ValidateTarget(target); err != nil {
		return nil, err
	}

	rootDir, err := destinationRoot(target, opts.Path)
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(Templates()))
	for _, template := range Templates() {
		destDir := filepath.Join(rootDir, template.Name)
		destPath := filepath.Join(destDir, "SKILL.md")

		result := Result{
			Target: target,
			Skill:  template.Name,
			Path:   destPath,
			DryRun: opts.DryRun,
		}

		if opts.DryRun {
			results = append(results, result)
			continue
		}

		_, statErr := os.Stat(destPath)
		if statErr == nil && !opts.Force {
			return nil, fmt.Errorf("%s already exists; pass --force to overwrite", destPath)
		}
		if statErr != nil && !os.IsNotExist(statErr) {
			return nil, statErr
		}

		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(destPath, []byte(template.Content), 0o644); err != nil {
			return nil, err
		}

		result.Installed = true
		result.Overwritten = statErr == nil
		results = append(results, result)
	}

	return results, nil
}

func destinationRoot(target Target, customPath string) (string, error) {
	if customPath != "" {
		return filepath.Clean(customPath), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch target {
	case TargetCodex:
		return filepath.Join(home, ".codex", "skills"), nil
	case TargetClaude:
		return filepath.Join(home, ".claude", "skills"), nil
	default:
		return "", fmt.Errorf("unsupported target: %s", target)
	}
}
