/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"fmt"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	logger = sthingsBase.StdOutFileLogger("/tmp/machineShop.log", "2006-01-02 15:04:05", 50, 3, 28)
)

func HandleOutput(outputFormat, destinationPath, renderedTemplate string) {

	switch outputFormat {
	default:
		logger.Error(outputFormat, "output format not defined")

	case "stdout":
		fmt.Println(string(renderedTemplate))

	case "file":

		if destinationPath == "" {
			logger.Warn("destinationPath empty")
			destinationPath = "/tmp/encrypted.yaml"
		}

		logger.Info("output file written to ", destinationPath)
		sthingsBase.WriteDataToFile(destinationPath, string(renderedTemplate))
	}

}

func HandleRenderOutput(renderedTemplates map[string]string, outputFormat, destinationPath string) {

	for appKey, renderedTemplate := range renderedTemplates {

		if outputFormat == "stdout" {
			fmt.Println(renderedTemplate)
		} else {
			// sthingsBase.CreateNestedDirectoryStructure(destinationPath, 0600)
			log.Info("output file written to ", destinationPath+"/"+appKey+".yaml")
			sthingsBase.WriteDataToFile(destinationPath+"/"+appKey+".yaml", renderedTemplate)
		}

	}

}
