package reverse_proxy

import (
	"net/http"

	"net/http/httputil"
	"net/url"

	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/innotech/hydra-go-clint/client"
)

type ReverseProxy interface {
	Run()
}

type HydraReverseProxy struct {
	AppId             string   `toml:"app_id"`
	HydraServers      []string `toml:"hydra_servers"`
	ProxyAddr         string   `toml:"proxy_addr"`
	HydraClientConfig struct {
		AppsCacheDuration              uint `toml:"apps_cache_duration"`
		DurationBetweenAllServersRetry uint `toml:"duration_between_all_servers_retry"`
		HydraServersCacheDuration      uint `toml:"hydra_servers_cache_duration"`
		MaxNumberOfRetries             uint `toml:"max_number_of_retries"`
	}
}

func (h *HydraReverseProxy) buildHydraClient() (*Client, error) {
	err := HydraClientFactory.Config(h.HydraServers)
	if err != nil {
		return nil, err
	}
	HydraClientFactory.WithAppsCacheDuration(h.HydraClientConfig.AppsCacheDuration).
		WithHydraServersCacheDuration(h.HydraClientConfig.HydraServersCacheDuration).
		WithMaxNumberOfRetriesPerHydraServer(h.HydraClientConfig.MaxNumberOfRetries).
		WaitBetweenAllServersRetry(h.HydraClientConfig.DurationBetweenAllServersRetry)
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

func (h *HydraReverseProxy) buildProxy() *httputil.ReverseProxy {
	hydraClient := h.buildHydraClient()
	director := func(req *http.Request) {
		//
		servers, err := hydraClient.Get()
		target, err := url.Parse(Get())
		if err != nil {
			panic(err.Error())
		}
		//
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

func (h *HydraReverseProxy) buildServer() (http.Server, error) {
	return http.Server{
		Addr:    h.ProxyAddr,
		Handler: h.buildProxy(),
	}
}

func (h *HydraReverseProxy) Run() {
	server, err := h.buildServer()
	if err != nil {
		panic(err.Error())
	}
	server.ListenAndServe()
}
