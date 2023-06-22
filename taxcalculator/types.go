package taxcalculator

import (
	"github.com/ybakhan/tax_calculator/taxclient"
)

type TaxCalculation struct {
	TotalTaxes    string       `json:"total_taxes"`
	EffectiveRate string       `json:"effective_rate"`
	TaxByBand     []*TaxByBand `json:"taxes_per_band"`
}

type TaxByBand struct {
	Tax  string                `json:"tax"`
	Band *taxclient.TaxBracket `json:"band"`
}
