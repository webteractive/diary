package project

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Resolution struct {
	Name string
	Root string
}

type Options struct {
	Project string
	WorkDir string
}

type configFile struct {
	Project string `yaml:"project"`
}

var unsafeProjectChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func Resolve(opts Options) (Resolution, error) {
	root, err := resolveRoot(opts.WorkDir)
	if err != nil {
		return Resolution{}, err
	}

	if opts.Project != "" {
		return Resolution{Name: Sanitize(opts.Project), Root: root}, nil
	}

	if configured := readConfiguredProject(root); configured != "" {
		return Resolution{Name: Sanitize(configured), Root: root}, nil
	}

	return Resolution{Name: Sanitize(filepath.Base(root)), Root: root}, nil
}

func Sanitize(name string) string {
	name = strings.TrimSpace(name)
	name = unsafeProjectChars.ReplaceAllString(name, "-")
	name = strings.Trim(name, ".-_")
	if name == "" {
		return "default"
	}
	return strings.ToLower(name)
}

func resolveRoot(workDir string) (string, error) {
	if workDir == "" {
		var err error
		workDir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	abs, err := filepath.Abs(workDir)
	if err != nil {
		return "", err
	}

	if gitRoot := findGitRoot(abs); gitRoot != "" {
		return gitRoot, nil
	}

	return abs, nil
}

func findGitRoot(dir string) string {
	for {
		if stat, err := os.Stat(filepath.Join(dir, ".git")); err == nil && stat.IsDir() {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func readConfiguredProject(root string) string {
	data, err := os.ReadFile(filepath.Join(root, ".diary", "config.yml"))
	if err != nil {
		return ""
	}

	var cfg configFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return ""
	}

	return cfg.Project
}
