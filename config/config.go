package config

import (
	"flag"
	"fmt"
)

// Config is the configuration information needed to run the application
type Config struct {
	SpecPath   string
	SpecURL    string
	ListenAddr string
}

var c Config

// Setup reads inputs like flags and environment variables to initialize a config struct for export
func Setup() error {
	if c != (Config{}) {
		return fmt.Errorf("config.Setup() has already been called")
	}

	c = Config{}

	flag.StringVar(&c.SpecPath, "spec-path", "", "Path to an OpenAPI spec file")
	flag.StringVar(&c.SpecURL, "spec-url", "", "URL to an OpenAPI spec file")
	flag.StringVar(&c.ListenAddr, "listen-addr", ":3000", "Address to listen on")
	flag.Parse()

	return nil
}

// Get returns a pointer to the config struct
func Get() *Config {
	return &c
}
