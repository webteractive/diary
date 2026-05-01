package render

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	diaryhash "diary/internal/hash"
	"diary/internal/storage"
)

func RecordMarkdown(w io.Writer, record storage.Record) error {
	_, err := fmt.Fprintf(w, "# Diary Context\n\nProject: %s\nRecord: %s\nHash: %s\nTimestamp: %s\n\n%s\n",
		record.Project,
		record.ID,
		diaryhash.Short(record.Hash, 12),
		record.Timestamp,
		strings.TrimSpace(record.Body),
	)
	return err
}

func RecordJSON(w io.Writer, record storage.Record) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(record)
}

func IndexMarkdown(w io.Writer, index storage.Index) error {
	if len(index.Records) == 0 {
		_, err := fmt.Fprintln(w, "No diary records found.")
		return err
	}
	for _, entry := range index.Records {
		_, err := fmt.Fprintf(w, "%s  %s  %s  %s\n",
			entry.ID,
			diaryhash.Short(entry.Hash, 12),
			entry.Type,
			entry.Preview,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func IndexJSON(w io.Writer, index storage.Index) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(index)
}

func ProjectsMarkdown(w io.Writer, projects []string) error {
	if len(projects) == 0 {
		_, err := fmt.Fprintln(w, "No diary projects found.")
		return err
	}
	for _, project := range projects {
		if _, err := fmt.Fprintln(w, project); err != nil {
			return err
		}
	}
	return nil
}
