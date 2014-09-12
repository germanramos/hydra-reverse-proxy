package testingtsupport_test

import (
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/gomega"

	"testing"
)

func TestTestingT(t *testing.T) {
	RegisterTestingT(t)
	Ω(true).Should(BeTrue())
}
