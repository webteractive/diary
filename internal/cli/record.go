package cli

import (
	"fmt"
	"io"
	"strings"

	"diary/internal/project"
	"diary/internal/storage"

	"github.com/spf13/cobra"
)

func (a app) recordCommand() *cobra.Command {
	var projectName string
	var recordType string
	var harness string
	var files []string
	var refs []string
	var tags []string
	var root string

	cmd := &cobra.Command{
		Use:   "record [message]",
		Short: "Record compact implementation context",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			message := strings.TrimSpace(strings.Join(args, " "))
			if message == "" {
				data, err := io.ReadAll(a.in)
				if err != nil {
					return err
				}
				message = strings.TrimSpace(string(data))
			}
			if message == "" {
				return fmt.Errorf("message is required")
			}

			resolved, err := project.Resolve(project.Options{Project: projectName})
			if err != nil {
				return err
			}

			store, err := storage.ResolveStore(storage.StoreOptions{
				Resolution:   resolved,
				RootOverride: root,
			})
			if err != nil {
				return err
			}

			record, err := storage.CreateRecord(storage.CreateRecordOptions{
				Project: resolved.Name,
				Paths:   store.Paths,
				Message: message,
				Type:    recordType,
				Harness: harness,
				Files:   files,
				Refs:    refs,
				Tags:    tags,
			})
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(a.out, "%s %s\n", record.ID, record.Hash)
			return err
		},
	}

	cmd.Flags().StringVar(&projectName, "project", "", "project name")
	cmd.Flags().StringVar(&recordType, "type", "context", "record type")
	cmd.Flags().StringVar(&harness, "harness", "unknown", "harness name")
	cmd.Flags().StringArrayVar(&files, "file", nil, "file reference")
	cmd.Flags().StringArrayVar(&refs, "ref", nil, "context reference")
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "tag")
	cmd.Flags().StringVar(&root, "root", "", "Diary storage root")

	return cmd
}
