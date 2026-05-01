package storage

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func RenderRecord(record Record) ([]byte, error) {
	frontmatter, err := yaml.Marshal(record)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	out.WriteString("---\n")
	out.Write(frontmatter)
	out.WriteString("---\n\n")
	out.WriteString(strings.TrimSpace(record.Body))
	out.WriteString("\n")

	return out.Bytes(), nil
}

func ParseRecord(data []byte) (Record, error) {
	text := string(data)
	if !strings.HasPrefix(text, "---\n") {
		return Record{}, fmt.Errorf("missing frontmatter")
	}

	rest := strings.TrimPrefix(text, "---\n")
	parts := strings.SplitN(rest, "\n---\n", 2)
	if len(parts) != 2 {
		return Record{}, fmt.Errorf("invalid frontmatter")
	}

	var record Record
	if err := yaml.Unmarshal([]byte(parts[0]), &record); err != nil {
		return Record{}, err
	}
	record.Body = strings.TrimSpace(parts[1])

	return record, nil
}
