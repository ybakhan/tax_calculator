package taxcalculator

import (
	"fmt"
	"math"

	"github.com/ybakhan/tax_calculator/taxclient"
)

func Calculate(brackets []*taxclient.TaxBracket, salary float32) *TaxCalculation {
	var taxByBand []*TaxByBand
	var totalTax float32
	for _, bracket := range brackets {
		tax := calculateByBracket(bracket, salary)
		if tax != 0 {
			taxByBand = append(taxByBand, &TaxByBand{
				format(tax),
				bracket,
			})
			totalTax += tax
		}
	}

	answer := &TaxCalculation{
		format(totalTax),
		format(totalTax / salary),
		taxByBand,
	}
	return answer
}

func calculateByBracket(bracket *taxclient.TaxBracket, salary float32) float32 {
	if salary == 0 || salary <= bracket.Min {
		return 0
	}

	if salary > bracket.Min && (salary <= bracket.Max || bracket.Max == 0) {
		return round(bracket.Rate * (salary - bracket.Min))
	}

	return round(bracket.Rate * (bracket.Max - bracket.Min))
}

func round(f float32) float32 {
	return float32(math.Round(float64(f*100)) / 100)
}

func format(f float32) string {
	return fmt.Sprintf("%.2f", f)
}
