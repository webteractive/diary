package cli

import (
	"encoding/json"
	"fmt"

	"diary/internal/project"
	"diary/internal/storage"

	"github.com/spf13/cobra"
)

func (a app) migrateCommand() *cobra.Command {
	var projectName string
	var from string
	var to string
	var force bool
	var deleteSource bool
	var dryRun bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate Diary records between storage locations",
		RunE: func(cmd *cobra.Command, args []string) error {
			resolved, err := project.Resolve(project.Options{Project: projectName})
			if err != nil {
				return err
			}

			result, err := storage.Migrate(storage.MigrateOptions{
				Resolution:   resolved,
				From:         from,
				To:           to,
				Force:        force,
				DeleteSource: deleteSource,
				DryRun:       dryRun,
			})
			if err != nil {
				return err
			}

			if jsonOutput {
				encoder := json.NewEncoder(a.out)
				encoder.SetIndent("", "  ")
				return encoder.Encode(result)
			}

			action := "migrated"
			if result.DryRun {
				action = "would migrate"
			}
			if _, err := fmt.Fprintf(a.out, "%s %d records from %s to %s\n", action, result.Records, result.FromPath, result.ToPath); err != nil {
				return err
			}
			if result.DeleteSource {
				_, err = fmt.Fprintf(a.out, "deleted source: %s\n", result.FromPath)
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&projectName, "project", "", "project name")
	cmd.Flags().StringVar(&from, "from", "", "source storage: project, user, or path")
	cmd.Flags().StringVar(&to, "to", "", "destination storage: project, user, or path")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing destination records")
	cmd.Flags().BoolVar(&deleteSource, "delete-source", false, "delete source project records after migration")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview migration without writing files")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit JSON")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")

	return cmd
}
