package calculator

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

// Measurement represents a numeric value paired with a unit exponent.
// UnitExp==1 corresponds to linear measurements (inches), UnitExp==0 is dimensionless.
type Measurement struct {
	Value   *big.Rat
	UnitExp int
}

// Layout captures the result of dividing two like-dimension measurements.
type Layout struct {
	Count     int
	Remainder Measurement
}

// Options configures expression evaluation. Placeholder for future precision/format controls.
type Options struct{}

var (
	errEmptyMeasurement = errors.New("measurement: value cannot be empty")
	errInvalidFraction  = errors.New("measurement: invalid fraction component")
	errInvalidNumber    = errors.New("measurement: invalid numeric component")
)

// ParseMeasurement parses a measurement string into a normalized Measurement expressed in inches.
func ParseMeasurement(input string) (Measurement, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return Measurement{}, errEmptyMeasurement
	}

	sign := 1
	switch {
	case strings.HasPrefix(s, "+"):
		s = strings.TrimSpace(s[1:])
	case strings.HasPrefix(s, "-"):
		sign = -1
		s = strings.TrimSpace(s[1:])
	}
	if s == "" {
		return Measurement{}, errEmptyMeasurement
	}

	total := new(big.Rat)
	feetValue, rest, err := splitFeet(s)
	if err != nil {
		return Measurement{}, err
	}
	if feetValue != nil {
		feetInches := new(big.Rat).Mul(feetValue, big.NewRat(12, 1))
		total.Add(total, feetInches)
	}

	if rest != "" {
		inchValue, err := parseInchComponent(rest)
		if err != nil {
			return Measurement{}, err
		}
		total.Add(total, inchValue)
	}

	if feetValue == nil && rest == "" {
		return Measurement{}, fmt.Errorf("measurement: could not parse %q", input)
	}

	if sign < 0 {
		total.Neg(total)
	}

	return Measurement{
		Value:   cloneRat(total),
		UnitExp: 1,
	}, nil
}

func splitFeet(s string) (*big.Rat, string, error) {
	idx := strings.IndexRune(s, '\'')
	if idx == -1 {
		return nil, s, nil
	}
	feetPart := strings.TrimSpace(s[:idx])
	if feetPart == "" {
		return nil, "", fmt.Errorf("measurement: missing feet value in %q", s)
	}
	feetInt, err := strconv.Atoi(feetPart)
	if err != nil {
		return nil, "", fmt.Errorf("measurement: invalid feet value %q", feetPart)
	}
	rest := strings.TrimSpace(s[idx+1:])
	return big.NewRat(int64(feetInt), 1), rest, nil
}

func parseInchComponent(s string) (*big.Rat, error) {
	value := new(big.Rat)
	clean := strings.TrimSpace(s)
	if clean == "" {
		return value, nil
	}

	switch {
	case strings.HasSuffix(clean, `"`) || strings.HasSuffix(clean, "”"):
		clean = strings.TrimSpace(strings.TrimSuffix(clean, `"`))
	case strings.HasSuffix(strings.ToLower(clean), "in"):
		clean = strings.TrimSpace(clean[:len(clean)-2])
	}

	if clean == "" {
		return value, nil
	}

	// Normalize hyphenated mixed numbers to space-delimited pieces.
	normalized := strings.ReplaceAll(clean, "-", " ")
	parts := strings.Fields(normalized)
	if len(parts) == 0 {
		return value, nil
	}

	for _, part := range parts {
		r, err := parseNumber(part)
		if err != nil {
			return nil, err
		}
		value.Add(value, r)
	}
	return value, nil
}

func parseNumber(fragment string) (*big.Rat, error) {
	if strings.Contains(fragment, "/") {
		items := strings.Split(fragment, "/")
		if len(items) != 2 {
			return nil, errInvalidFraction
		}
		numStr := strings.TrimSpace(items[0])
		denStr := strings.TrimSpace(items[1])
		if numStr == "" || denStr == "" {
			return nil, errInvalidFraction
		}
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, errInvalidFraction
		}
		den, err := strconv.Atoi(denStr)
		if err != nil || den == 0 {
			return nil, errInvalidFraction
		}
		return big.NewRat(int64(num), int64(den)), nil
	}

	if strings.IndexFunc(fragment, func(r rune) bool {
		return !(unicode.IsDigit(r) || r == '.')
	}) != -1 {
		return nil, errInvalidNumber
	}

	r := new(big.Rat)
	if _, ok := r.SetString(fragment); !ok {
		return nil, errInvalidNumber
	}
	return r, nil
}

func cloneRat(r *big.Rat) *big.Rat {
	if r == nil {
		return nil
	}
	return new(big.Rat).Set(r)
}
