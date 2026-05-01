package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	diaryhash "diary/internal/hash"
)

type CreateRecordOptions struct {
	Root    string
	Project string
	Paths   Paths
	Message string
	Type    string
	Harness string
	Files   []string
	Refs    []string
	Tags    []string
	Now     time.Time
}

func CreateRecord(opts CreateRecordOptions) (Record, error) {
	if strings.TrimSpace(opts.Message) == "" {
		return Record{}, fmt.Errorf("message is required")
	}
	if opts.Type == "" {
		opts.Type = "context"
	}
	if opts.Harness == "" {
		opts.Harness = "unknown"
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}

	paths := opts.Paths
	if paths.Project == "" {
		paths = NewPaths(opts.Root, opts.Project)
	}
	if err := os.MkdirAll(paths.RecordsDir, 0o755); err != nil {
		return Record{}, err
	}
	projectName := opts.Project
	if projectName == "" {
		projectName = paths.Project
	}

	record := Record{
		ID:        makeID(opts.Now, opts.Harness),
		Project:   projectName,
		Type:      opts.Type,
		Timestamp: opts.Now.UTC().Format(time.RFC3339),
		Harness:   opts.Harness,
		Files:     opts.Files,
		Refs:      opts.Refs,
		Tags:      opts.Tags,
		Body:      strings.TrimSpace(opts.Message),
	}

	hash, err := diaryhash.Content(record.hashFields(), record.Body)
	if err != nil {
		return Record{}, err
	}
	record.Hash = hash

	data, err := RenderRecord(record)
	if err != nil {
		return Record{}, err
	}

	record.Path = filepath.Join(paths.RecordsDir, record.ID+".md")
	if err := os.WriteFile(record.Path, data, 0o644); err != nil {
		return Record{}, err
	}
	if err := os.WriteFile(paths.Latest, data, 0o644); err != nil {
		return Record{}, err
	}
	if err := RebuildIndex(paths); err != nil {
		return Record{}, err
	}

	return record, nil
}

func (record Record) hashFields() map[string]any {
	return map[string]any{
		"id":          record.ID,
		"project":     record.Project,
		"parent_hash": record.ParentHash,
		"type":        record.Type,
		"timestamp":   record.Timestamp,
		"harness":     record.Harness,
		"files":       record.Files,
		"refs":        record.Refs,
		"tags":        record.Tags,
	}
}

func makeID(now time.Time, harness string) string {
	return fmt.Sprintf("%s-%s-%s", now.UTC().Format("2006-01-02T150405Z"), slug(harness), randomSuffix())
}

func randomSuffix() string {
	var bytes [3]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "000000"
	}
	return hex.EncodeToString(bytes[:])
}

func slug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if ok {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
