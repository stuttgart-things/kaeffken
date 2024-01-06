/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneGitRepository(t *testing.T) {

	assert := assert.New(t)

	repoworks := make(map[string]string)
	repoworks["repository"] = "https://github.com/stuttgart-things/kaeffken.git"
	repoworks["branch"] = "main"
	repoworks["gitCommitID"] = "09de9ff7b5c76aff8bb32f68cfb0bbe49cd5a7a8"

	repoworksNot := make(map[string]string)
	repoworksNot["repository"] = "https://github.com/stuttgart-things/kaeffken.git"
	repoworksNot["branch"] = "test"
	repoworksNot["gitCommitID"] = "09de9ff7b5c76aff8bb32f68cfb0bbe49cd5a7a8"

	_, cloned := CloneGitRepository(repoworks)
	assert.Equal(cloned, true)

	_, cloned = CloneGitRepository(repoworksNot)
	assert.Equal(cloned, false)

}
