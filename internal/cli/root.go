package cli

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

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
		SilenceUsage: true,
	}
	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(a.recordCommand())
	root.AddCommand(a.getCommand())
	root.AddCommand(a.listCommand())
	root.AddCommand(a.installSkillsCommand())
	root.AddCommand(a.selfUpdateCommand())

	return root
}
