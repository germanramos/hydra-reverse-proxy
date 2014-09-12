package config

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/innotech/hydra-reverse-proxy/vendors/github.com/BurntSushi/toml"
)

type Config struct {
	AppId          string `toml:"app_id"`
	ConfigFilePath string
	HydraServers   []string `toml:"hydra_servers"`
	ProxyAddr      string   `toml:"proxy_addr"`
	HydraClient    struct {
		AppsCacheDuration              uint `toml:"apps_cache_duration"`
		DurationBetweenAllServersRetry uint `toml:"duration_between_all_servers_retry"`
		HydraServersCacheDuration      uint `toml:"hydra_servers_cache_duration"`
		MaxNumberOfRetries             uint `toml:"max_number_of_retries"`
	}
}

func New() *Config {
	c := new(Config)

	return c
}

// Load configures hydra-reverse-proxy, it can be loaded from both system file,
// custom file or command line arguments and the values extracted from
// files they can be overriden with the command line arguments.
func (c *Config) Load(arguments []string) error {
	var path string
	f := flag.NewFlagSet("hydra-reverse-proxy", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&path, "config", "", "path to config file")
	f.Parse(arguments)

	if path != "" {
		// Load from config file specified in arguments.
		if err := c.LoadFile(path); err != nil {
			return err
		}
	} else {
		// Load from system file.
		if err := c.LoadSystemFile(); err != nil {
			return err
		}

	}

	// Load from command line flags.
	if err := c.LoadFlags(arguments); err != nil {
		return err
	}

	return nil
}

// LoadSystemFile loads from the default hydra configuration file path if it exists.
func (c *Config) LoadSystemFile() error {
	if _, err := os.Stat(c.ConfigFilePath); os.IsNotExist(err) {
		return nil
	}
	return c.LoadFile(c.ConfigFilePath)
}

// LoadFile loads configuration from a file.
func (c *Config) LoadFile(path string) error {
	_, err := toml.DecodeFile(path, &c)
	return err
}

// LoadFlags loads configuration from command line flags.
func (c *Config) LoadFlags(arguments []string) error {
	var hydraServers, ignoredString string

	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&c.AppId, "app-id", c.AppId, "")
	f.StringVar(&hydraServers, "hydra-servers", "", "")
	f.StringVar(&c.ProxyAddr, "proxy-addr", c.ProxyAddr, "")
	f.UintVar(&c.HydraClient.AppsCacheDuration, "apps-cache-duration", c.HydraClient.AppsCacheDuration, "")
	f.UintVar(&c.HydraClient.DurationBetweenAllServersRetry, "duration-between-servers-retries", c.HydraClient.DurationBetweenAllServersRetry, "")
	f.UintVar(&c.HydraClient.HydraServersCacheDuration, "hydra-servers-cache-duration", c.HydraClient.HydraServersCacheDuration, "")
	f.UintVar(&c.HydraClient.MaxNumberOfRetries, "max-number-of-retries", c.HydraClient.MaxNumberOfRetries, "")
	// f.BoolVar(&c.Verbose, "v", c.Verbose, "")
	// f.BoolVar(&c.Verbose, "verbose", c.Verbose, "")

	// BEGIN IGNORED FLAGS
	f.StringVar(&ignoredString, "config", "", "")
	// BEGIN IGNORED FLAGS

	if err := f.Parse(arguments); err != nil {
		return err
	}

	// Convert some parameters to lists.
	if hydraServers != "" {
		c.HydraServers = strings.Split(hydraServers, ",")
	}

	return nil
}
