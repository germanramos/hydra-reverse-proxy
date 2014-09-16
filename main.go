package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/innotech/hydra-reverse-proxy/log"
	"github.com/innotech/hydra-reverse-proxy/reverse_proxy"
)

func main() {
	proxy, err := HydraReverseProxyFactory.Build(os.Args[1:])
	if err != nil {
		log.Fatal(err.Error() + "\n")
	}
	proxy.Run()
}
