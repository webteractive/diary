package cli

import (
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const banner = `
     ___
    / _ \  diary
   / /_\ \ local AI context
   \  _  / records that remember
    \_/\/
`

type app struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func Execute() error {
	return New(os.Stdin, os.Stdout, os.Stderr).Execute()
}

func New(in io.Reader, out, errOut io.Writer) *cobra.Command {
	a := app{in: in, out: out, err: errOut}
	root := &cobra.Command{
		Use:          "diary",
		Short:        "Record and retrieve local AI harness context",
		Long:         strings.Trim(banner, "\n") + "\n\nRecord and retrieve local AI harness context.",
		Version:      Version,
		SilenceUsage: true,
	}
	root.SetVersionTemplate("diary version {{.Version}}\n")
	root.SetIn(in)
	root.SetOut(out)
	root.SetErr(errOut)
	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(a.recordCommand())
	root.AddCommand(a.getCommand())
	root.AddCommand(a.listCommand())
	root.AddCommand(a.renameCommand())
	root.AddCommand(a.initCommand())
	root.AddCommand(a.migrateCommand())
	root.AddCommand(a.installSkillsCommand())
	root.AddCommand(a.selfUpdateCommand())

	return root
}
