package models

// Define the CR struct
type CR struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Namespace  string `yaml:"namespace"`
}

// Define the SourceRef struct
type SourceRef struct {
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
}

// Define the Spec struct
type FluxSpec struct {
	Prune         bool      `yaml:"prune"`
	Wait          bool      `yaml:"wait"`
	Interval      string    `yaml:"interval"`
	RetryInterval string    `yaml:"retryInterval"`
	Timeout       string    `yaml:"timeout"`
	SourceRef     SourceRef `yaml:"sourceRef"`
}

// Define the FluxKustomization struct
type FluxKustomization struct {
	CR   CR       `yaml:"cr"`
	Spec FluxSpec `yaml:"spec"`
}

// Define the AppDefaults struct
type AppDefaults struct {
	FluxKustomization FluxKustomization `yaml:"fluxKustomization"`
}

// FLUX-DEFAULTS.YAML
// Define the Variables type
type Variables map[string]interface{}

// Define the FluxComponent struct
type FluxComponent struct {
	Repository string              `yaml:"repository"`
	Revision   string              `yaml:"revision"`
	Path       string              `yaml:"path"`
	Variables  map[string]Variable `yaml:"variables"`
	Spec       FluxSpec            `yaml:"spec"`
}

type FluxApp struct {
	Name      string              `yaml:"name"`
	Spec      Spec                `yaml:"spec"`
	Variables map[string]Variable `yaml:"variables"`
}

// Define the FluxDefaults struct with dynamic keys
type FluxDefaults struct {
	FluxAppDefaults map[string]FluxComponent `yaml:"fluxAppDefaults"`
}

// Define the Apps struct with dynamic keys
type Apps struct {
	Flux map[string]FluxApp `yaml:"flux"`
}

var templateKustomization = `
apiVersion: {{.APIVersion}}
kind: {{.Kind}}
metadata:
  name: {{.Metadata.Name}}
  namespace: {{.Metadata.Namespace}}
spec:
  interval: {{.Spec.Interval}}
  retryInterval: {{.Spec.RetryInterval}}
  timeout: {{.Spec.Timeout}}
  sourceRef:
    kind: {{.Spec.SourceRef.Kind}}
    name: {{.Spec.SourceRef.Name}}
  path: {{.Spec.Path}}
  prune: {{.Spec.Prune}}
  wait: {{.Spec.Wait}}
  postBuild:
    substitute:
{{- range $key, $value := .Spec.PostBuild.Substitute }}
      {{$key}}: {{$value}}
{{- end }}
`

// Define the Metadata struct
type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type Variable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Define the PostBuild struct
type PostBuild struct {
	Substitute map[string]interface{} `yaml:"substitute"`
}

// Define the Spec struct
type Spec struct {
	Interval      string    `yaml:"interval"`
	RetryInterval string    `yaml:"retryInterval"`
	Timeout       string    `yaml:"timeout"`
	SourceRef     SourceRef `yaml:"sourceRef"`
	Path          string    `yaml:"path"`
	Prune         bool      `yaml:"prune"`
	Wait          bool      `yaml:"wait"`
	PostBuild     PostBuild `yaml:"postBuild"`
}

// Define the Kustomization struct
type Kustomization struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}
