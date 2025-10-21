package cli

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/rogwilco/dcalc/internal/calculator"
)

func formatMeasurement(m calculator.Measurement, format, precision string) string {
	switch strings.ToLower(format) {
	case "decimal", "dec", "dec-in", "decimal-inches":
		return formatDecimal(m)
	case "feet-inches", "ft-in", "feet":
		return formatFeetInches(m)
	default:
		return formatFractional(m)
	}
}

func formatLayout(layout calculator.Layout, format, precision string) string {
	remainder := formatMeasurement(layout.Remainder, format, precision)
	return fmt.Sprintf("%d + %s remainder", layout.Count, remainder)
}

func formatFractional(m calculator.Measurement) string {
	if m.UnitExp == 0 {
		return m.Value.RatString()
	}
	sign := ""
	value := cloneRat(m.Value)
	if value.Sign() < 0 {
		sign = "-"
		value.Neg(value)
	}
	num := new(big.Int).Set(value.Num())
	den := new(big.Int).Set(value.Denom())

	whole := new(big.Int).Quo(num, den)
	rem := new(big.Int).Sub(num, new(big.Int).Mul(whole, den))

	var builder strings.Builder
	builder.WriteString(sign)
	if whole.Sign() > 0 {
		builder.WriteString(whole.String())
	}
	if rem.Sign() > 0 {
		if whole.Sign() > 0 {
			builder.WriteString("-")
		}
		reduced := new(big.Rat).SetFrac(rem, den)
		if builder.Len() > len(sign) {
			builder.WriteString(reduced.RatString())
		} else {
			builder.WriteString(reduced.RatString())
		}
	}
	if builder.Len() == len(sign) {
		if whole.Sign() == 0 && rem.Sign() == 0 {
			builder.WriteString("0")
		} else {
			builder.WriteString(whole.String())
		}
	}
	builder.WriteString(`"`)
	return builder.String()
}

func formatDecimal(m calculator.Measurement) string {
	if m.UnitExp == 0 {
		return m.Value.FloatString(6)
	}
	return fmt.Sprintf("%s\"", m.Value.FloatString(6))
}

func formatFeetInches(m calculator.Measurement) string {
	if m.UnitExp != 1 {
		return m.Value.RatString()
	}

	sign := ""
	value := cloneRat(m.Value)
	if value.Sign() < 0 {
		sign = "-"
		value.Neg(value)
	}

	num := new(big.Int).Set(value.Num())
	den := new(big.Int).Set(value.Denom())

	twelve := big.NewInt(12)
	denFeet := new(big.Int).Mul(den, twelve)
	feet := new(big.Int).Quo(num, denFeet)

	remNum := new(big.Int).Sub(num, new(big.Int).Mul(feet, denFeet))
	remainder := new(big.Rat).SetFrac(remNum, den)

	var builder strings.Builder
	builder.WriteString(sign)
	builder.WriteString(feet.String())
	builder.WriteString("' ")
	// Reuse fractional formatter but ensure we don't append another unit.
	remainderStr := formatFractional(calculator.Measurement{
		Value:   remainder,
		UnitExp: 1,
	})
	builder.WriteString(remainderStr)
	return builder.String()
}

func cloneRat(r *big.Rat) *big.Rat {
	return new(big.Rat).Set(r)
}
