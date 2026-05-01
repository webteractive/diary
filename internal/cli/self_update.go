package cli

import (
	"encoding/json"
	"fmt"

	"diary/internal/update"

	"github.com/spf13/cobra"
)

func (a app) selfUpdateCommand() *cobra.Command {
	var version string
	var repo string
	var dryRun bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Update diary from GitHub releases",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := update.SelfUpdate(update.Options{
				Repo:    repo,
				Version: version,
				DryRun:  dryRun,
			})
			if err != nil {
				return err
			}

			if jsonOutput {
				encoder := json.NewEncoder(a.out)
				encoder.SetIndent("", "  ")
				return encoder.Encode(result)
			}

			if dryRun {
				_, err = fmt.Fprintf(a.out, "would update %s from %s\n", result.Path, result.URL)
				return err
			}

			_, err = fmt.Fprintf(a.out, "updated %s to %s\n", result.Path, result.Version)
			return err
		},
	}

	cmd.Flags().StringVar(&version, "version", "latest", "version to install, for example v0.0.1")
	cmd.Flags().StringVar(&repo, "repo", update.DefaultRepo, "GitHub repository")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview update without replacing the binary")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit JSON")

	return cmd
}
