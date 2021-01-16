package covid

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

const fileItaly = "dpc-covid19-ita-andamento-nazionale.json"
const fileRegions = "dpc-covid19-ita-regioni.json"
const fileProvinces = "dpc-covid19-ita-province.json"

const dateLayout = "2006-01-02T15:04:05"

func (d *GitCoviddi) HasNewData() (bool, error) {
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

func (d *GitCoviddi) Load(commit plumbing.Hash) error {
	commitObject, err := d.repository.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("could not load commit %s: %w", commit.String(), err)
	}
	tree, err := commitObject.Tree()
	if err != nil {
		return fmt.Errorf("could not load tree of commit %s: %w", commit.String(), err)
	}

	jsonTree, err := tree.Tree("dati-json")
	if err != nil {
		return fmt.Errorf("could not open dati-json: %w", err)
	}

	data := new(Snapshot)

	getJson := func(filename string) (*gabs.Container, error) {
		gitFile, err := jsonTree.File(filename)
		if err != nil {
			return nil, fmt.Errorf("could not read %s in commit %s: %w", filename, commit.String(), err)
		}
		reader, err := gitFile.Reader()
		if err != nil {
			return nil, fmt.Errorf("could not open %s in commit %s: %w", filename, commit.String(), err)
		}
		defer reader.Close()
		container, err := gabs.ParseJSONBuffer(reader)
		if err != nil{
			return nil, fmt.Errorf("could not parse json in file %s in commit %s: %w", filename, commit.String(), err)
		}
		return container, nil
	}

	data.Italy, err = getJson(fileItaly)
	if err != nil {
		return err
	}
	data.Regions, err = getJson(fileRegions)
	if err != nil {
		return err
	}
	
	data.Provinces, err = getJson(fileProvinces)
	if err != nil {
		return err
	}

	if err = stuzzica(data); err != nil {
		return err
	}

	data.Ref.Updated = commitObject.Author.When
	data.Ref.Permalink = strings.Replace(GitCommitUrl, "%COMMIT%", commit.String(), -1)
	data.Ref.Hash = commit.String()
	d.Data = data

	return nil
}

func stuzzica(data *Snapshot) error {
	var err error

	var p *gabs.Container
	p = nil
	err = stuzzicaFields(data.Italy, []string{
		"totale_ospedalizzati",
		"terapia_intensiva",
		"ricoverati_con_sintomi",
		"isolamento_domiciliare",
		"dimessi_guariti",
		"deceduti",
		"tamponi",
		"casi_testati",
	}, func(container *gabs.Container) *gabs.Container {
		prev := p
		p = container
		return prev
	})
	if err != nil {
		return fmt.Errorf("could not enrich italy data - %w", err)
	}

	regions := make(map[float64]*gabs.Container)
	err = stuzzicaFields(data.Regions, []string{
		"totale_ospedalizzati",
		"terapia_intensiva",
		"ricoverati_con_sintomi",
		"isolamento_domiciliare",
		"dimessi_guariti",
		"deceduti",
		"tamponi",
		"casi_testati",
	}, func(container *gabs.Container) *gabs.Container {
		regionCode, ok := container.S("codice_regione").Data().(float64)
		if !ok {
			regionCode = -1
		}
		prev, ok := regions[regionCode]
		if !ok || regionCode < 0{
			prev = nil
		}
		regions[regionCode] = container
		return prev
	})
	if err != nil {
		return fmt.Errorf("could not enrich region data - %w", err)
	}

	provinces := make(map[string]*gabs.Container)
	err = stuzzicaFields(data.Provinces, []string{
		"totale_casi",
	}, func(container *gabs.Container) *gabs.Container {
		province, ok := container.S("sigla_provincia").Data().(string)
		if !ok {
			province = "none"
		}
		prev, ok := provinces[province]
		if !ok || province == "none" {
			prev = nil
		}
		provinces[province] = container
		return prev
	})
	if err != nil {
		return fmt.Errorf("could not enrich province covid - %w", err)
	}

	return nil
}

func stuzzicaFields(container *gabs.Container, fields []string, previous func(*gabs.Container)*gabs.Container) error {
	for i, item := range container.Children() {
		p := previous(item)
		for _, field := range fields {
			if p == nil {
				_, err := item.Set(item.S(field), field + "_dt")
				if err != nil {
					return fmt.Errorf("could not set dt field %s at item %s: %w", field, i, err)
				}
				continue
			}
			previousValue, ok := p.S(field).Data().(float64)
			if !ok {
				if p.S(field).Data() == nil {
					previousValue = 0
				}else{
					return fmt.Errorf("field %s at previous item of %d is not a float64", field, i)
				}
			}
			currentValue, ok := item.S(field).Data().(float64)
			if !ok {
				if  item.S(field).Data() == nil {
					// field is not available yet
					item.Set(currentValue - previousValue, field + "_dt")
					continue
				}
				return fmt.Errorf("field %s at item %s is not a float64", field, i)
			}
			_, err := item.Set(currentValue - previousValue, field + "_dt")
			if err != nil {
				return fmt.Errorf("could not set dt field %s at item %s: %w", field, i, err)
			}
		}
	}
	return nil
}