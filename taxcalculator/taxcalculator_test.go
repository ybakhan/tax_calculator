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
