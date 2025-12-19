package main

import (
	"log"
	"os"
)

type Config struct {
	//ConfigFileName string
	cleanup        []func() error
	Name           string
	Version        string
	BuildTimestamp string
	BoltDBName     string
	FlatFileName   string

	PostgreSQL struct {
		Host     string
		Port     string
		Database string
		User     string
		Password string
	}
}

func NewConfig(name, version, buildTimestamp string) *Config {
	cfg := &Config{Name: name, Version: version, BuildTimestamp: buildTimestamp}
	cfg.PostgreSQL.Host = "localhost"
	cfg.PostgreSQL.Port = "5432"
	cfg.PostgreSQL.Database = name
	cfg.PostgreSQL.User = name + "app"
	cfg.PostgreSQL.Password = ""
	cfg.BoltDBName = name + ".db"
	cfg.FlatFileName = name + ".json"
	return cfg
}

func (cfg *Config) CleanUp() {
	log.Print("Cleaning up.")
	for _, fun := range cfg.cleanup {
		fun()
	}
}

func (cfg *Config) AddCleanUpFn(f func() error) {
	cfg.cleanup = append(cfg.cleanup, f)
}

func (cfg *Config) FatalExit() {
	cfg.CleanUp()
	os.Exit(1)
}
