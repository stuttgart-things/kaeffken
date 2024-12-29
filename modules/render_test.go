/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"fmt"
	"testing"
)

// TestRenderTemplate tests the RenderTemplate function
func TestRenderTemplate(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		tmpl     string
		data     interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple template with map",
			tmpl:     "Hello, {{.Name}}!",
			data:     map[string]string{"Name": "World"},
			expected: "Hello, World!",
			wantErr:  false,
		},
		{
			name: "Template with struct",
			tmpl: "Name: {{.Name}}, Age: {{.Age}}",
			data: struct {
				Name string
				Age  int
			}{Name: "Alice", Age: 30},
			expected: "Name: Alice, Age: 30",
			wantErr:  false,
		},
		{
			name:     "Template with missing field",
			tmpl:     "Hello, {{.Name}}!",
			data:     struct{ Age int }{Age: 30},
			expected: "",
			wantErr:  true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderTemplate(tt.tmpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("RenderTemplate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestRenderAliases tests the RenderAliases function
func TestRenderAliases(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		aliases   []string
		allValues map[string]interface{}
		expected  map[string]interface{}
	}{
		{
			name:      "Simple alias",
			aliases:   []string{"{{.Name}}:{{.Age}}"},
			allValues: map[string]interface{}{"Name": "Alice", "Age": 30},
			expected:  map[string]interface{}{"Age": 30, "Alice": 30, "Name": "Alice"},
		},
		{
			name:      "Multiple aliases",
			aliases:   []string{"{{.Name}}:{{.Age}}", "{{.City}}:{{.Country}}"},
			allValues: map[string]interface{}{"Name": "Alice", "Age": 30, "City": "Berlin", "Country": "Germany"},
			expected:  map[string]interface{}{"Age": 30, "Alice": 30, "Berlin": "Germany", "City": "Berlin", "Country": "Germany", "Name": "Alice"},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderAliases(tt.aliases, tt.allValues)
			fmt.Println("RESULT: ", result)

		})
	}
}

// // Helper function to compare two maps
// func mapsEqual(a, b map[string]interface{}) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for key, valueA := range a {
// 		if valueB, exists := b[key]; !exists || valueA != valueB {
// 			return false
// 		}
// 	}
// 	return true
// }
