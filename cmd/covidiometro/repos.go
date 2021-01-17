package main

import (
	"fmt"
	"github.com/admiral0/covidiometro"
	"github.com/admiral0/covidiometro/covid"
	"github.com/admiral0/covidiometro/vaccines"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"time"
)

type DataSources struct {
	covid *covid.GitCoviddi
	vaccines *vaccines.GitVaccines

	updateStoppers []func()
}

func (r *DataSources) StopUpdates() {
	for _, stopPlease := range r.updateStoppers {
		stopPlease()
	}
}



type Repo interface {
	Head() (plumbing.Hash, error)
	Load(plumbing.Hash) error
	HasNewData() (bool, error)
	WalkUp(plumbing.Hash, func(covidiometro.RefInfo,error))error
}

func InitializeRepositories(basedir string) (*DataSources, error) {
	ds := new(DataSources)
	var err error
	ds.covid, err = covid.New(basedir)
	if err != nil {
		return nil, fmt.Errorf("could not initialize covid repo: %w", err)
	}
	ds.vaccines, err = vaccines.New(basedir)
	if err != nil {
		return nil, fmt.Errorf("could not initialize vaccines repo: %w", err)
	}

	err = updateRepo(ds.covid, func(info covidiometro.RefInfo, err error) {
		log.Warn().Str("repo","covid").Str("commit", info.Hash).Err(err).
			Msg("initial commit load error")
	})
	if err != nil {
		return nil, fmt.Errorf("could not perform initial load for covid: %w", err)
	}
	err = updateRepo(ds.vaccines, func(info covidiometro.RefInfo, err error) {
		log.Warn().Str("repo","vaccines").Str("commit", info.Hash).Err(err).
			Msg("initial commit load error")
	})
	if err != nil {
		return nil, fmt.Errorf("could not perform initial load for vaccines: %w", err)
	}

	ds.updateStoppers = append(ds.updateStoppers, setUpAutoUpdates(ds.covid, func(info covidiometro.RefInfo, err error) {
		log.Warn().Str("repo","covid").Str("commit", info.Hash).Err(err).
			Msg("commit load error")
	}))
	ds.updateStoppers = append(ds.updateStoppers, setUpAutoUpdates(ds.vaccines, func(info covidiometro.RefInfo, err error) {
		log.Warn().Str("repo","vaccines").Str("commit", info.Hash).Err(err).
			Msg("commit load error")
	}))

	return ds, nil
}

func updateRepo(repo Repo, errorHandler func(covidiometro.RefInfo,error)) error {
	h, err := repo.Head()
	if err != nil {
		return fmt.Errorf("could not get reference to HEAD: %w", err)
	}
	needsUpdate, err := repo.HasNewData()
	if err != nil {
		return fmt.Errorf("could not check for update: %w", err)
	}
	if needsUpdate {
		err = repo.WalkUp(h, errorHandler)
		if err != nil {
			return err
		}
	}
	return nil
}

func setUpAutoUpdates(repo Repo, errorHandler func(covidiometro.RefInfo,error)) func() {
	stopper := make(chan bool)
	ticker := time.NewTicker(5* time.Minute)

	go func() {
		for {
			select {
			case <-stopper:
				return
			case <-ticker.C:
				err := updateRepo(repo, errorHandler)
				if err != nil {
					panic(fmt.Errorf("fatal error during update: %w", err))
				}
			}
		}
	}()

	return func() {
		ticker.Stop()
		stopper <- true
	}
}