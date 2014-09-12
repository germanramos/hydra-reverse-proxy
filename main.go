package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Config struct {
	sourceAddr      string
	destinationAddr string
}

func main() {
	sourceAddress := ":3000"
	destinationUrlString := "http://127.0.0.1:7772"
	destinationUrl, _ := url.Parse(destinationUrlString)
	proxyHandler := httputil.NewSingleHostReverseProxy(destinationUrl)
	server := http.Server{
		Addr:    sourceAddress,
		Handler: proxyHandler,
	}
	server.ListenAndServe()
}

// loadAppsFromJSON reads application configuration from json file
// func (a *ApplicationsConfig) loadAppsFromJSON(pathToFile string) error {
// 	fileContent, err := ioutil.ReadFile(pathToFile)
// 	if err != nil {
// 		return err
// 	}
// 	if err := json.Unmarshal(fileContent, &(a.Apps)); err != nil {
// 		return err
// 	}
// 	return nil
// }
