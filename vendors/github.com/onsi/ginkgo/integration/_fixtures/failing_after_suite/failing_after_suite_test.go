package failing_before_suite_test

import (
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/onsi/ginkgo"
)

var _ = Describe("FailingBeforeSuite", func() {
	It("should run", func() {
		println("A TEST")
	})

	It("should run", func() {
		println("A TEST")
	})
})
