/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"os"
	"reflect"
	"testing"
)

// Define the structs to match the YAML structures
type Variable struct {
	Name         string `yaml:"name"`
	DefaultValue string `yaml:"defaultValue"`
}

type AppDefaults struct {
	Repository string              `yaml:"repository"`
	Revision   string              `yaml:"revision"`
	Path       string              `yaml:"path"`
	Variables  map[string]Variable `yaml:"variables"`
}

type FluxDefaults struct {
	FluxAppDefaults map[string]AppDefaults `yaml:"fluxAppDefaults"`
}

func TestReadYAMLFile(t *testing.T) {
	// Create a temporary YAML file with sample data
	yamlContent := `
fluxAppDefaults:
  app1:
    repository: "repo1"
    revision: "rev1"
    path: "path1"
    variables:
      var1:
        name: "Variable 1"
        defaultValue: "Value 1"
      var2:
        name: "Variable 2"
        defaultValue: "Value 2"
  app2:
    repository: "repo2"
    revision: "rev2"
    path: "path2"
    variables:
      var3:
        name: "Variable 3"
        defaultValue: "Value 3"
      var4:
        name: "Variable 4"
        defaultValue: "Value 4"
`

	tmpFile, err := os.CreateTemp("", "test*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Use the readYAMLFile function to read and parse the temporary YAML file
	var fluxDefaults FluxDefaults
	fluxDefaults, err = ReadYAMLFileFromDisk[FluxDefaults](tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read YAML file: %v", err)
	}

	// Define the expected data
	expected := FluxDefaults{
		FluxAppDefaults: map[string]AppDefaults{
			"app1": {
				Repository: "repo1",
				Revision:   "rev1",
				Path:       "path1",
				Variables: map[string]Variable{
					"var1": {Name: "Variable 1", DefaultValue: "Value 1"},
					"var2": {Name: "Variable 2", DefaultValue: "Value 2"},
				},
			},
			"app2": {
				Repository: "repo2",
				Revision:   "rev2",
				Path:       "path2",
				Variables: map[string]Variable{
					"var3": {Name: "Variable 3", DefaultValue: "Value 3"},
					"var4": {Name: "Variable 4", DefaultValue: "Value 4"},
				},
			},
		},
	}

	// Verify that the parsed data matches the expected data
	if !reflect.DeepEqual(fluxDefaults, expected) {
		t.Errorf("Parsed data does not match expected data.\nGot: %+v\nExpected: %+v", fluxDefaults, expected)
	}
}
