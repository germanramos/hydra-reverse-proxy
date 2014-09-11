package config_test

import (
	. "github.com/innotech/hydra-reverse-proxy/config"
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"strconv"
)

func WithTempFile(content string, fn func(string)) {
	f, _ := ioutil.TempFile("", "")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())
	fn(f.Name())
}

var _ = Describe("Config", func() {
	Describe("loading from TOML", func() {
		Context("when the TOML file exists", func() {
			const (
				APP_ID                             string = "app1"
				HYDRA_SERVER_1                     string = "127.0.0.1:4000"
				HYDRA_SERVER_2                     string = "127.0.0.1:4001"
				PROXY_ADDR                         string = ":3000"
				APPS_CACHE_DURATION                uint   = 10000
				DURATION_BETWEEN_ALL_SERVERS_RETRY uint   = 15000
				HYDRA_SERVERS_CACHE_DURATION       uint   = 20000
				MAX_NUMBER_OF_RETRIES              uint   = 25000
			)
			fileContent := `
				app_id = "` + APP_ID + `"
				hydra_servers = ["` + HYDRA_SERVER_1 + `","` + HYDRA_SERVER_2 + `"]
				proxy_addr = "` + PROXY_ADDR + `"
				[hydraClient]
				apps_cache_duration = ` + strconv.FormatInt(int64(APPS_CACHE_DURATION), 10) + `
				duration_between_all_servers_retry = ` + strconv.FormatInt(int64(DURATION_BETWEEN_ALL_SERVERS_RETRY), 10) + `
				hydra_servers_cache_duration = ` + strconv.FormatInt(int64(HYDRA_SERVERS_CACHE_DURATION), 10) + `
				max_number_of_retries = ` + strconv.FormatInt(int64(MAX_NUMBER_OF_RETRIES), 10) + `
			`
			WithTempFile(fileContent, func(pathToFile string) {
				c := New()
				err := c.LoadFile(pathToFile)
				It("should be loaded successfully", func() {
					Expect(err).To(BeNil(), "error should be nil")
					Expect(c.AppId).To(Equal(APP_ID))
					Expect(c.HydraServers).To(ContainElement(HYDRA_SERVER_1))
					Expect(c.HydraServers).To(ContainElement(HYDRA_SERVER_2))
					Expect(c.ProxyAddr).To(Equal(PROXY_ADDR))
					Expect(c.HydraClient.AppsCacheDuration).To(Equal(APPS_CACHE_DURATION))
					Expect(c.HydraClient.DurationBetweenAllServersRetry).To(Equal(DURATION_BETWEEN_ALL_SERVERS_RETRY))
					Expect(c.HydraClient.HydraServersCacheDuration).To(Equal(HYDRA_SERVERS_CACHE_DURATION))
					Expect(c.HydraClient.MaxNumberOfRetries).To(Equal(MAX_NUMBER_OF_RETRIES))
				})
			})
		})
	})

	// Check system configuration values
	Describe("loading without flags", func() {
		Context("when default system cofig file exists", func() {
			systemProxyAddr := "127.0.0.1:77710"
			systemFileContent := `proxy_addr = "` + systemProxyAddr + `"`
			WithTempFile(systemFileContent, func(pathToSystemFile string) {
				c := New()
				c.ConfigFilePath = pathToSystemFile
				err := c.Load([]string{})
				It("should be loaded successfully", func() {
					Expect(err).ToNot(HaveOccurred())
				})
				It("should be override the default configuration", func() {
					Expect(c.ProxyAddr).To(Equal(systemProxyAddr), "c.ProxyAddr should be equal "+systemProxyAddr)
				})
			})
		})
	})

	Describe("loading from flags", func() {
		Context("when bad flag exists", func() {
			c := New()
			err := c.LoadFlags([]string{"-bad-flag"})
			It("should be throw an error", func() {
				Expect(err).To(HaveOccurred(), "No bad flag are allowed")
			})
			It("should be have an specific error message", func() {
				Expect(err.Error()).To(Equal("flag provided but not defined: -bad-flag"))
			})
		})
		Context("When -app-id flag exists", func() {
			const APP_ID string = "app35"
			c := New()
			err := c.LoadFlags([]string{"-app-id", APP_ID})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.AppId).To(Equal(APP_ID))
			})
		})
		Context("When -hydra-servers flag exists", func() {
			const HYDRA_SERVER_1 string = "203.0.113.101:7001"
			const HYDRA_SERVER_2 string = "203.0.113.102:7001"
			c := New()
			err := c.LoadFlags([]string{"-hydra-servers", HYDRA_SERVER_1 + "," + HYDRA_SERVER_2})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.HydraServers).To(HaveLen(2))
				Expect(c.HydraServers).To(ContainElement(HYDRA_SERVER_1))
				Expect(c.HydraServers).To(ContainElement(HYDRA_SERVER_2))
			})
		})
		Context("When -proxy-addr flag exists", func() {
			const PROXY_ADDR string = "127.0.0.1:7444"
			c := New()
			err := c.LoadFlags([]string{"-proxy-addr", PROXY_ADDR})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.ProxyAddr).To(Equal(PROXY_ADDR))
			})
		})
		Context("When -apps-cache-duration flag exists", func() {
			const APPS_CACHE_DURATION uint = 50000
			c := New()
			err := c.LoadFlags([]string{"-apps-cache-duration", strconv.FormatUint(uint64(APPS_CACHE_DURATION), 10)})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.HydraClient.AppsCacheDuration).To(Equal(APPS_CACHE_DURATION))
			})
		})
		Context("When -duration-between-servers-retries flag exists", func() {
			const DURATION_BETWEEN_SERVERS_RETRIES uint = 51000
			c := New()
			err := c.LoadFlags([]string{"-duration-between-servers-retries", strconv.FormatUint(uint64(DURATION_BETWEEN_SERVERS_RETRIES), 10)})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.HydraClient.DurationBetweenAllServersRetry).To(Equal(DURATION_BETWEEN_SERVERS_RETRIES))
			})
		})
		Context("When -hydra-servers-cache-duration flag exists", func() {
			const HYDRA_SERVERS_CACHE_DURATION uint = 53000
			c := New()
			err := c.LoadFlags([]string{"-hydra-servers-cache-duration", strconv.FormatUint(uint64(HYDRA_SERVERS_CACHE_DURATION), 10)})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.HydraClient.HydraServersCacheDuration).To(Equal(HYDRA_SERVERS_CACHE_DURATION))
			})
		})
		Context("When -max-number-of-retries flag exists", func() {
			const MAX_NUMBER_OF_RETRIES uint = 3
			c := New()
			err := c.LoadFlags([]string{"-max-number-of-retries", strconv.FormatUint(uint64(MAX_NUMBER_OF_RETRIES), 10)})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.HydraClient.MaxNumberOfRetries).To(Equal(MAX_NUMBER_OF_RETRIES))
			})
		})
		// Context("When -verbose flag exists", func() {
		// 	c := New()
		// 	c.Verbose = false
		// 	err := c.LoadFlags([]string{"-verbose"})
		// 	It("should be loaded successfully", func() {
		// 		Expect(err).NotTo(HaveOccurred())
		// 		Expect(c.Verbose).To(BeTrue())
		// 	})
		// })
		Context("when -config flag exists", func() {
			Context("and no more flags exist", func() {
				systemProxyAddr := "127.0.0.1:87720"
				systemFileContent := `proxy_addr = "` + systemProxyAddr + `"`
				customProxyAddr := systemProxyAddr + "0"
				customFileContent := `proxy_addr = "` + customProxyAddr + `"`
				WithTempFile(systemFileContent, func(pathToSystemFile string) {
					WithTempFile(customFileContent, func(pathToCustomFile string) {
						c := New()
						c.ConfigFilePath = pathToSystemFile
						err := c.Load([]string{"-config", pathToCustomFile})
						It("should be loaded successfully", func() {
							Expect(err).To(BeNil(), "error should be nil")
							Expect(c.ConfigFilePath).To(Equal(pathToSystemFile))
						})
						It("should be override the default configuration loaded from default system configuration file", func() {
							Expect(c.ProxyAddr).To(Equal(customProxyAddr), "c.ProxyAddr should be equal "+customProxyAddr)
						})
					})
				})
			})
			Context("and also more valid flags exist", func() {
				customProxyAddr := "127.0.0.1:87720"
				customFileContent := `proxy_addr = "` + customProxyAddr + `"`
				proxyAddrCustomFlag := customProxyAddr + "0"
				WithTempFile(customFileContent, func(pathToCustomFile string) {
					c := New()
					err := c.Load([]string{"-proxy-addr", proxyAddrCustomFlag, "-config", pathToCustomFile})
					It("should be loaded successfully", func() {
						Expect(err).To(BeNil(), "error should be nil")
					})
					It("should be override the configuration loaded from custom configuration file", func() {
						Expect(c.ProxyAddr).To(Equal(proxyAddrCustomFlag), "c.ProxyAddr should be equal "+proxyAddrCustomFlag)
					})
				})
			})
		})
		Context("when default system cofig file doesn't exist", func() {
			systemProxyAddr := "127.0.0.1:87720"
			systemFileContent := `proxy_addr = "` + systemProxyAddr + `"`
			customPublicAPIAddr := systemProxyAddr + "0"
			WithTempFile(systemFileContent, func(pathToSystemFile string) {
				c := New()
				c.ConfigFilePath = pathToSystemFile
				err := c.Load([]string{"-proxy-addr", customPublicAPIAddr})
				It("should be loaded successfully", func() {
					Expect(err).To(BeNil(), "error should be nil")
				})
				It("should be override the default configuration loaded from default system configuration file", func() {
					Expect(c.ProxyAddr).To(Equal(customPublicAPIAddr), "c.ProxyAddr should be equal "+customPublicAPIAddr)
				})
			})
		})
	})
})
