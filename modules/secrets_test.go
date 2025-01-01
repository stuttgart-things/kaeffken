/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"fmt"
	"testing"
)

func TestCreateSecretsMap(t *testing.T) {

	flatYAML := `
host: localhost
port: 5432
username: dbuser
password: dbpassword
name: example_db`

	dictYAML := `
secrets:
  host: localhost
  port: 5432
  username: dbuser
  password: dbpassword
  name: example_db`

	// NEEDED SECRETS
	wantedSecrets := map[string]interface{}{
		"host": "hostname",
		"port": "portNumber",
	}

	secretsFlat := CreateSecretsMap([]byte(flatYAML), wantedSecrets)
	fmt.Println("RESULT flatYAML", secretsFlat)

	secretsDict := CreateSecretsMap([]byte(dictYAML), wantedSecrets)
	fmt.Println("RESULT dictYAML", secretsDict)

}
