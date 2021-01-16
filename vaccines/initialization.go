package vaccines

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
	"path"
)

const GitDataURL = "https://github.com/italia/covid19-opendata-vaccini.git"
const GitCommitUrl = "https://github.com/italia/covid19-opendata-vaccini/commit/%COMMIT%"
const SubPath = "git-vaccines"

type GitVaccines struct {
	repository *git.Repository
	lastHEAD *plumbing.Reference
	Data *Snapshot
}

func newRepository(gitdir string) (*git.Repository, error){
	r, err := git.PlainClone(gitdir, true, &git.CloneOptions{
		URL: GitDataURL,
	})
	if err != nil {
		return nil, fmt.Errorf("could not clone %s: %w", GitDataURL, err)
	}
	return r, nil
}

func New(basedir string) (*GitVaccines, error) {
	gitdir := path.Join(basedir, SubPath)
	var r *git.Repository
	var err error
	if _, exists := os.Stat(gitdir); os.IsNotExist(exists) {
		r, err = newRepository(gitdir)
	}else{
		r, err = openAndPull(gitdir)
	}
	if err != nil {
		return nil, err
	}
	h, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("could not get reference to HEAD: %w", err)
	}

	g := new(GitVaccines)
	g.repository = r
	g.lastHEAD = h
	return g, nil
}

func openAndPull(gitdir string) (*git.Repository, error) {
	r, err := git.PlainOpen(gitdir)
	if err != nil {
		return nil, fmt.Errorf("could not open repo in %s: %w", gitdir, err)
	}
	err = r.Fetch(&git.FetchOptions{})
	if err != nil && err.Error() != "already up-to-date" {
		return nil, fmt.Errorf("could not fetch: %w", err)
	}
	return r, nil
}

