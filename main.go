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
	cfg := NewConfig(Name, Version, BuildTimestamp)
	defer cfg.CleanUp()

	db, err := NewPostgresSink(cfg)
	if err != nil {
		log.Printf("Could not get db connection: %v", err)
		cfg.FatalExit()
	}

	RunRestserver(cfg, db)

}
