package datastore

import (
	"covidiometro/util"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"time"
)
import "gopkg.in/src-d/go-git.v4/storage/memory"
import "encoding/json"

type Dati struct {
	Italia     []map[string]interface{}
	ItaliaRaw  string
	lastUpdate time.Time
	ttl        time.Time
}

const GitRepository = "https://github.com/pcm-dpc/COVID-19.git"
const DatiItalia = "dpc-covid19-ita-andamento-nazionale.json"

var update, _ = time.ParseDuration("1h")
var dati = Dati{
	Italia:     nil,
	ItaliaRaw:  "",
	lastUpdate: time.Unix(0, 0),
}

func clone() (*object.Tree, error) {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: GitRepository,
		Depth:1,
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

func Updater(consumer chan Dati) {
	for ; true;  {
		if dati.ttl.Before(time.Now()){
			log.Println("Updating repo - data is old")
			updateAll()
			log.Println("Update done")
		}
		consumer <- dati
	}
}

func updateAll(){
	tree, err := clone()
	util.ErrFatal(err)
	italiaFile, err := tree.File(DatiItalia)
	util.ErrFatal(err)
	italia, err := italiaFile.Contents()
	util.ErrFatal(err)
	util.ErrFatal(json.Unmarshal([]byte(italia), &dati.Italia))
	dati.ItaliaRaw = italia
	now := time.Now()
	dati.lastUpdate = now
	dati.ttl = now.Add(update)
}