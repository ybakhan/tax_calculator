package taxcalculator

import (
	"context"
	"fmt"

	"github.com/ybakhan/tax_calculator/taxclient"
)

func InitializeTaxCalculator(tc taxclient.TaxClient) TaxCalculator {
	return &taxCalculator{tc}
}

func (tc *taxCalculator) Calculate(ctx context.Context, year, salary string) (*TaxCalculation, error) {
	brackets, resp, err := tc.TaxClient.GetBrackets(ctx, year);
	if err != nil {
		return nil, err
	}

	if (resp == taxclient.Failed || resp == taxclient.NotFound) {
		return nil, nil
	}

	for _, bracket := range brackets {
		fmt.Print(bracket)
	}

	return nil, nil
}