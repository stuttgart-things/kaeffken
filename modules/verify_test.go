/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyValues(t *testing.T) {

	assert := assert.New(t)

	type test struct {
		mandatoryFlags []string
		values         map[string]string
		want           bool
	}

	values1 := make(map[string]string)
	values1["repository"] = "https://github.com/stuttgart-things/stuttgart-things.git"
	values1["branch"] = "main"

	tests := []test{
		{mandatoryFlags: []string{"repository", "branch", "clusterName", "envPath"}, values: values1, want: true},
		{mandatoryFlags: []string{"repository"}, values: values1, want: true},
	}

	for _, tc := range tests {
		validValues := VerifyValues(tc.values, tc.mandatoryFlags)
		assert.Equal(validValues, tc.want)
	}

}

func TestSetAppParameter(t *testing.T) {
	// Define test cases
	tests := []struct {
		name              string
		appValue          string
		appDefault        string
		technologyDefault string
		expected          string
	}{
		{
			name:              "App value is set",
			appValue:          "appValue",
			appDefault:        "appDefault",
			technologyDefault: "technologyDefault",
			expected:          "appValue",
		},
		{
			name:              "App value is empty, app default is set",
			appValue:          "",
			appDefault:        "appDefault",
			technologyDefault: "technologyDefault",
			expected:          "appDefault",
		},
		{
			name:              "App value and app default are empty, technology default is set",
			appValue:          "",
			appDefault:        "",
			technologyDefault: "technologyDefault",
			expected:          "technologyDefault",
		},
		{
			name:              "All values are empty",
			appValue:          "",
			appDefault:        "",
			technologyDefault: "",
			expected:          "",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SetAppParameter(tt.appValue, tt.appDefault, tt.technologyDefault)
			if result != tt.expected {
				t.Errorf("SetAppParameter(%q, %q, %q) = %q; want %q", tt.appValue, tt.appDefault, tt.technologyDefault, result, tt.expected)
			}
		})
	}
}
