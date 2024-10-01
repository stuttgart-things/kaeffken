/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package models

type Secret struct {
	Name string   `yaml:"name"`
	Data []string `yaml:"secretData"`
}

var K8sSecret = `---
apiVersion: v1
kind: Secret
metadata:
  name: {{ .metaName }}
  namespace: {{ .metaNamespace }}
type: Opaque
stringData:
{{- range $key, $value := .Data }}
  {{ $key }}: {{ $value }}
{{- end }}
`
