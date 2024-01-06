/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
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
	for key, value := range values {
		if value != "" {
			log.Info(strings.ToUpper(key)+": ", value)

		} else {
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

// if mandatoryUnset {
// 	log.Error("KAEFFKEN EXITED")
// 	os.Exit(3)
// }
