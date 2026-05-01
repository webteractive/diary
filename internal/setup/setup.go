package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"diary/internal/install"
	"diary/internal/project"
)

type Scope string

const (
	ScopeGlobal  Scope = "global"
	ScopeProject Scope = "project"
)

const (
	blockStart = "<diary>"
	blockEnd   = "</diary>"
)

type Options struct {
	Target        install.Target
	Scope         Scope
	InstallSkills bool
	Force         bool
	DryRun        bool
	WorkDir       string
}

type Result struct {
	Target      install.Target `json:"target"`
	Scope       Scope          `json:"scope,omitempty"`
	Kind        string         `json:"kind"`
	Name        string         `json:"name"`
	Path        string         `json:"path"`
	Installed   bool           `json:"installed"`
	DryRun      bool           `json:"dry_run"`
	Overwritten bool           `json:"overwritten"`
	Unchanged   bool           `json:"unchanged"`
}

func Init(opts Options) ([]Result, error) {
	if opts.Scope == "" {
		opts.Scope = ScopeGlobal
	}
	if err := validateScope(opts.Scope); err != nil {
		return nil, err
	}

	targets, err := install.ExpandTargets(opts.Target)
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(targets))
	for _, target := range targets {
		result, err := initInstruction(target, opts)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if opts.InstallSkills {
		skillResults, err := install.Install(install.Options{
			Target: opts.Target,
			Force:  opts.Force,
			DryRun: opts.DryRun,
		})
		if err != nil {
			return nil, err
		}
		for _, skillResult := range skillResults {
			results = append(results, Result{
				Target:      skillResult.Target,
				Kind:        "skill",
				Name:        skillResult.Skill,
				Path:        skillResult.Path,
				Installed:   skillResult.Installed,
				DryRun:      skillResult.DryRun,
				Overwritten: skillResult.Overwritten,
			})
		}
	}

	return results, nil
}

func initInstruction(target install.Target, opts Options) (Result, error) {
	path, err := instructionPath(target, opts.Scope, opts.WorkDir)
	if err != nil {
		return Result{}, err
	}

	result := Result{
		Target: target,
		Scope:  opts.Scope,
		Kind:   "instruction",
		Name:   "diary",
		Path:   path,
		DryRun: opts.DryRun,
	}

	content := instructionBlock()
	existing, readErr := os.ReadFile(path)
	if readErr != nil && !os.IsNotExist(readErr) {
		return Result{}, readErr
	}

	if opts.DryRun {
		result.Installed = true
		return result, nil
	}

	next, status, err := mergeInstruction(string(existing), content, readErr == nil, opts.Force)
	if err != nil {
		return Result{}, fmt.Errorf("%s: %w", path, err)
	}
	if status == "unchanged" {
		result.Unchanged = true
		return result, nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return Result{}, err
	}
	if err := os.WriteFile(path, []byte(next), 0o644); err != nil {
		return Result{}, err
	}

	result.Installed = true
	result.Overwritten = status == "overwritten"
	return result, nil
}

func validateScope(scope Scope) error {
	switch scope {
	case ScopeGlobal, ScopeProject:
		return nil
	default:
		return fmt.Errorf("unsupported scope: %s", scope)
	}
}

func instructionPath(target install.Target, scope Scope, workDir string) (string, error) {
	if err := install.ValidateTarget(target); err != nil {
		return "", err
	}

	filename := instructionFilename(target)
	if scope == ScopeProject {
		resolved, err := project.Resolve(project.Options{WorkDir: workDir})
		if err != nil {
			return "", err
		}
		return filepath.Join(resolved.Root, filename), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch target {
	case install.TargetCodex:
		return filepath.Join(home, ".codex", filename), nil
	case install.TargetClaude:
		return filepath.Join(home, ".claude", filename), nil
	default:
		return "", fmt.Errorf("unsupported target: %s", target)
	}
}

func instructionFilename(target install.Target) string {
	if target == install.TargetClaude {
		return "CLAUDE.md"
	}
	return "AGENTS.md"
}

func instructionBlock() string {
	return blockStart + "\n" +
		"## Diary\n\n" +
		"Before ending a meaningful coding session, ask the user whether to run the `diary-record` skill to record the latest changes.\n\n" +
		"Diary records should be managed through the `diary` CLI. Run `diary get`, `diary list`, and `diary record` from the project directory so Diary can resolve the current project and use the user's configured Diary storage location, including a private Diary repository when configured.\n" +
		blockEnd + "\n"
}

func mergeInstruction(existing, block string, exists bool, force bool) (string, string, error) {
	if !exists || strings.TrimSpace(existing) == "" {
		return block, "installed", nil
	}

	start := strings.Index(existing, blockStart)
	end := strings.Index(existing, blockEnd)
	if start >= 0 && end >= 0 && end > start {
		end += len(blockEnd)
		if end < len(existing) && existing[end] == '\n' {
			end++
		}

		current := existing[start:end]
		if current == block {
			return existing, "unchanged", nil
		}
		if !force {
			return "", "", fmt.Errorf("Diary instructions already exist; pass --force to overwrite")
		}

		return existing[:start] + block + existing[end:], "overwritten", nil
	}
	if start >= 0 || end >= 0 {
		return "", "", fmt.Errorf("partial Diary instruction block found; fix the file or pass --force after restoring the block markers")
	}

	separator := "\n\n"
	if strings.HasSuffix(existing, "\n\n") {
		separator = ""
	} else if strings.HasSuffix(existing, "\n") {
		separator = "\n"
	}

	return existing + separator + block, "installed", nil
}
