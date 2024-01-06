/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"fmt"
	"path/filepath"

	billy "github.com/go-git/go-billy/v5"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

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
