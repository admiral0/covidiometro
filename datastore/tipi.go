package datastore

import (
	"sync"
	"time"
)

const GitRepository = "https://github.com/pcm-dpc/COVID-19.git"

const DatiItalia = "dpc-covid19-ita-andamento-nazionale.json"
const DatiRegioni = "dpc-covid19-ita-regioni.json"
const DatiProvince = "dpc-covid19-ita-province.json"

type Dati struct {
	AndroidBundle []byte

	Italia   []map[string]interface{}
	Regioni  []map[string]interface{}
	Province []map[string]interface{}

	MappaRegioni map[int]string

	lastUpdate time.Time
	ttl        time.Time
}

type DataHolder struct {
	dati *Dati
	mux  sync.Mutex
}

func Inizializzazione() {
	dati := NuoviDati()
	Holder.dati = dati
}

func NuoviDati() *Dati {
	dati := Dati{
		AndroidBundle: nil,
		Italia:        nil,
		Regioni:       nil,
		Province:      nil,
		MappaRegioni: make(map[int]string),
		lastUpdate:    time.Unix(0, 0),
		ttl:           time.Now(),
	}
	return &dati
}

func (h *DataHolder) Get() *Dati {
	h.mux.Lock()
	p := h.dati
	h.mux.Unlock()
	return p
}

func (h *DataHolder) Put(d *Dati) {
	h.mux.Lock()
	h.dati = d
	h.mux.Unlock()
}
