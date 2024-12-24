/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"testing"
	"time"

	homerun "github.com/stuttgart-things/homerun-library"
)

func TestSendMessageToHomerun(t *testing.T) {

	// Define the address and token
	address := "https://homerun.homerun-dev.sthings-vsphere.labul.sva.de/generic"
	authToken := "IhrGeheimerToken"

	// Create a Message instance
	message := homerun.Message{
		Title:           "Memory System Alert",
		Message:         "Memory usage is high",
		Severity:        "warning",
		Author:          "monitoring-system",
		Timestamp:       time.Now().UTC().Format(time.RFC3339), // Generate current timestamp
		System:          "flux",
		Tags:            "cpu,usage,alert",
		AssigneeAddress: "admin@example.com",
		AssigneeName:    "Admin",
		Artifacts:       "Admin",
		Url:             "Admin",
	}

	// Send the message
	err, respCode := SendMessageToHomerun(address, authToken, message)
	if err != nil {
		t.Errorf("Error sending message: %v", err)
	}

	if respCode != "200 OK" {
		t.Errorf("Unexpected response code: %s", respCode)
	}

}
