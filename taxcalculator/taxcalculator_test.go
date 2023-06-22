package taxcalculator

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ybakhan/tax_calculator/taxclient"
)

func TestCalculate(t *testing.T) {
	file, _ := ioutil.ReadFile("./input/taxbrackets.json")
	var taxBrackets taxclient.TaxBrackets
	err := json.Unmarshal([]byte(file), &taxBrackets)
	if err != nil {
		t.Errorf("Error parsing test input file: %v", err)
	}

	tests := map[string]struct {
		Salary         float32
		TaxCalculation *TaxCalculation
	}{
		"calculate over one band": {
			50196,
			&TaxCalculation{
				"7529.40",
				"0.15",
				[]*TaxByBand{
					{"7529.40", taxBrackets.Data[0]},
				},
			},
		},
		"calculate over one band boundary": {
			50197,
			&TaxCalculation{
				"7529.55",
				"0.15",
				[]*TaxByBand{
					{"7529.55", taxBrackets.Data[0]},
				},
			},
		},
		"calculate over two bands": {
			55000,
			&TaxCalculation{
				"8514.17",
				"0.15",
				[]*TaxByBand{
					{"7529.55", taxBrackets.Data[0]},
					{"984.62", taxBrackets.Data[1]},
				},
			},
		},
		"calculate over two bands boundary": {
			100392,
			&TaxCalculation{
				"17819.52",
				"0.18",
				[]*TaxByBand{
					{"7529.55", taxBrackets.Data[0]},
					{"10289.97", taxBrackets.Data[1]},
				},
			},
		},
		"calculate over three bands": {
			100393,
			&TaxCalculation{
				"17819.78",
				"0.18",
				[]*TaxByBand{
					{"7529.55", taxBrackets.Data[0]},
					{"10289.97", taxBrackets.Data[1]},
					{"0.26", taxBrackets.Data[2]},
				},
			},
		},
		"calculate over five bands": {
			221709,
			&TaxCalculation{
				"51344.50",
				"0.23",
				[]*TaxByBand{
					{"7529.55", taxBrackets.Data[0]},
					{"10289.97", taxBrackets.Data[1]},
					{"14360.58", taxBrackets.Data[2]},
					{"19164.07", taxBrackets.Data[3]},
					{"0.33", taxBrackets.Data[4]},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			taxCalculation := Calculate(taxBrackets.Data, test.Salary)
			assert.Equal(t, test.TaxCalculation, taxCalculation)
		})
	}
}

func TestCalculateByBracket(t *testing.T) {
	tests := map[string]struct {
		Bracket *taxclient.TaxBracket
		Salary  float32
		Tax     float32
	}{
		"first bracket": {
			Salary: 55000,
			Bracket: &taxclient.TaxBracket{
				Min:  0,
				Max:  50197,
				Rate: 0.15,
			},
			Tax: 7529.55,
		},
		"second bracket": {
			Salary: 55000,
			Bracket: &taxclient.TaxBracket{
				Min:  50197,
				Max:  100392,
				Rate: 0.205,
			},
			Tax: 984.62,
		},
		"out of bracket": {
			Salary: 50197,
			Bracket: &taxclient.TaxBracket{
				Min:  50197,
				Max:  100392,
				Rate: 0.205,
			},
			Tax: 0,
		},
		"bracket boundary": {
			Salary: 50197,
			Bracket: &taxclient.TaxBracket{
				Min:  0,
				Max:  50197,
				Rate: 0.15,
			},
			Tax: 7529.55,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tax := calculateByBracket(test.Bracket, test.Salary)
			assert.Equal(t, test.Tax, tax)
		})
	}
}

// type TaxCalculator interface {
// 	Calculate(ctx context.Context, year, salary string) (*TaxCalculation, taxclient.GetTaxBracketsResponse, error)
// }

// type taxCalculator struct {
// 	TaxClient taxclient.TaxClient
// }

// type mockTaxClient struct {
// 	mock.Mock
// }

// func (tc *mockTaxClient) GetBrackets(ctx context.Context, year string) ([]*taxclient.TaxBracket, taxclient.GetTaxBracketsResponse, error) {
// 	args := tc.Called(ctx, year)
// 	return args.Get(0).([]*taxclient.TaxBracket), args.Get(1).(taxclient.GetTaxBracketsResponse), args.Error(2)
// }

// func InitializeTaxCalculator(tc taxclient.TaxClient) TaxCalculator {
// 	return &taxCalculator{tc}
// }

// func (tc *taxCalculator) Calculate(ctx context.Context, year, salary string) (*TaxCalculation, taxclient.GetTaxBracketsResponse, error) {
// 	brackets, resp, err := tc.TaxClient.GetBrackets(ctx, year)
// 	if err != nil {
// 		return nil, resp, err
// 	}

// 	if resp == taxclient.Failed || resp == taxclient.NotFound {
// 		return nil, resp, nil
// 	}

// 	salaryF, err := strconv.ParseFloat(salary, 32)
// 	if err != nil {
// 		//TODO log err
// 		return nil, resp, err
// 	}

// 	var taxByBand []*TaxByBand
// 	var totalTax float32
// 	for _, bracket := range brackets {
// 		tax := calculateByBracket(bracket, float32(salaryF))
// 		taxByBand = append(taxByBand, &TaxByBand{
// 			format(tax),
// 			bracket,
// 		})
// 		totalTax += tax
// 	}

// 	answer := &TaxCalculation{
// 		format(totalTax),
// 		format(totalTax / float32(salaryF)),
// 		taxByBand,
// 	}
// 	return answer, resp, nil
// }
