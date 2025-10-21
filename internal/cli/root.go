package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rogwilco/dcalc/internal/calculator"
)

const (
	flagFormat    = "format"
	flagPrecision = "precision"
)

// Execute runs the CLI application.
func Execute() error {
	return NewRootCommand().Execute()
}

// NewRootCommand builds the Cobra root command.
func NewRootCommand() *cobra.Command {
	var format string
	var precision string

	rootCmd := &cobra.Command{
		Use:   "dcalc [expression]",
		Short: "Dimensional calculator for architectural measurements",
		Long:  "dcalc evaluates fractional inch expressions and performs unit conversions for construction math.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			expr := strings.Join(args, " ")
			return runEvaluate(cmd, expr, format, precision)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&format, flagFormat, "f", "fractional", "output format: fractional|feet-inches|decimal")
	rootCmd.PersistentFlags().StringVar(&precision, flagPrecision, "1/16", "output precision (reserved for future use)")

	rootCmd.AddCommand(newEvalCommand(&format, &precision))
	rootCmd.AddCommand(newConvertCommand(&format, &precision))
	rootCmd.SilenceUsage = true

	return rootCmd
}

func runEvaluate(cmd *cobra.Command, expression, format, precision string) error {
	result, err := calculator.EvaluateExpression(expression, calculator.Options{})
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	fmt.Fprintln(out, formatMeasurement(result.Value, format, precision))
	if result.Layout != nil {
		layoutStr := formatLayout(*result.Layout, format, precision)
		fmt.Fprintf(out, "layout: %s\n", layoutStr)
	}
	return nil
}
