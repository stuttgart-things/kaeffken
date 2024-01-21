/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	billy "github.com/go-git/go-billy/v5"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

type Kustomization struct {
	Name, Namespace, Version, Path string
	SourceRef                      SourceRef
	VersionTarget                  VersionTarget
	Substitute                     Substitute
}

type DefaultKustomizations struct {
	Defaults []Kustomization
}

type Substitute struct {
	Variables        map[string]string
	Secrets          map[string]string
	SecretsResources []string
}

type VersionTarget struct {
	Name, Kind, Namespace, Path, Version string
}

type SourceRef struct {
	Name, Kind string
}

// LOAD CLUSTERFILE - DEFAULT IS <ROOT>/<ENV>/<LAB>/clusters.yaml
func LoadDataFromRepository(repository billy.Filesystem, filePath string) (loadedFile string) {

	loadedFile = ""
	file := filepath.Base(filePath)
	folder := filepath.Dir(filePath)

	fmt.Println(file)
	fmt.Println(folder)

	// GET FILELIST
	fileList, directoryList := sthingsCli.GetFileListFromGitRepository(folder, repository)
	fmt.Println(fileList, directoryList)

	if sthingsBase.CheckForStringInSlice(fileList, file) {
		loadedFile = sthingsCli.ReadFileContentFromGitRepo(repository, filePath)
		return
	} else {
		log.Error("FILE DOES NOT EXIST IN REPOSITORY: ", filePath)
		return
	}
}

func LoadDefaultKustomizations(fileContent string) (defaults DefaultKustomizations) {

	err := json.Unmarshal([]byte(fileContent), &defaults)
	if err != nil {
		log.Fatal(err)
	}

	return defaults

}

func LoadClustersfile(yamlFileContent string) {

	fmt.Println(yamlFileContent)
	var allClusters Clusters

	allClusters = sthingsCli.ReadInlineYamlToObject([]byte(yamlFileContent), allClusters).(Clusters)
	fmt.Println(allClusters.ClusterProfile)

	for _, cluster := range allClusters.ClusterProfile {

		fmt.Println(cluster)
		for key, value := range cluster {
			fmt.Println(key, value)

		}

	}

}

type Clusters struct {
	ClusterProfile []map[string]Cluster `mapstructure:"clusters"`
}

type Cluster struct {
	Cloud string   `mapstructure:"cloud"`
	Ips   []string `mapstructure:"ips"`
}
