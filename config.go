package main

import (
	"log"
	"os"
)

type Config struct {
	ConfigFileName string
	cleanup        []func() error
	Name           string
	Version        string
	BuildTimestamp string

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
	cfg.PostgreSQL.Database = Name
	cfg.PostgreSQL.User = Name + "app"
	cfg.PostgreSQL.Password = ""
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
