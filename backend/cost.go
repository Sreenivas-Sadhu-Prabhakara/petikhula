package backend

import "fmt"

// Input holds the carton-vs-loose decision inputs for one SKU.
type Input struct {
	CartonPricePerUnit float64 `json:"cartonPricePerUnit"`
	LoosePricePerUnit  float64 `json:"loosePricePerUnit"`
	CartonSize         int     `json:"cartonSize"`         // units per sealed carton
	SpoilagePct        float64 `json:"spoilagePct"`        // expected breakage/spoilage % on a carton
	CapitalCostRatePct float64 `json:"capitalCostRatePct"` // annual cost of tied-up capital, %
	WeeklySellThrough  float64 `json:"weeklySellThrough"`  // units sold per week
}

// Result is the carton-vs-loose comparison.
type Result struct {
	CartonCostPerUsableUnit float64 `json:"cartonCostPerUsableUnit"`
	LooseCostPerUsableUnit  float64 `json:"looseCostPerUsableUnit"`
	HoldingCostPerUnit      float64 `json:"holdingCostPerUnit"`
	WeeksToSellCarton       float64 `json:"weeksToSellCarton"`
	Winner                  string  `json:"winner"`             // "carton" | "loose" | "tie"
	BreakEvenWeeklyOfftake  float64 `json:"breakEvenWeeklyOfftake"` // sell-through where carton==loose
}

// Validate reports whether the Input is well formed.
func (in Input) Validate() error {
	if in.CartonPricePerUnit < 0 || in.LoosePricePerUnit < 0 {
		return fmt.Errorf("prices cannot be negative")
	}
	if in.CartonSize <= 0 {
		return fmt.Errorf("carton size must be positive")
	}
	if in.SpoilagePct < 0 || in.SpoilagePct >= 100 {
		return fmt.Errorf("spoilage %% must be between 0 and 100")
	}
	if in.CapitalCostRatePct < 0 {
		return fmt.Errorf("capital cost rate cannot be negative")
	}
	if in.WeeklySellThrough <= 0 {
		return fmt.Errorf("weekly sell-through must be positive")
	}
	return nil
}

// TrueCostPerUsableUnit nets spoilage and capital-holding cost into a real cost
// per usable unit for the carton, versus the loose price, and finds the weekly
// off-take at which the two break even.
func TrueCostPerUsableUnit(in Input) Result {
	spoilageFactor := 1 - in.SpoilagePct/100
	cartonSpoilageAdj := in.CartonPricePerUnit / spoilageFactor

	weeksToSell := float64(in.CartonSize) / in.WeeklySellThrough
	// Average capital tied over the sell-down period is ~half the carton value.
	holding := in.CartonPricePerUnit * (in.CapitalCostRatePct / 100) * (weeksToSell / 52) / 2
	cartonCost := cartonSpoilageAdj + holding
	looseCost := in.LoosePricePerUnit

	winner := "tie"
	if cartonCost < looseCost {
		winner = "carton"
	} else if looseCost < cartonCost {
		winner = "loose"
	}

	// Break-even weekly off-take: solve cartonSpoilageAdj + k/W = loose for W,
	// where k = CartonPricePerUnit*(rate/100)*(CartonSize/52)/2.
	breakEven := 0.0
	gap := looseCost - cartonSpoilageAdj
	k := in.CartonPricePerUnit * (in.CapitalCostRatePct / 100) * (float64(in.CartonSize) / 52) / 2
	if gap > 0 && k > 0 {
		breakEven = k / gap
	}

	return Result{
		CartonCostPerUsableUnit: cartonCost,
		LooseCostPerUsableUnit:  looseCost,
		HoldingCostPerUnit:      holding,
		WeeksToSellCarton:       weeksToSell,
		Winner:                  winner,
		BreakEvenWeeklyOfftake:  breakEven,
	}
}
