package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rogwilco/dcalc/internal/calculator"
)

func newConvertCommand(format, precision *string) *cobra.Command {
	var to string

	cmd := &cobra.Command{
		Use:   "convert <measurement>",
		Short: "Convert a measurement into another unit representation",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := strings.Join(args, " ")
			measurement, err := calculator.ParseMeasurement(input)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			switch to {
			case "", "fractional":
				fmt.Fprintln(out, formatMeasurement(measurement, *format, *precision))
			case "feet-inches", "ft-in":
				fmt.Fprintln(out, formatMeasurement(measurement, "feet-inches", *precision))
			case "decimal", "dec-in":
				fmt.Fprintln(out, formatMeasurement(measurement, "decimal", *precision))
			default:
				return fmt.Errorf("unknown target format %q", to)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&to, "to", "", "output format (fractional|feet-inches|decimal)")
	return cmd
}
