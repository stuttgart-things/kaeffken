/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"fmt"
	"testing"
)

func TestCreateSecretsMap(t *testing.T) {

	yamlData := `
secrets:
  host: localhost
  port: 5432
  username: dbuser
  password: dbpassword
  name: example_db
`

	// NEEDED SECRETS
	wantedSecrets := map[string]interface{}{
		"host": "hostname",
		"port": "portNumber",
	}

	secrets := CreateSecretsMap([]byte(yamlData), wantedSecrets)
	fmt.Println(secrets)

}
