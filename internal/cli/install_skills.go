package cli

import (
	"encoding/json"
	"fmt"

	"diary/internal/install"

	"github.com/spf13/cobra"
)

func (a app) installSkillsCommand() *cobra.Command {
	var target string
	var path string
	var force bool
	var dryRun bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "install-skills",
		Short: "Install Diary usage skills for supported AI harnesses",
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := install.Install(install.Options{
				Target: install.Target(target),
				Path:   path,
				Force:  force,
				DryRun: dryRun,
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
				} else if result.Overwritten {
					action = "overwrote"
				}
				if _, err := fmt.Fprintf(a.out, "%s %s %s skill: %s\n", action, result.Target, result.Skill, result.Path); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "target harness: codex, claude, or all")
	cmd.Flags().StringVar(&path, "path", "", "custom destination directory for SKILL.md")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing skill")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview installation without writing files")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit JSON")
	_ = cmd.MarkFlagRequired("target")

	return cmd
}
