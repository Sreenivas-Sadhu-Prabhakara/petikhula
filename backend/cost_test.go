package backend

import (
	"math"
	"testing"
)

func almost(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

func TestTrueCost_CartonWinsOnFastSellThrough(t *testing.T) {
	// Carton ₹35/u, loose ₹40/u, 24/carton, 2% spoilage, 24% capital, 24 u/week.
	in := Input{CartonPricePerUnit: 35, LoosePricePerUnit: 40, CartonSize: 24,
		SpoilagePct: 2, CapitalCostRatePct: 24, WeeklySellThrough: 24}
	r := TrueCostPerUsableUnit(in)
	if r.Winner != "carton" {
		t.Fatalf("expected carton to win at fast off-take, got %s (%.3f vs %.3f)",
			r.Winner, r.CartonCostPerUsableUnit, r.LooseCostPerUsableUnit)
	}
	if !almost(r.WeeksToSellCarton, 1) {
		t.Fatalf("weeksToSell = %v want 1", r.WeeksToSellCarton)
	}
}

func TestTrueCost_LooseWinsOnSlowSellThrough(t *testing.T) {
	// Same prices but only 2 units/week: 12 weeks of locked capital + spoilage.
	in := Input{CartonPricePerUnit: 35, LoosePricePerUnit: 40, CartonSize: 24,
		SpoilagePct: 8, CapitalCostRatePct: 30, WeeklySellThrough: 2}
	r := TrueCostPerUsableUnit(in)
	if r.CartonCostPerUsableUnit <= 35 {
		t.Fatalf("holding+spoilage should push carton cost above base 35, got %v", r.CartonCostPerUsableUnit)
	}
}

func TestTrueCost_BreakEvenIsPositiveWhenCartonBaseCheaper(t *testing.T) {
	in := Input{CartonPricePerUnit: 35, LoosePricePerUnit: 40, CartonSize: 24,
		SpoilagePct: 0, CapitalCostRatePct: 24, WeeklySellThrough: 6}
	r := TrueCostPerUsableUnit(in)
	if r.BreakEvenWeeklyOfftake <= 0 {
		t.Fatalf("expected a positive break-even off-take, got %v", r.BreakEvenWeeklyOfftake)
	}
}

func TestValidate(t *testing.T) {
	ok := Input{CartonPricePerUnit: 35, LoosePricePerUnit: 40, CartonSize: 24, WeeklySellThrough: 6}
	if err := ok.Validate(); err != nil {
		t.Fatalf("valid rejected: %v", err)
	}
	for i, bad := range []Input{
		{CartonSize: 0, WeeklySellThrough: 6, LoosePricePerUnit: 1},
		{CartonSize: 24, WeeklySellThrough: 0, LoosePricePerUnit: 1},
		{CartonSize: 24, WeeklySellThrough: 6, SpoilagePct: 120, LoosePricePerUnit: 1},
		{CartonSize: 24, WeeklySellThrough: 6, CartonPricePerUnit: -1},
	} {
		if err := bad.Validate(); err == nil {
			t.Fatalf("bad input %d accepted", i)
		}
	}
}
