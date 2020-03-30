package datastore

import (
	"covidiometro/util"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func Clone() (*object.Tree, error) {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:   GitRepository,
		Depth: 1,
	})
	util.ErrFatal(err)

	ref, err := r.Head()
	util.ErrFatal(err)
	commit, err := r.CommitObject(ref.Hash())
	util.ErrFatal(err)

	tree, err := commit.Tree()
	util.ErrFatal(err)
	dati, err := tree.Tree("dati-json")
	util.ErrFatal(err)
	return dati, nil
}
