/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	homerun "github.com/stuttgart-things/homerun-library"
)

func SendMessageToHomerun(homerunURL, authToken string, message homerun.Message) (err error, responseCode string) {
	// MARSHAL THE MESSAGE STRUCT TO JSON
	payloadBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return err, ""
	}

	// CREATE THE HTTP REQUEST
	req, err := http.NewRequest("POST", homerunURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return err, ""
	}

	// SET HEADERS
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", authToken)

	// CREATE AN HTTP CLIENT THAT IGNORES SSL CERTIFICATE VERIFICATION
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disable SSL verification
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: customTransport,
	}

	// MAKE THE REQUEST
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return err, ""
	}
	defer resp.Body.Close()

	// PRINT THE RESPONSE STATUS
	fmt.Printf("Response status: %s\n", resp.Status)

	return nil, resp.Status
}
