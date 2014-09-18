package reverse_proxy_test

import (
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/gomega"

	. "github.com/innotech/hydra-reverse-proxy/reverse_proxy"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("Hydra-Reverse-Proxy", func() {
	Context("when none hydra server is accesible", func() {
		It("should not be able to run", func() {
			appId := "app1"

			// SERVICE MOCK
			var serviceResponseBody []byte = []byte("SUCCESS")
			servicePath := "/text/success"
			serviceMux := http.NewServeMux()
			serviceMux.Handle(servicePath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "text/plain")
				w.Write(serviceResponseBody)
			}))
			ts1 := httptest.NewServer(serviceMux)
			defer ts1.Close()
			time.Sleep(time.Duration(200) * time.Millisecond)
			// END SERVICE MOCK

			hydraServerAddr := "http://localhost:7772"
			// PROXY
			proxyAddr := ":3000"
			proxy, err := HydraReverseProxyFactory.Build([]string{"-app-id", appId, "-hydra-servers", hydraServerAddr, "-proxy-addr", proxyAddr})
			Expect(err).ToNot(HaveOccurred())
			go func() {
				proxy.Run()
			}()
			time.Sleep(time.Duration(200) * time.Millisecond)
			// END PROXY

			res, err := http.Get("http://localhost" + proxyAddr + servicePath)
			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())
		})
	})
	Context("when at least one hydra server is accesible", func() {
		Context("when first service is not accesible", func() {
			Context("when client request is correct", func() {
				It("should response with an internal server error", func() {
					appPath := "/app/"
					appId := "app1"

					noAccessibleServiceAddr := "http://localhost:9999"

					// SERVICE MOCK
					var serviceResponseBody []byte = []byte("SECOND SERVICE REACHED")
					servicePath := "/text/success"
					serviceMux := http.NewServeMux()
					serviceMux.Handle(servicePath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "text/plain")
						w.Write(serviceResponseBody)
					}))
					ts1 := httptest.NewServer(serviceMux)
					defer ts1.Close()
					time.Sleep(time.Duration(200) * time.Millisecond)
					// END SERVICE MOCK

					// HYDRA MOCK
					hydraId := "hydra"
					var hydraServerAddr *string
					hydraMux := http.NewServeMux()
					hydraMux.Handle(appPath+hydraId, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						sortedInstanceUris := []string{*hydraServerAddr}
						jsonOutput, _ := json.Marshal(sortedInstanceUris)
						w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "application/json")
						w.Write(jsonOutput)
					}))
					hydraMux.Handle(appPath+appId, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						sortedInstanceUris := []string{noAccessibleServiceAddr, ts1.URL}
						jsonOutput, _ := json.Marshal(sortedInstanceUris)
						w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "application/json")
						w.Write(jsonOutput)
					}))
					ts2 := httptest.NewServer(hydraMux)
					hydraServerAddr = &ts2.URL
					defer ts2.Close()
					time.Sleep(time.Duration(200) * time.Millisecond)
					// END HYDRA MOCK

					// PROXY
					proxyAddr := ":3000"
					proxy, err := HydraReverseProxyFactory.Build([]string{"-app-id", appId, "-hydra-servers", *hydraServerAddr, "-proxy-addr", proxyAddr})
					Expect(err).ToNot(HaveOccurred())
					// wg.Add(1)
					go func() {
						proxy.Run()
					}()
					time.Sleep(time.Duration(200) * time.Millisecond)
					// END PROXY

					res, err := http.Get("http://localhost" + proxyAddr + servicePath)
					Expect(err).ToNot(HaveOccurred())
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})
		Context("when first service is accesible", func() {
			Context("when client request is correct", func() {
				It("should route the response successfully", func() {
					appPath := "/app/"
					appId := "app1"

					// SERVICE MOCK
					var serviceResponseBody []byte = []byte("SUCCESS")
					servicePath := "/text/success"
					serviceMux := http.NewServeMux()
					serviceMux.Handle(servicePath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "text/plain")
						w.Write(serviceResponseBody)
					}))
					ts1 := httptest.NewServer(serviceMux)
					defer ts1.Close()
					time.Sleep(time.Duration(200) * time.Millisecond)
					// END SERVICE MOCK

					// HYDRA MOCK
					hydraId := "hydra"
					var hydraServerAddr *string
					hydraMux := http.NewServeMux()
					hydraMux.Handle(appPath+hydraId, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						sortedInstanceUris := []string{*hydraServerAddr}
						jsonOutput, _ := json.Marshal(sortedInstanceUris)
						w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "application/json")
						w.Write(jsonOutput)
					}))
					hydraMux.Handle(appPath+appId, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						sortedInstanceUris := []string{ts1.URL}
						jsonOutput, _ := json.Marshal(sortedInstanceUris)
						w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "application/json")
						w.Write(jsonOutput)
					}))
					ts2 := httptest.NewServer(hydraMux)
					hydraServerAddr = &ts2.URL
					defer ts2.Close()
					time.Sleep(time.Duration(200) * time.Millisecond)
					// END HYDRA MOCK

					// PROXY
					proxyAddr := ":3001"
					proxy, err := HydraReverseProxyFactory.Build([]string{"-app-id", appId, "-hydra-servers", *hydraServerAddr, "-proxy-addr", proxyAddr})
					Expect(err).ToNot(HaveOccurred())
					go func() {
						proxy.Run()
					}()
					time.Sleep(time.Duration(200) * time.Millisecond)
					// END PROXY

					res, err := http.Get("http://localhost" + proxyAddr + servicePath)
					Expect(err).ToNot(HaveOccurred())
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					bodyText, err := ioutil.ReadAll(res.Body)
					res.Body.Close()
					Expect(err).ToNot(HaveOccurred())
					Expect(bodyText).To(Equal(serviceResponseBody))
				})
			})
			Context("when client request is incorrect", func() {
				It("should route the response successfully", func() {
					appPath := "/app/"
					appId := "app1"

					// SERVICE MOCK
					var serviceResponseBody []byte = []byte("SUCCESS")
					servicePath := "/text/success"
					serviceMux := http.NewServeMux()
					serviceMux.Handle(servicePath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "text/plain")
						w.Write(serviceResponseBody)
					}))
					ts1 := httptest.NewServer(serviceMux)
					defer ts1.Close()
					time.Sleep(time.Duration(200) * time.Millisecond)
					// END SERVICE MOCK

					// HYDRA MOCK
					hydraId := "hydra"
					var hydraServerAddr *string
					hydraMux := http.NewServeMux()
					hydraMux.Handle(appPath+hydraId, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						sortedInstanceUris := []string{*hydraServerAddr}
						jsonOutput, _ := json.Marshal(sortedInstanceUris)
						w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "application/json")
						w.Write(jsonOutput)
					}))
					hydraMux.Handle(appPath+appId, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						sortedInstanceUris := []string{ts1.URL}
						jsonOutput, _ := json.Marshal(sortedInstanceUris)
						w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "application/json")
						w.Write(jsonOutput)
					}))
					ts2 := httptest.NewServer(hydraMux)
					hydraServerAddr = &ts2.URL
					defer ts2.Close()
					time.Sleep(time.Duration(200) * time.Millisecond)
					// END HYDRA MOCK

					// PROXY
					proxyAddr := ":3002"
					proxy, err := HydraReverseProxyFactory.Build([]string{"-app-id", appId, "-hydra-servers", *hydraServerAddr, "-proxy-addr", proxyAddr})
					Expect(err).ToNot(HaveOccurred())
					go func() {
						proxy.Run()
					}()
					time.Sleep(time.Duration(200) * time.Millisecond)
					// END PROXY

					badServicePath := "/text/failure"
					res, err := http.Get("http://localhost" + proxyAddr + badServicePath)
					Expect(err).ToNot(HaveOccurred())
					Expect(res.StatusCode).To(Equal(http.StatusNotFound))
				})
			})
		})
	})
})
