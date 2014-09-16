package main

import (
	"github.com/innotech/hydra-reverse-proxy/log"
	. "github.com/innotech/hydra-reverse-proxy/reverse_proxy"

	"os"
)

func main() {
	proxy, err := HydraReverseProxyFactory.Build(os.Args[1:])
	if err != nil {
		log.Fatal(err.Error() + "\n")
	}
	proxy.Run()
}
