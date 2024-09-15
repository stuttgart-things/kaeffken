package models

// func test() {

// 	kustomization := Kustomization{
// 		APIVersion: "kustomize.toolkit.fluxcd.io/v1",
// 		Kind:       "Kustomization",
// 		Metadata: struct {
// 			Name      string `yaml:"name"`
// 			Namespace string `yaml:"namespace"`
// 		}{
// 			Name:      "ingress-nginx",
// 			Namespace: "flux-system",
// 		},
// 		Spec: struct {
// 			Interval      string `yaml:"interval"`
// 			RetryInterval string `yaml:"retryInterval"`
// 			Timeout       string `yaml:"timeout"`
// 			SourceRef     struct {
// 				Kind string `yaml:"kind"`
// 				Name string `yaml:"name"`
// 			} `yaml:"sourceRef"`
// 			Path      string `yaml:"path"`
// 			Prune     bool   `yaml:"prune"`
// 			Wait      bool   `yaml:"wait"`
// 			PostBuild struct {
// 				Substitute map[string]string `yaml:"substitute"`
// 			} `yaml:"postBuild"`
// 		}{
// 			Interval:      "1h",
// 			RetryInterval: "1m",
// 			Timeout:       "5m",
// 			SourceRef: struct {
// 				Kind string `yaml:"kind"`
// 				Name string `yaml:"name"`
// 			}{
// 				Kind: "GitRepository",
// 				Name: "flux-system",
// 			},
// 			Path:  "./infra/ingress-nginx",
// 			Prune: true,
// 			Wait:  true,
// 			PostBuild: struct {
// 				Substitute map[string]string `yaml:"substitute"`
// 			}{
// 				Substitute: map[string]string{
// 					"INGRESS_NGINX_NAMESPACE":     "ingress-nginx",
// 					"INGRESS_NGINX_CHART_VERSION": "4.11.2",
// 				},
// 			},
// 		},
// 	}

// }
