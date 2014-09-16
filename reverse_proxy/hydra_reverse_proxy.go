package reverse_proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/innotech/hydra-go-client/client"

	"github.com/innotech/hydra-reverse-proxy/log"
)

type ReverseProxy interface {
	Run()
}

type HydraReverseProxy struct {
	AppId        string   `toml:"app_id"`
	HydraServers []string `toml:"hydra_servers"`
	ProxyAddr    string   `toml:"proxy_addr"`
	HydraClient  struct {
		AppsCacheDuration              uint `toml:"apps_cache_duration"`
		DurationBetweenAllServersRetry uint `toml:"duration_between_all_servers_retry"`
		HydraServersCacheDuration      uint `toml:"hydra_servers_cache_duration"`
		MaxNumberOfRetries             uint `toml:"max_number_of_retries"`
	}
}

// buildHydraClient builds a hydra client configured completely.
func (h *HydraReverseProxy) buildHydraClient() *Client {
	if err := HydraClientFactory.Config(h.HydraServers); err != nil {
		log.Fatal(err.Error())
	}
	HydraClientFactory.WithAppsCacheDuration(time.Duration(h.HydraClient.AppsCacheDuration) * time.Millisecond).
		WithHydraServersCacheDuration(time.Duration(h.HydraClient.HydraServersCacheDuration) * time.Millisecond).
		WithMaxNumberOfRetriesPerHydraServer(h.HydraClient.MaxNumberOfRetries).
		WaitBetweenAllServersRetry(time.Duration(h.HydraClient.DurationBetweenAllServersRetry) * time.Millisecond)

	return HydraClientFactory.Build()
}

func (h *HydraReverseProxy) singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func (h *HydraReverseProxy) getURLTarget(serverURLs []string) *url.URL {
	var target *url.URL
	var err error
	for _, serverURL := range serverURLs {
		target, err = url.Parse(serverURL)
		if err != nil {
			log.Warn("Bad server URL " + serverURL)
		} else {
			return target
		}
	}
	target, _ = url.Parse("")
	return target
}

// buildProxy builds a reverse proxy which redirects requests from clients to
// one of the servers linked to the destination service monitored by the hydra client.
func (h *HydraReverseProxy) buildProxy() *httputil.ReverseProxy {
	hydraClient := h.buildHydraClient()
	director := func(req *http.Request) {
		var target *url.URL
		serverURLs, err := hydraClient.Get(h.AppId, false)
		if err != nil || err == nil && len(serverURLs) == 0 {
			target, _ = url.Parse("")
		} else {
			target = h.getURLTarget(serverURLs)
		}
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = h.singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

// buildServer builds and returns a http server which is handled by the proxy.
func (h *HydraReverseProxy) buildServer() http.Server {
	return http.Server{
		Addr:    h.ProxyAddr,
		Handler: h.buildProxy(),
	}
}

// Run launches the proxy server and kept listening requests.
func (h *HydraReverseProxy) Run() {
	server := h.buildServer()
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
