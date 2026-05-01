package cli

import (
	"fmt"

	"diary/internal/project"
	"diary/internal/render"
	"diary/internal/storage"

	"github.com/spf13/cobra"
)

func (a app) getCommand() *cobra.Command {
	var projectName string
	var id string
	var hashPrefix string
	var jsonOutput bool
	var maxChars int
	var root string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve prompt-ready context",
		RunE: func(cmd *cobra.Command, args []string) error {
			resolved, err := project.Resolve(project.Options{Project: projectName})
			if err != nil {
				return err
			}
			resolveStore := storage.ResolveStore
			if projectName == "" {
				resolveStore = storage.ResolveStoreForRoot
			}
			store, err := resolveStore(storage.StoreOptions{
				Resolution:   resolved,
				RootOverride: root,
			})
			if err != nil {
				return err
			}
			paths := store.Paths

			var record storage.Record
			switch {
			case id != "":
				record, err = storage.FindByID(paths, id)
			case hashPrefix != "":
				record, err = storage.FindByHashPrefix(paths, hashPrefix)
			default:
				record, err = storage.Latest(paths)
			}
			if err != nil {
				return err
			}
			if maxChars > 0 && len(record.Body) > maxChars {
				record.Body = storage.Preview(record.Body, maxChars)
			}

			if jsonOutput {
				return render.RecordJSON(a.out, record)
			}
			return render.RecordMarkdown(a.out, record)
		},
	}

	cmd.Flags().StringVar(&projectName, "project", "", "project name")
	cmd.Flags().StringVar(&id, "id", "", "record id")
	cmd.Flags().StringVar(&hashPrefix, "hash", "", "record hash prefix")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit JSON")
	cmd.Flags().IntVar(&maxChars, "max-chars", 0, "maximum body characters")
	cmd.Flags().StringVar(&root, "root", "", "Diary storage root")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if id != "" && hashPrefix != "" {
			return fmt.Errorf("--id and --hash cannot be used together")
		}
		return nil
	}

	return cmd
}
