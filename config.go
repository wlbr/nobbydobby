package main

import (
	"log"
	"os"
)

type Config struct {
	ConfigFileName string
	cleanup        []func() error

	PostgreSQL struct {
		Host     string
		Port     string
		Database string
		User     string
		Password string
	}
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
