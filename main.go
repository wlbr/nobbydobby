package main

import (
	"log"
)

var (
	Name           = "unknown"
	Version        = "0.0.0-dev"
	BuildTimestamp = "unknown"
)

func main() {
	log.Print("Name: ", Name, "    Version: ", Version, "    Timestamp:  ", BuildTimestamp)
	cfg := NewConfig(Name, Version, BuildTimestamp)
	defer cfg.CleanUp()

	//db, err := NewPostgresSink(cfg)
	//db, err := NewFlatFileDB(cfg.FlatFileName)
	db, err := NewBoltDatabase(cfg.BoltDBName)
	if err != nil {
		log.Printf("Could not get db connection: %v", err)
		cfg.FatalExit()
	}
	defer db.Close()

	RunRestserver(cfg, db)

}
