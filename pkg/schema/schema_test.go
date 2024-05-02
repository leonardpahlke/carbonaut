package schema_test

import (
	"testing"

	"carbonaut.cloud/pkg/schema"
)

func TestConvertToKilowatt(t *testing.T) {
	testCases := []struct {
		name        string
		amount      float64
		unit        schema.Unit
		expected    float64
		description string
	}{
		{
			name:        "Test Microwatt to Kilowatt Conversion",
			amount:      1000,
			unit:        schema.MICROWATT,
			expected:    0.000001,
			description: "Testing conversion from microwatt to kilowatt",
		},
		{
			name:        "Test Milliwatt to Kilowatt Conversion",
			amount:      1000,
			unit:        schema.MILLIWATT,
			expected:    1e-3,
			description: "Testing conversion from milliwatt to kilowatt",
		},
		{
			name:        "Test Kilowatt to Kilowatt Conversion",
			amount:      1,
			unit:        schema.KILOWATT,
			expected:    1,
			description: "Testing conversion from kilowatt to kilowatt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			energy := schema.Energy{
				Amount: tc.amount,
				Unit:   tc.unit,
			}
			result := energy.ConvertToKilowatt()
			if result != tc.expected {
				t.Errorf("Expected: %f, Got: %f. %s", tc.expected, result, tc.description)
			}
		})
	}
}

func TestConvertToMilliwatt(t *testing.T) {
	testCases := []struct {
		name        string
		amount      float64
		unit        schema.Unit
		expected    float64
		description string
	}{
		{
			name:        "Test Milliwatt to Milliwatt Conversion",
			amount:      1000,
			unit:        schema.MILLIWATT,
			expected:    1000,
			description: "Testing conversion from milliwatt to milliwatt",
		},
		{
			name:        "Test Kilowatt to Milliwatt Conversion",
			amount:      1,
			unit:        schema.KILOWATT,
			expected:    1000.000000,
			description: "Testing conversion from kilowatt to milliwatt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			energy := schema.Energy{
				Amount: tc.amount,
				Unit:   tc.unit,
			}
			result := energy.ConvertToMilliwatt()
			if result != tc.expected {
				t.Errorf("Expected: %f, Got: %f. %s", tc.expected, result, tc.description)
			}
		})
	}
}
