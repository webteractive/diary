package cli

import (
	"fmt"
	"os"
	"time"

	"diary/internal/project"
	"diary/internal/storage"

	"github.com/spf13/cobra"
)

func (a app) renameCommand() *cobra.Command {
	var root string

	cmd := &cobra.Command{
		Use:   "rename <project-name>",
		Short: "Rename the mapped diary project for this checkout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resolved, err := project.Resolve(project.Options{Project: args[0]})
			if err != nil {
				return err
			}

			diaryRoot := root
			if diaryRoot == "" {
				diaryRoot = os.Getenv(storage.EnvRoot)
			}
			if diaryRoot == "" {
				diaryRoot, err = storage.DefaultRoot()
				if err != nil {
					return err
				}
			}

			result, err := storage.RenameProject(diaryRoot, resolved, time.Now().UTC())
			if err != nil {
				return err
			}

			if result.Old.ID == result.New.ID {
				_, err = fmt.Fprintf(a.out, "Project already named %s (%s)\n", result.New.Name, result.New.ID)
				return err
			}
			_, err = fmt.Fprintf(a.out, "Renamed project %s -> %s\n", result.Old.ID, result.New.ID)
			return err
		},
	}

	cmd.Flags().StringVar(&root, "root", "", "Diary storage root")

	return cmd
}
