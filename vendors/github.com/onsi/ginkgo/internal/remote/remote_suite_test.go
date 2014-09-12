package remote_test

import (
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/gomega"

	"testing"
)

func TestRemote(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Remote Spec Forwarding Suite")
}
