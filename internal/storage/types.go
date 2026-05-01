package storage

type Record struct {
	ID         string   `yaml:"id" json:"id"`
	Project    string   `yaml:"project" json:"project"`
	Hash       string   `yaml:"hash" json:"hash"`
	ParentHash string   `yaml:"parent_hash,omitempty" json:"parent_hash,omitempty"`
	Type       string   `yaml:"type" json:"type"`
	Timestamp  string   `yaml:"timestamp" json:"timestamp"`
	Harness    string   `yaml:"harness,omitempty" json:"harness,omitempty"`
	Files      []string `yaml:"files,omitempty" json:"files,omitempty"`
	Refs       []string `yaml:"refs,omitempty" json:"refs,omitempty"`
	Tags       []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	Body       string   `yaml:"-" json:"body"`
	Path       string   `yaml:"-" json:"path,omitempty"`
}

type Index struct {
	Project   string       `json:"project"`
	UpdatedAt string       `json:"updated_at"`
	Records   []IndexEntry `json:"records"`
}

type IndexEntry struct {
	ID        string   `json:"id"`
	Hash      string   `json:"hash"`
	Timestamp string   `json:"timestamp"`
	Project   string   `json:"project"`
	Type      string   `json:"type"`
	Preview   string   `json:"preview"`
	Files     []string `json:"files,omitempty"`
	Refs      []string `json:"refs,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Path      string   `json:"path"`
}
