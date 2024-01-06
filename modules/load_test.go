/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDataFromRepository(t *testing.T) {

	assert := assert.New(t)

	type test struct {
		filePath string
		want     bool
	}

	tests := []test{
		{filePath: "cmd/deploy.go", want: true},
		{filePath: "cmd/deploy.go1", want: false},
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
