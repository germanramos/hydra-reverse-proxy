package reverse_proxy_test

import (
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/gomega"

	"testing"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reverse Proxy Suite")
}
