package models

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
