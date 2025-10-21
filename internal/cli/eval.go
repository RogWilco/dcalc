package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

func newEvalCommand(format, precision *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eval <expression>",
		Short: "Evaluate an architectural expression explicitly",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			expr := strings.Join(args, " ")
			return runEvaluate(cmd, expr, *format, *precision)
		},
	}
	return cmd
}
