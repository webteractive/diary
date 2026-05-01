package cli

import (
	"encoding/json"
	"path/filepath"

	"diary/internal/project"
	"diary/internal/render"
	"diary/internal/storage"

	"github.com/spf13/cobra"
)

func (a app) listCommand() *cobra.Command {
	var projectName string
	var projectsOnly bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List diary projects and records",
		RunE: func(cmd *cobra.Command, args []string) error {
			resolved, err := project.Resolve(project.Options{Project: projectName})
			if err != nil {
				return err
			}

			if projectsOnly {
				projects, err := storage.Projects(resolved.Root)
				if err != nil {
					return err
				}
				if jsonOutput {
					return json.NewEncoder(a.out).Encode(projects)
				}
				return render.ProjectsMarkdown(a.out, projects)
			}

			paths := storage.NewPaths(resolved.Root, resolved.Name)
			index, err := storage.ReadIndex(paths)
			if err != nil {
				if records, readErr := storage.ReadRecords(paths); readErr == nil {
					index = storage.Index{Project: resolved.Name}
					for _, record := range records {
						index.Records = append(index.Records, storage.IndexEntry{
							ID:        record.ID,
							Hash:      record.Hash,
							Timestamp: record.Timestamp,
							Project:   record.Project,
							Type:      record.Type,
							Preview:   storage.Preview(record.Body, 100),
							Files:     record.Files,
							Refs:      record.Refs,
							Tags:      record.Tags,
							Path:      filepath.Base(record.Path),
						})
					}
				} else {
					return err
				}
			}

			if jsonOutput {
				return render.IndexJSON(a.out, index)
			}
			return render.IndexMarkdown(a.out, index)
		},
	}

	cmd.Flags().StringVar(&projectName, "project", "", "project name")
	cmd.Flags().BoolVar(&projectsOnly, "projects", false, "list projects")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit JSON")

	return cmd
}
