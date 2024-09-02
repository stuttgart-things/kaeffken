/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ReadYAMLFile[T any](filename string) (T, error) {
	var data T

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
