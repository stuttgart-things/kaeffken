/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
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
