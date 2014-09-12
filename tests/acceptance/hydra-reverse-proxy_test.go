package acceptance_test

import (
	// . "github.com/innotech/hydra-reverse-proxy/tests/acceptance"

	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/gomega"

	. "github.com/innotech/hydra-reverse-proxy"

	"time"
)

var _ = Describe("Hydra-Reverse-Proxy", func() {
	Context("when first service is accesible", func() {
		Context("when client request is correct", func() {
			It("should route the response successfully", func() {
				serviceServers := []string{"http://www.q1.com", "http://www.q2.com", "http://www.q3.com"}
				servicePath := "/app/app1"

				ts := httptest.NewServer(http.HandlerFunc(servicePath, func(w http.ResponseWriter, r *http.Request) {
					sortedInstanceUris := []string{"http://www.q1.com", "http://www.q2.com", "http://www.q3.com"}
					jsonOutput, _ = json.Marshal(sortedInstanceUris)
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonOutput)
				}))
				defer ts.Close()

				// PROXY
				var proxyAddr string = ":3000"
				go func() {
					proxy := NewHydraReverseProxy(proxyAddr, ts.URL)
					proxy.Run()
				}()
				time.Sleep(time.Duration(200) * time.Millisecond)
				// END PROXY

				res, err := http.Get("http://localhost" + proxyAddr + servicePath)
				if err != nil {
					log.Fatal(err)
				}
				instances, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				if err != nil {
					log.Fatal(err)
				}

				Expect(instances).ToNot(BeNil())
				Expect(instances).To(HaveLen(3))
				Expect(instances).To(ConsistOf(serviceServers))
			})
		})
	})
})
