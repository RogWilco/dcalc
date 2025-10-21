package calculator

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"unicode"
)

// Result captures the outcome of evaluating an expression.
type Result struct {
	Value  Measurement
	Layout *Layout
}

var (
	errUnexpectedToken = errors.New("calculator: unexpected token in expression")
	errMismatchedParen = errors.New("calculator: mismatched parentheses")
	errDivideByZero    = errors.New("calculator: divide by zero")
	errUnitMismatch    = errors.New("calculator: unit mismatch for operation")
)

type tokenKind int

const (
	tokenOperand tokenKind = iota
	tokenOperator
	tokenLParen
	tokenRParen
)

type token struct {
	kind  tokenKind
	value string
}

// EvaluateExpression parses and evaluates an arithmetic expression returning the resulting measurement.
func EvaluateExpression(expr string, _ Options) (Result, error) {
	toks, err := tokenize(expr)
	if err != nil {
		return Result{}, err
	}
	if len(toks) == 0 {
		return Result{}, fmt.Errorf("calculator: empty expression")
	}

	rpn, err := toRPN(toks)
	if err != nil {
		return Result{}, err
	}

	res, layout, err := evalRPN(rpn)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Value:  res,
		Layout: layout,
	}, nil
}

func tokenize(expr string) ([]token, error) {
	var toks []token
	var buf strings.Builder
	runes := []rune(expr)

	flush := func() {
		if buf.Len() == 0 {
			return
		}
		toks = append(toks, token{
			kind:  tokenOperand,
			value: buf.String(),
		})
		buf.Reset()
	}

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case unicode.IsSpace(r):
			flush()
		case r == '-' && isBinarySurrounded(runes, i):
			flush()
			toks = append(toks, token{kind: tokenOperator, value: "-"})
		case r == '/' && isBinarySurrounded(runes, i):
			flush()
			toks = append(toks, token{kind: tokenOperator, value: string(r)})
		case r == '+' || r == '*':
			flush()
			toks = append(toks, token{kind: tokenOperator, value: string(r)})
		case r == '/' || r == '-':
			// Treat as part of operand (fractions, mixed numbers).
			buf.WriteRune(r)
		case r == '(':
			flush()
			toks = append(toks, token{kind: tokenLParen, value: string(r)})
		case r == ')':
			flush()
			toks = append(toks, token{kind: tokenRParen, value: string(r)})
		default:
			buf.WriteRune(r)
		}
	}
	flush()

	return toks, nil
}

func isBinaryMinus(runes []rune, idx int) bool {
	return isBinarySurrounded(runes, idx)
}

func isBinarySurrounded(runes []rune, idx int) bool {
	if idx <= 0 || idx >= len(runes)-1 {
		return false
	}
	prev := runes[idx-1]
	next := runes[idx+1]
	return unicode.IsSpace(prev) && unicode.IsSpace(next)
}

func toRPN(tokens []token) ([]token, error) {
	var output []token
	var stack []token

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	for _, tk := range tokens {
		switch tk.kind {
		case tokenOperand:
			output = append(output, tk)
		case tokenOperator:
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.kind != tokenOperator {
					break
				}
				topPrec := precedence[top.value]
				curPrec := precedence[tk.value]
				if topPrec >= curPrec {
					output = append(output, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, tk)
		case tokenLParen:
			stack = append(stack, tk)
		case tokenRParen:
			found := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.kind == tokenLParen {
					found = true
					break
				}
				output = append(output, top)
			}
			if !found {
				return nil, errMismatchedParen
			}
		default:
			return nil, errUnexpectedToken
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top.kind == tokenLParen || top.kind == tokenRParen {
			return nil, errMismatchedParen
		}
		output = append(output, top)
	}

	return output, nil
}

func evalRPN(tokens []token) (Measurement, *Layout, error) {
	var stack []Measurement
	var layout *Layout

	for _, tk := range tokens {
		switch tk.kind {
		case tokenOperand:
			operand, err := parseOperand(tk.value)
			if err != nil {
				return Measurement{}, nil, err
			}
			stack = append(stack, operand)
		case tokenOperator:
			if len(stack) < 2 {
				return Measurement{}, nil, errUnexpectedToken
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch tk.value {
			case "+":
				if left.UnitExp != right.UnitExp {
					return Measurement{}, nil, errUnitMismatch
				}
				stack = append(stack, Measurement{
					Value:   new(big.Rat).Add(left.Value, right.Value),
					UnitExp: left.UnitExp,
				})
			case "-":
				if left.UnitExp != right.UnitExp {
					return Measurement{}, nil, errUnitMismatch
				}
				stack = append(stack, Measurement{
					Value:   new(big.Rat).Sub(left.Value, right.Value),
					UnitExp: left.UnitExp,
				})
			case "*":
				stack = append(stack, Measurement{
					Value:   new(big.Rat).Mul(left.Value, right.Value),
					UnitExp: left.UnitExp + right.UnitExp,
				})
			case "/":
				if right.Value.Sign() == 0 {
					return Measurement{}, nil, errDivideByZero
				}
				res := Measurement{
					Value:   new(big.Rat).Quo(left.Value, right.Value),
					UnitExp: left.UnitExp - right.UnitExp,
				}
				stack = append(stack, res)
				if left.UnitExp == right.UnitExp && left.UnitExp != 0 {
					layout = computeLayout(left, right)
				}
			default:
				return Measurement{}, nil, errUnexpectedToken
			}
		default:
			return Measurement{}, nil, errUnexpectedToken
		}
	}

	if len(stack) != 1 {
		return Measurement{}, nil, errUnexpectedToken
	}

	result := stack[0]
	if result.Value == nil {
		result.Value = new(big.Rat)
	}

	if layout != nil && result.UnitExp != 0 {
		layout = nil
	}

	return Measurement{
		Value:   cloneRat(result.Value),
		UnitExp: result.UnitExp,
	}, layout, nil
}

func parseOperand(val string) (Measurement, error) {
	candidate := strings.TrimSpace(val)
	if candidate == "" {
		return Measurement{}, fmt.Errorf("calculator: empty operand")
	}

	if isScalarLiteral(candidate) {
		r := new(big.Rat)
		if _, ok := r.SetString(candidate); !ok {
			return Measurement{}, fmt.Errorf("calculator: invalid scalar %q", candidate)
		}
		return Measurement{Value: r, UnitExp: 0}, nil
	}

	m, err := ParseMeasurement(candidate)
	if err != nil {
		return Measurement{}, err
	}
	return m, nil
}

func isScalarLiteral(s string) bool {
	hasDigit := false
	dotCount := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			hasDigit = true
			continue
		}
		if r == '.' {
			dotCount++
			if dotCount > 1 {
				return false
			}
			continue
		}
		return false
	}
	return hasDigit
}

func computeLayout(numerator, denominator Measurement) *Layout {
	if denominator.Value.Sign() == 0 {
		return nil
	}

	ratio := new(big.Rat).Quo(cloneRat(numerator.Value), cloneRat(denominator.Value))
	num := new(big.Int).Set(ratio.Num())
	den := new(big.Int).Set(ratio.Denom())
	if den.Sign() == 0 {
		return nil
	}
	quot := new(big.Int).Quo(num, den)

	count := int(quot.Int64())
	if count < 0 {
		// For negative counts fall back to default behavior, layout less helpful.
		return nil
	}

	multiplier := new(big.Rat).Mul(cloneRat(denominator.Value), new(big.Rat).SetInt(quot))
	remainder := new(big.Rat).Sub(cloneRat(numerator.Value), multiplier)
	remainder = normalizePositiveRemainder(remainder, cloneRat(denominator.Value))

	return &Layout{
		Count: count,
		Remainder: Measurement{
			Value:   remainder,
			UnitExp: numerator.UnitExp,
		},
	}
}

func normalizePositiveRemainder(remainder, divisor *big.Rat) *big.Rat {
	if remainder.Sign() >= 0 {
		return remainder
	}
	divAbs := cloneRat(divisor)
	if divAbs.Sign() < 0 {
		divAbs.Neg(divAbs)
	}
	for remainder.Sign() < 0 {
		remainder.Add(remainder, divAbs)
	}
	return remainder
}
