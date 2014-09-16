package reverse_proxy

import (
	// "errors"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/innotech/hydra-reverse-proxy/vendors/github.com/BurntSushi/toml"
)

const (
	DefaultConfigFilePath = "/etc/hydra-reverse-proxy.conf"
)

type HydraProxyBuilder struct {
	ConfigFilePath string
}

var HydraReverseProxyFactory *HydraProxyBuilder = &HydraProxyBuilder{
	ConfigFilePath: DefaultConfigFilePath,
}

var hydraReverseProxy *HydraReverseProxy

// Build configures hydra-reverse-proxy, it can be loaded from both system file,
// custom file or command line arguments and the values extracted from
// files they can be overriden with the command line arguments.
func (h *HydraProxyBuilder) Build(arguments []string) (*HydraReverseProxy, error) {
	var path string
	f := flag.NewFlagSet("hydra-reverse-proxy", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&path, "config", "", "path to config file")
	f.Parse(arguments)

	hydraReverseProxy = new(HydraReverseProxy)

	if path != "" {
		// Load from config file specified in arguments.
		if err := h.loadFile(path); err != nil {
			return nil, err
		}
	} else {
		// Load from system file.
		if err := h.loadSystemFile(); err != nil {
			return nil, err
		}
	}

	// Load from command line flags.
	if err := h.loadFlags(arguments); err != nil {
		return nil, err
	}

	// Validate proxy
	// if err := h.validateHydraProxy(); err != nil {
	// 	return nil, err
	// }

	return hydraReverseProxy, nil
}

// loadSystemFile loads from the default hydra-reverse-proxy configuration
// file path (/etc/hydra-reverse-proxy.conf) if it exists.
func (h *HydraProxyBuilder) loadSystemFile() error {
	if _, err := os.Stat(h.ConfigFilePath); os.IsNotExist(err) {
		return nil
	}
	return h.loadFile(h.ConfigFilePath)
}

// loadFile loads configuration from a file.
func (h *HydraProxyBuilder) loadFile(path string) error {
	_, err := toml.DecodeFile(path, &hydraReverseProxy)
	return err
}

// loadFlags loads configuration from command line flags.
func (h *HydraProxyBuilder) loadFlags(arguments []string) error {
	var hydraServers, ignoredString string

	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&hydraReverseProxy.AppId, "app-id", hydraReverseProxy.AppId, "")
	f.StringVar(&hydraServers, "hydra-servers", "", "")
	f.StringVar(&hydraReverseProxy.ProxyAddr, "proxy-addr", hydraReverseProxy.ProxyAddr, "")
	f.UintVar(&hydraReverseProxy.HydraClient.AppsCacheDuration, "apps-cache-duration", hydraReverseProxy.HydraClient.AppsCacheDuration, "")
	f.UintVar(&hydraReverseProxy.HydraClient.DurationBetweenAllServersRetry, "duration-between-servers-retries", hydraReverseProxy.HydraClient.DurationBetweenAllServersRetry, "")
	f.UintVar(&hydraReverseProxy.HydraClient.HydraServersCacheDuration, "hydra-servers-cache-duration", hydraReverseProxy.HydraClient.HydraServersCacheDuration, "")
	f.UintVar(&hydraReverseProxy.HydraClient.MaxNumberOfRetries, "max-number-of-retries", hydraReverseProxy.HydraClient.MaxNumberOfRetries, "")
	// f.BoolVar(&h.Verbose, "v", h.Verbose, "")
	// f.BoolVar(&h.Verbose, "verbose", h.Verbose, "")

	// BEGIN IGNORED FLAGS
	f.StringVar(&ignoredString, "config", "", "")
	// BEGIN IGNORED FLAGS

	if err := f.Parse(arguments); err != nil {
		return err
	}

	// Convert some parameters to lists.
	if hydraServers != "" {
		hydraReverseProxy.HydraServers = strings.Split(hydraServers, ",")
	}

	return nil
}

// func (h *HydraProxyBuilder) validateHydraProxy() error {
// 	if hydraReverseProxy.AppId == "" {
// 		return errors.New("application ID can not be empty")
// 	}
// 	if len(hydraReverseProxy.HydraServers) == 0 {
// 		return errors.New("It must be configured at least one hydra server")
// 	}
// 	if hydraReverseProxy.ProxyAddr == "" {
// 		return errors.New("It must be configured one valid endpoint address proxy")
// 	}
// 	return nil
// }
