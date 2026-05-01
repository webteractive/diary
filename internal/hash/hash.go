package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
)

func Content(fields map[string]any, body string) (string, error) {
	normalized := make(map[string]any, len(fields))
	for key, value := range fields {
		if key == "hash" {
			continue
		}
		normalized[key] = value
	}

	payload, err := json.Marshal(normalized)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256([]byte(string(payload) + "\n" + normalizeBody(body)))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func Short(value string, length int) string {
	value = strings.TrimPrefix(value, "sha256:")
	if length <= 0 || len(value) <= length {
		return value
	}
	return value[:length]
}

func normalizeBody(body string) string {
	body = strings.ReplaceAll(body, "\r\n", "\n")
	body = strings.ReplaceAll(body, "\r", "\n")
	return strings.TrimSpace(body) + "\n"
}
