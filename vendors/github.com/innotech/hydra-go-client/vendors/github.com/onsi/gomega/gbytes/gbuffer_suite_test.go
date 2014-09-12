package gbytes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"

	"testing"
)

func TestGbytes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gbytes Suite")
}
