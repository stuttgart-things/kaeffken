/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	billy "github.com/go-git/go-billy/v5"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

// LOAD CLUSTERFILE - DEFAULT IS <ROOT>/<ENV>/<LAB>/clusters.yaml
func CloneGitRepository(values map[string]string) (repository billy.Filesystem, repositoryCloned bool) {
	repository, repositoryCloned = sthingsCli.CloneGitRepository(values["repository"], values["branch"], values["commitID"], nil)

	if !repositoryCloned {
		log.Error("GIT REPOSITORY CAN NOT BE CLONED: ", values["repository"])
		repositoryCloned = false
	}

	return repository, repositoryCloned
}
