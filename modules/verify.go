/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"strings"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	log = sthingsBase.StdOutFileLogger("/tmp/machineShop.log", "2006-01-02 15:04:05", 50, 3, 28)
)

// VERIFY / OUTPUT ALL VALUES
func VerifyValues(values map[string]string, mandatoryFlags []string) (validValues bool) {

	validValues = true

	// RUN OVER ALL VALUES
	for key, value := range values {

		// OUTPUT IF VALUE IS SET
		if value != "" {
			log.Info(strings.ToUpper(key)+": ", value)

		} else {
			// CHECK IF UNSET VALUE IS MANDATORY OR NOT
			if sthingsBase.CheckForStringInSlice(mandatoryFlags, key) {
				log.Error(strings.ToUpper(key) + " is unset")
				validValues = false
			} else {
				log.Warn(strings.ToUpper(key) + " is unset")
			}
		}
	}

	return validValues
}

func SetAppParameter(appValue, appDefault, technologyDefault string) string {
	if appValue != "" {
		return appValue
	} else if appDefault != "" {
		return appDefault
	} else {
		return technologyDefault
	}
}
