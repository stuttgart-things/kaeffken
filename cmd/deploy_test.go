/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	expectedFileList      = []string{".gitignore", "LICENSE", "README.md"}
	expectedDirectoryList = []string{}
)

func TestCloneGitRepository(t *testing.T) {

	gitRepository = "https://github.com/stuttgart-things/kaeffken.git" //pragma: allowlist secre
	gitBranch = "main"
	gitCommitID = "09de9ff7b5c76aff8bb32f68cfb0bbe49cd5a7a8" //pragma: allowlist secret

	assert := assert.New(t)
	_, cloned := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
	assert.Equal(cloned, true)
	fmt.Println("TEST SUCCESSFULLY")
}

func TestGetFileListFromGitRepository(t *testing.T) {

	var fileList []string
	var directoryList []string

	gitRepository = "https://github.com/stuttgart-things/kaeffken.git"
	gitBranch = "main"
	gitCommitID = "09de9ff7b5c76aff8bb32f68cfb0bbe49cd5a7a8" // pragma: allowlist secret

	repo, cloned := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)

	if cloned {
		fileList, directoryList = sthingsCli.GetFileListFromGitRepository("", repo)
		fmt.Println(fileList, directoryList)
	}

	if !reflect.DeepEqual(fileList, expectedFileList) && reflect.DeepEqual(directoryList, expectedDirectoryList) {
		t.Errorf("EXPECTED LISTS DIFFER")
	} else {
		fmt.Println("TEST SUCCESSFULLY")
	}

}

func TestReadFileContentFromGitRepo(t *testing.T) {

	gitRepository = "https://github.com/stuttgart-things/kaeffken.git"
	gitBranch = "main"
	gitCommitID = "09de9ff7b5c76aff8bb32f68cfb0bbe49cd5a7a8" //pragma: allowlist secret

	assert := assert.New(t)
	expectedReadMe := "# kaeffken\ngitops cluster management cli \n"

	repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
	readMe := sthingsCli.ReadFileContentFromGitRepo(repo, "README.md")
	fmt.Println(readMe)
	fmt.Println(expectedReadMe)

	assert.Equal(readMe, expectedReadMe)
	fmt.Println("TEST SUCCESSFULLY")

}
