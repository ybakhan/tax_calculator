package taxcalculator

import (
	"context"

	"github.com/ybakhan/tax_calculator/taxclient"
)

type TaxCalculator interface {
	Calculate(ctx context.Context, year, salary string) (*TaxCalculation, error)
}

type taxCalculator struct {
	TaxClient taxclient.TaxClient
}

type TaxCalculation struct {
	TotalTaxes    string `json:"total_taxes"`
	EffectiveRate string `json:"effective_rate"`
}

type TaxByBand struct {
	Tax  string   	 			`json:"tax"`
	Band taxclient.TaxBracket	`json:"band"`
}