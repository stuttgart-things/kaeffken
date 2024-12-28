/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package modules

import (
	"gopkg.in/yaml.v3"
)

func CreateSecretsMap(yamlData []byte, wantedSecrets map[string]interface{}) (secrets map[string]interface{}) {

	secrets = make(map[string]interface{})

	// CREATE A VARIABLE TO HOLD THE PARSED DATA
	var data map[string]interface{}

	// PARSE THE YAML
	err := yaml.Unmarshal([]byte(yamlData), &data)
	if err != nil {
		log.Fatalf("ERROR PARSING YAML: %v", err)
	}

	// CREATE THE SECRETS MAP
	for _, value := range data {

		switch v := value.(type) {
		case map[string]interface{}:
			for k, v := range v {

				if wantedSecrets == nil {
					secrets[k] = v
				}

				// CHECK IF THE KEY IS IN THE WANTED SECRETS
				if _, ok := wantedSecrets[k]; ok {
					secrets[wantedSecrets[k].(string)] = v
				}

				// fmt.Println(k, v)
			}
		}

	}

	return
}
