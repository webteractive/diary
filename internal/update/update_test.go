package update

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelfUpdateDryRunUsesProvidedVersion(t *testing.T) {
	exe := filepath.Join(t.TempDir(), "diary")
	if err := os.WriteFile(exe, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}

	result, err := SelfUpdate(Options{
		Repo:       "webteractive/diary",
		Version:    "v0.0.1",
		DryRun:     true,
		Executable: exe,
		Goos:       "darwin",
		Goarch:     "arm64",
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.Version != "v0.0.1" {
		t.Fatalf("expected version v0.0.1, got %q", result.Version)
	}
	if result.Asset != "diary_v0.0.1_darwin_arm64.tar.gz" {
		t.Fatalf("unexpected asset: %s", result.Asset)
	}
	if result.Updated || result.Downloaded {
		t.Fatalf("dry-run should not update or download: %#v", result)
	}
}

func TestSelfUpdateRejectsWindows(t *testing.T) {
	_, err := SelfUpdate(Options{
		Version:    "v0.0.1",
		DryRun:     true,
		Executable: filepath.Join(t.TempDir(), "diary.exe"),
		Goos:       "windows",
		Goarch:     "amd64",
	})
	if err == nil || !strings.Contains(err.Error(), "windows") {
		t.Fatalf("expected windows unsupported error, got %v", err)
	}
}

func TestLatestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/webteractive/diary/releases/latest" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"tag_name":"v0.0.1"}`))
	}))
	defer server.Close()

	version, err := latestVersionFromURL(server.Client(), server.URL+"/repos/webteractive/diary/releases/latest")
	if err != nil {
		t.Fatal(err)
	}
	if version != "v0.0.1" {
		t.Fatalf("expected v0.0.1, got %q", version)
	}
}
