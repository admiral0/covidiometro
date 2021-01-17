package vaccines

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"strings"
)

func (d *GitVaccines) HasNewData() (bool, error) {
	if d.Data == nil {
		return true, nil
	}
	err := d.repository.Fetch(&git.FetchOptions{})
	if err != nil && err.Error() != "already up-to-date" {
		return false, fmt.Errorf("could not fetch during update check: %w", err)
	}
	h, err := d.repository.Head()
	if err != nil {
		return false, fmt.Errorf("could not get HEAD during update check: %w", err)
	}

	return d.lastHEAD.Hash() != h.Hash(), nil
}

func (d *GitVaccines) Load(hash plumbing.Hash) error {
	commitObject, err := d.repository.CommitObject(hash)
	if err != nil {
		return fmt.Errorf("could not load commit %s: %w", hash.String(), err)
	}
	tree, err := commitObject.Tree()
	if err != nil {
		return fmt.Errorf("could not load tree of commit %s: %w", hash.String(), err)
	}

	jsonTree, err := tree.Tree("dati")
	if err != nil {
		return fmt.Errorf("could not open dati: %w", err)
	}

	data := new(Snapshot)

	filename := "somministrazioni-vaccini-latest.json"
	gitFile, err := jsonTree.File(filename)
	if err != nil {
		return fmt.Errorf("could not read %s in commit %s: %w", filename, hash.String(), err)
	}
	reader, err := gitFile.Reader()
	if err != nil {
		return fmt.Errorf("could not open %s in commit %s: %w", filename, hash.String(), err)
	}
	defer reader.Close()
	data.Regions, err = gabs.ParseJSONBuffer(reader)
	if err != nil{
		return fmt.Errorf("could not parse json in file %s in commit %s: %w", filename, hash.String(), err)
	}

	data.Ref.Updated = commitObject.Author.When
	data.Ref.Permalink = strings.Replace(GitCommitUrl, "%COMMIT%", hash.String(), -1)
	data.Ref.Hash = hash.String()

	d.Data = data

	return nil
}