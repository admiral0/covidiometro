package vaccines

import (
	"fmt"
	"github.com/admiral0/covidiometro"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"strings"
)

func GetRefInfo(commit *object.Commit) covidiometro.RefInfo {
	return covidiometro.RefInfo{
		Updated:   commit.Author.When,
		Hash:      commit.Hash.String(),
		Permalink: strings.Replace(GitCommitUrl, "%COMMIT%", commit.String(), -1),
	}
}

func GetMinimalRefInfo(hash plumbing.Hash) covidiometro.RefInfo {
	return covidiometro.RefInfo{
		Hash:      hash.String(),
		Permalink: strings.Replace(GitCommitUrl, "%COMMIT%", hash.String(), -1),
	}
}

func (d *GitVaccines) WalkUp(hash plumbing.Hash, errorHandler func(covidiometro.RefInfo, error)) error {
	toCheck := []plumbing.Hash{hash}
	for len(toCheck)>0 {
		newChecks := make([]plumbing.Hash, 0)
		for _, h := range toCheck {
			commit, err := d.repository.CommitObject(h)
			if err != nil {
				errorHandler(GetMinimalRefInfo(hash), err)
			}
			err = d.Load(commit.Hash)
			if err == nil {
				return nil
			}
			errorHandler(GetRefInfo(commit), err)
			for _, parent := range commit.ParentHashes {
				newChecks = append(newChecks, parent)
			}
		}
		toCheck = newChecks
	}
	return fmt.Errorf("no loadable commits")
}

func (d *GitVaccines) Head() (plumbing.Hash, error) {
	ref, err := d.repository.Head()
	if err != nil {
		return [20]byte{}, fmt.Errorf("could not open HEAD: %w", err)
	}
	return ref.Hash(), nil
}