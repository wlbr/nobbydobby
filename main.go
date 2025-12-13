package main

import (
	"log"
)

func main() {
	cfg := &Config{}
	defer cfg.CleanUp()
	cfg.PostgreSQL.Host = "localhost"
	cfg.PostgreSQL.Port = "5432"
	cfg.PostgreSQL.Database = "felix"
	cfg.PostgreSQL.User = "felixapp"
	cfg.PostgreSQL.Password = ""

	db, err := NewPostgresSink(cfg)
	if err != nil {
		log.Printf("Could not get db connection: %v", err)
		cfg.FatalExit()
	}

	RunRestserver(cfg, db)

}
