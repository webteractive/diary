package update

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const DefaultRepo = "webteractive/diary"

type Options struct {
	Repo       string
	Version    string
	DryRun     bool
	Executable string
	Goos       string
	Goarch     string
	Client     *http.Client
}

type Result struct {
	Version    string `json:"version"`
	Asset      string `json:"asset"`
	URL        string `json:"url"`
	Path       string `json:"path"`
	DryRun     bool   `json:"dry_run"`
	Updated    bool   `json:"updated"`
	NeedsSudo  bool   `json:"needs_sudo"`
	Downloaded bool   `json:"downloaded"`
}

type releaseResponse struct {
	TagName string `json:"tag_name"`
}

func SelfUpdate(opts Options) (Result, error) {
	if opts.Repo == "" {
		opts.Repo = DefaultRepo
	}
	if opts.Goos == "" {
		opts.Goos = runtime.GOOS
	}
	if opts.Goarch == "" {
		opts.Goarch = runtime.GOARCH
	}
	if opts.Client == nil {
		opts.Client = &http.Client{Timeout: 30 * time.Second}
	}

	if opts.Goos == "windows" {
		return Result{}, fmt.Errorf("self-update is not supported on windows yet")
	}

	executable := opts.Executable
	if executable == "" {
		var err error
		executable, err = os.Executable()
		if err != nil {
			return Result{}, err
		}
	}
	executable, err := filepath.Abs(executable)
	if err != nil {
		return Result{}, err
	}

	version := opts.Version
	if version == "" || version == "latest" {
		version, err = latestVersion(opts.Client, opts.Repo)
		if err != nil {
			return Result{}, err
		}
	}

	asset := fmt.Sprintf("diary_%s_%s_%s.tar.gz", version, opts.Goos, opts.Goarch)
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", opts.Repo, version, asset)

	result := Result{
		Version: version,
		Asset:   asset,
		URL:     url,
		Path:    executable,
		DryRun:  opts.DryRun,
	}

	if opts.DryRun {
		result.NeedsSudo = !isWritable(filepath.Dir(executable))
		return result, nil
	}

	tmpdir, err := os.MkdirTemp("", "diary-update-*")
	if err != nil {
		return Result{}, err
	}
	defer os.RemoveAll(tmpdir)

	archivePath := filepath.Join(tmpdir, asset)
	if err := download(opts.Client, url, archivePath); err != nil {
		return Result{}, err
	}
	result.Downloaded = true

	binaryPath, err := extractTarGz(archivePath, tmpdir)
	if err != nil {
		return Result{}, err
	}
	if err := os.Chmod(binaryPath, 0o755); err != nil {
		return Result{}, err
	}

	if isWritable(filepath.Dir(executable)) {
		if err := replace(binaryPath, executable); err != nil {
			return Result{}, err
		}
		result.Updated = true
		return result, nil
	}

	result.NeedsSudo = true
	if err := sudoMove(binaryPath, executable); err != nil {
		return Result{}, err
	}
	result.Updated = true
	return result, nil
}

func latestVersion(client *http.Client, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	return latestVersionFromURL(client, url)
}

func latestVersionFromURL(client *http.Client, url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("latest release request failed: %s", resp.Status)
	}

	var release releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	if release.TagName == "" {
		return "", errors.New("latest release response did not include tag_name")
	}
	return release.TagName, nil
}

func download(client *http.Client, url, dest string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func extractTarGz(archivePath, destDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if filepath.Base(header.Name) != "diary" {
			continue
		}

		path := filepath.Join(destDir, "diary")
		out, err := os.Create(path)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return "", err
		}
		if err := out.Close(); err != nil {
			return "", err
		}
		return path, nil
	}

	return "", errors.New("archive did not contain diary binary")
}

func replace(source, dest string) error {
	backup := dest + ".old"
	_ = os.Remove(backup)
	if _, err := os.Stat(dest); err == nil {
		if err := os.Rename(dest, backup); err != nil {
			return err
		}
	}
	if err := os.Rename(source, dest); err != nil {
		if _, statErr := os.Stat(backup); statErr == nil {
			_ = os.Rename(backup, dest)
		}
		return err
	}
	_ = os.Remove(backup)
	return nil
}

func isWritable(dir string) bool {
	file, err := os.CreateTemp(dir, ".diary-write-test-*")
	if err != nil {
		return false
	}
	name := file.Name()
	file.Close()
	os.Remove(name)
	return true
}

func sudoMove(source, dest string) error {
	if _, err := exec.LookPath("sudo"); err != nil {
		return fmt.Errorf("%s is not writable and sudo is unavailable", filepath.Dir(dest))
	}
	cmd := exec.Command("sudo", "mv", source, dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sudo mv failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
