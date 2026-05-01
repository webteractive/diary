package cli

import (
	"encoding/json"
	"fmt"

	"diary/internal/install"
	"diary/internal/setup"

	"github.com/spf13/cobra"
)

func (a app) initCommand() *cobra.Command {
	var target string
	var scope string
	var installSkills bool
	var force bool
	var dryRun bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Install Diary reminder instructions for supported AI harnesses",
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := setup.Init(setup.Options{
				Target:        install.Target(target),
				Scope:         setup.Scope(scope),
				InstallSkills: installSkills,
				Force:         force,
				DryRun:        dryRun,
			})
			if err != nil {
				return err
			}

			if jsonOutput {
				encoder := json.NewEncoder(a.out)
				encoder.SetIndent("", "  ")
				return encoder.Encode(results)
			}

			for _, result := range results {
				action := "installed"
				if result.DryRun {
					action = "would install"
				} else if result.Unchanged {
					action = "unchanged"
				} else if result.Overwritten {
					action = "overwrote"
				}

				if result.Kind == "skill" {
					if _, err := fmt.Fprintf(a.out, "%s %s %s skill: %s\n", action, result.Target, result.Name, result.Path); err != nil {
						return err
					}
					continue
				}

				if _, err := fmt.Fprintf(a.out, "%s %s %s instruction: %s\n", action, result.Target, result.Scope, result.Path); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "target harness: codex, claude, or all")
	cmd.Flags().StringVar(&scope, "scope", string(setup.ScopeGlobal), "instruction scope: global or project")
	cmd.Flags().BoolVar(&installSkills, "install-skills", false, "also install Diary usage skills globally")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing Diary-managed instructions or skills")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview initialization without writing files")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit JSON")
	_ = cmd.MarkFlagRequired("target")

	return cmd
}
