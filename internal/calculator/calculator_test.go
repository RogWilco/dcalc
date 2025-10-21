package calculator

import (
	"math/big"
	"testing"
)

func TestParseMeasurement_MixedNumber(t *testing.T) {
	got, err := ParseMeasurement(`5-3/8"`)
	if err != nil {
		t.Fatalf("ParseMeasurement returned error: %v", err)
	}

	want := big.NewRat(43, 8) // 5 3/8 inches
	if got.UnitExp != 1 {
		t.Fatalf("UnitExp = %d, want 1", got.UnitExp)
	}
	if got.Value.Cmp(want) != 0 {
		t.Fatalf("Value = %s, want %s", got.Value.RatString(), want.RatString())
	}
}

func TestParseMeasurement_FeetInches(t *testing.T) {
	got, err := ParseMeasurement(`6' 11-1/2"`)
	if err != nil {
		t.Fatalf("ParseMeasurement returned error: %v", err)
	}

	want := big.NewRat(167, 2) // 6 feet 11.5 inches = 167/2 inches
	if got.UnitExp != 1 {
		t.Fatalf("UnitExp = %d, want 1", got.UnitExp)
	}
	if got.Value.Cmp(want) != 0 {
		t.Fatalf("Value = %s, want %s", got.Value.RatString(), want.RatString())
	}
}

func TestParseMeasurement_Decimal(t *testing.T) {
	got, err := ParseMeasurement(`3.125"`)
	if err != nil {
		t.Fatalf("ParseMeasurement returned error: %v", err)
	}

	want := big.NewRat(25, 8) // 3.125 inches reduced
	if got.UnitExp != 1 {
		t.Fatalf("UnitExp = %d, want 1", got.UnitExp)
	}
	if got.Value.Cmp(want) != 0 {
		t.Fatalf("Value = %s, want %s", got.Value.RatString(), want.RatString())
	}
}

func TestParseMeasurement_MixedWithoutUnit(t *testing.T) {
	got, err := ParseMeasurement(`4-5/8`)
	if err != nil {
		t.Fatalf("ParseMeasurement returned error: %v", err)
	}

	want := big.NewRat(37, 8)
	if got.UnitExp != 1 {
		t.Fatalf("UnitExp = %d, want 1", got.UnitExp)
	}
	if got.Value.Cmp(want) != 0 {
		t.Fatalf("Value = %s, want %s", got.Value.RatString(), want.RatString())
	}
}

func TestParseMeasurement_NegativeMixedNumber(t *testing.T) {
	got, err := ParseMeasurement(`-2 7/16"`)
	if err != nil {
		t.Fatalf("ParseMeasurement returned error: %v", err)
	}

	want := big.NewRat(-39, 16)
	if got.UnitExp != 1 {
		t.Fatalf("UnitExp = %d, want 1", got.UnitExp)
	}
	if got.Value.Cmp(want) != 0 {
		t.Fatalf("Value = %s, want %s", got.Value.RatString(), want.RatString())
	}
}

func TestEvaluateExpression_Addition(t *testing.T) {
	res, err := EvaluateExpression(`1-1/2 + 4-5/8`, Options{})
	if err != nil {
		t.Fatalf("EvaluateExpression returned error: %v", err)
	}

	want := big.NewRat(49, 8) // 6 1/8 inches
	if res.Value.UnitExp != 1 {
		t.Fatalf("UnitExp = %d, want 1", res.Value.UnitExp)
	}
	if res.Value.Value.Cmp(want) != 0 {
		t.Fatalf("Value = %s, want %s", res.Value.Value.RatString(), want.RatString())
	}
	if res.Layout != nil {
		t.Fatalf("Layout = %#v, want nil", res.Layout)
	}
}

func TestEvaluateExpression_DivisionProducesLayout(t *testing.T) {
	res, err := EvaluateExpression(`43-5/8 / 4-7/8`, Options{})
	if err != nil {
		t.Fatalf("EvaluateExpression returned error: %v", err)
	}

	want := big.NewRat(349, 39)
	if res.Value.UnitExp != 0 {
		t.Fatalf("UnitExp = %d, want 0 for dimensionless ratio", res.Value.UnitExp)
	}
	if res.Value.Value.Cmp(want) != 0 {
		t.Fatalf("Value = %s, want %s", res.Value.Value.RatString(), want.RatString())
	}
	if res.Layout == nil {
		t.Fatalf("Layout is nil, want quotient/remainder details")
	}
	if res.Layout.Count != 8 {
		t.Fatalf("Layout.Count = %d, want 8", res.Layout.Count)
	}
	wantRemainder := big.NewRat(37, 8) // 4 5/8 inches
	if res.Layout.Remainder.UnitExp != 1 {
		t.Fatalf("Remainder.UnitExp = %d, want 1", res.Layout.Remainder.UnitExp)
	}
	if res.Layout.Remainder.Value.Cmp(wantRemainder) != 0 {
		t.Fatalf("Remainder.Value = %s, want %s", res.Layout.Remainder.Value.RatString(), wantRemainder.RatString())
	}
}

func TestEvaluateExpression_LengthDividedByScalarKeepsUnits(t *testing.T) {
	res, err := EvaluateExpression(`12" / 3`, Options{})
	if err != nil {
		t.Fatalf("EvaluateExpression returned error: %v", err)
	}

	want := big.NewRat(4, 1)
	if res.Value.UnitExp != 1 {
		t.Fatalf("UnitExp = %d, want 1", res.Value.UnitExp)
	}
	if res.Value.Value.Cmp(want) != 0 {
		t.Fatalf("Value = %s, want %s", res.Value.Value.RatString(), want.RatString())
	}
	if res.Layout != nil {
		t.Fatalf("Layout = %#v, want nil", res.Layout)
	}
}

func TestEvaluateExpression_SoloMeasurement(t *testing.T) {
	res, err := EvaluateExpression(`4-5/8`, Options{})
	if err != nil {
		t.Fatalf("EvaluateExpression returned error: %v", err)
	}
	if res.Value.Value.RatString() != "37/8" {
		t.Fatalf("single operand result = %s", res.Value.Value.RatString())
	}
}
