package datastore

import (
	"log"
	"time"
)

var update, _ = time.ParseDuration("1h")
var Holder DataHolder

func Updater() {
	for true {
		d := Holder.Get()
		if d.ttl.Before(time.Now()) {
			log.Println("Updating repo - data is old")
			Holder.Put(ScaricaDati())
			log.Println("Update done")
		}
	}
}
