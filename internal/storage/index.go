package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func RebuildIndex(paths Paths) error {
	records, err := ReadRecords(paths)
	if err != nil {
		return err
	}

	index := Index{
		Project:   paths.Project,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Records:   make([]IndexEntry, 0, len(records)),
	}

	for _, record := range records {
		index.Records = append(index.Records, IndexEntry{
			ID:        record.ID,
			Hash:      record.Hash,
			Timestamp: record.Timestamp,
			Project:   record.Project,
			Type:      record.Type,
			Preview:   Preview(record.Body, 100),
			Files:     record.Files,
			Refs:      record.Refs,
			Tags:      record.Tags,
			Path:      record.Path,
		})
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(paths.Index, data, 0o644)
}

func ReadIndex(paths Paths) (Index, error) {
	data, err := os.ReadFile(paths.Index)
	if err != nil {
		return Index{}, err
	}
	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return Index{}, err
	}
	return index, nil
}

func ReadRecords(paths Paths) ([]Record, error) {
	entries, err := os.ReadDir(paths.RecordsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	records := make([]Record, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		path := filepath.Join(paths.RecordsDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		record, err := ParseRecord(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		record.Path = path
		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp > records[j].Timestamp
	})

	return records, nil
}

func FindByID(paths Paths, id string) (Record, error) {
	records, err := ReadRecords(paths)
	if err != nil {
		return Record{}, err
	}
	for _, record := range records {
		if record.ID == id {
			return record, nil
		}
	}
	return Record{}, fmt.Errorf("record not found: %s", id)
}

func FindByHashPrefix(paths Paths, prefix string) (Record, error) {
	records, err := ReadRecords(paths)
	if err != nil {
		return Record{}, err
	}
	prefix = strings.TrimPrefix(prefix, "sha256:")
	var matches []Record
	for _, record := range records {
		if strings.HasPrefix(strings.TrimPrefix(record.Hash, "sha256:"), prefix) {
			matches = append(matches, record)
		}
	}
	if len(matches) == 0 {
		return Record{}, fmt.Errorf("record not found for hash prefix: %s", prefix)
	}
	if len(matches) > 1 {
		ids := make([]string, 0, len(matches))
		for _, match := range matches {
			ids = append(ids, match.ID)
		}
		return Record{}, fmt.Errorf("hash prefix is ambiguous; matches: %s", strings.Join(ids, ", "))
	}
	return matches[0], nil
}

func Latest(paths Paths) (Record, error) {
	data, err := os.ReadFile(paths.Latest)
	if err != nil {
		if os.IsNotExist(err) {
			return Record{}, fmt.Errorf("no diary records found for project %q", paths.Project)
		}
		return Record{}, err
	}
	record, err := ParseRecord(data)
	if err != nil {
		return Record{}, err
	}
	record.Path = paths.Latest
	return record, nil
}

func Projects(root string) ([]string, error) {
	projectsDir := filepath.Join(root, ".diary", "projects")
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var projects []string
	for _, entry := range entries {
		if entry.IsDir() {
			projects = append(projects, entry.Name())
		}
	}
	sort.Strings(projects)
	return projects, nil
}

func Preview(body string, max int) string {
	body = strings.Join(strings.Fields(body), " ")
	if max <= 0 || len(body) <= max {
		return body
	}
	return strings.TrimSpace(body[:max-1]) + "…"
}
