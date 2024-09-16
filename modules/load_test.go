/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadClustersfile(t *testing.T) {

	LoadClustersfile(yamlExample)
}

func TestLoadDefaultKustomizations(t *testing.T) {
	assert := assert.New(t)
	expectedDefaultApps := "ingress-nginxlonghornmetallbcert-manager"
	var allDefaultApps string

	infraDefaults := LoadDefaultKustomizations(infraCatalog)

	for _, kustomization := range infraDefaults.Defaults {
		fmt.Println("CHECKING FOR", kustomization.Name)
		allDefaultApps = allDefaultApps + kustomization.Name
	}
	assert.Equal(allDefaultApps, expectedDefaultApps)
}

func TestLoadDataFromRepository(t *testing.T) {
	assert := assert.New(t)

	type test struct {
		filePath string
		want     bool
	}

	tests := []test{
		{filePath: "cmd/deploy.go", want: true},
		// {filePath: "cmd/deploy.go1", want: false},
	}

	repoworks := make(map[string]string)
	repoworks["repository"] = "https://github.com/stuttgart-things/kaeffken.git"
	repoworks["branch"] = "main"
	repoworks["gitCommitID"] = "2f2ada234ffee467195439dacfe4a5579be3f66a"

	repository, cloned := CloneGitRepository(repoworks)

	for _, tc := range tests {
		var fileLoaded bool

		if cloned {
			loadedFile := LoadDataFromRepository(repository, tc.filePath)

			if loadedFile != "" {
				fileLoaded = true
			}
		}

		assert.Equal(fileLoaded, tc.want)
	}
}

var (
	yamlExample = `clusters:
  - sthings-app1:
      cloud: vsphere
      ips:
        - 12.23.44.22
  - sthings-dev1:
      cloud: vsphere
      ips:
        - 2324.43.2.1`

	infraCatalog = `{
		"Defaults": [
		 {
		  "Name": "ingress-nginx",
		  "Namespace": "flux-system",
		  "Version": "1.2.3",
		  "Path": "./infra/ingress-nginx",
		  "SourceRef": {
		   "Name": "flux-system",
		   "Kind": "GitRepository"
		  },
		  "VersionTarget": {
		   "Name": "ingress-nginx",
		   "Kind": "HelmRelease",
		   "Namespace": "ingress-nginx",
		   "Path": "/spec/chart/spec/version",
		   "Version": "1.2.3"
		  },
		  "Substitute": {
		   "Variables": null,
		   "Secrets": null,
		   "SecretsResources": null
		  }
		 },
		 {
		  "Name": "longhorn",
		  "Namespace": "flux-system",
		  "Version": "4.8.3",
		  "Path": "./infra/longhorn",
		  "SourceRef": {
		   "Name": "flux-system",
		   "Kind": "GitRepository"
		  },
		  "VersionTarget": {
		   "Name": "longhorn",
		   "Kind": "HelmRelease",
		   "Namespace": "ingress-nginx",
		   "Path": "/spec/chart/spec/version",
		   "Version": "4.8.3"
		  },
		  "Substitute": {
		   "Variables": null,
		   "Secrets": null,
		   "SecretsResources": null
		  }
		 },
		 {
		  "Name": "metallb",
		  "Namespace": "flux-system",
		  "Version": "4.8.3",
		  "Path": "./infra/metallb",
		  "SourceRef": {
		   "Name": "flux-system",
		   "Kind": "GitRepository"
		  },
		  "VersionTarget": {
		   "Name": "metallb",
		   "Kind": "HelmRelease",
		   "Namespace": "ingress-nginx",
		   "Path": "/spec/chart/spec/version",
		   "Version": "4.8.3"
		  },
		  "Substitute": {
		   "Variables": {
			"IP_RANGE": ""
		   },
		   "Secrets": null,
		   "SecretsResources": null
		  }
		 },
		 {
		  "Name": "cert-manager",
		  "Namespace": "flux-system",
		  "Version": "5.1.3",
		  "Path": "./infra/cert-manager",
		  "SourceRef": {
		   "Name": "flux-system",
		   "Kind": "GitRepository"
		  },
		  "VersionTarget": {
		   "Name": "cert-manager",
		   "Kind": "HelmRelease",
		   "Namespace": "ingress-nginx",
		   "Path": "/spec/chart/spec/version",
		   "Version": "5.1.3"
		  },
		  "Substitute": {
		   "Variables": {
			"HOSTNAME": "defaultVar+cert-manager",
			"INGRESS_DOMAIN": "clusterVars+ingressDomain",
			"IP_RANGE": "invokeMethod+GetVIP"
		   },
		   "Secrets": {
			"APPROLE": ""
		   },
		   "SecretsResources": [
			"flux-vault-secrets"
		   ]
		  }
		 }
		]
	   }`
)
