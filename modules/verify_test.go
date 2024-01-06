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
	values1["branch"] = ""

	tests := []test{
		{mandatoryFlags: []string{"repository", "branch", "clusterName", "envPath"}, values: values1, want: false},
		{mandatoryFlags: []string{"repository"}, values: values1, want: true},
	}

	for _, tc := range tests {
		validValues := VerifyValues(tc.values, tc.mandatoryFlags)
		assert.Equal(validValues, tc.want)
	}

}
