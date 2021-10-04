package repository

import "errors"

// repoToOrg associates repo to org.
var repoToOrg = map[string]string{
	"rke2": "rancher",
	"k3s":  "k3s-io",
}

// OrgFromRepo
func OrgFromRepo(repo string) (string, error) {
	if repo, ok := repoToOrg[repo]; ok {
		return repo, nil
	}

	return "", errors.New("repo not found: " + repo)
}

// IsValidRepo determines if the given
// repository is valid for this program
// to operate against.
func IsValidRepo(repo string) bool {
	for r := range repoToOrg {
		if repo == r {
			return true
		}
	}

	return false
}
