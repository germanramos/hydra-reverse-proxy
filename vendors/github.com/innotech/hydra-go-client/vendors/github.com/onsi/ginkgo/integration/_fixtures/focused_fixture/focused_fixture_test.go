package focused_fixture_test

import (
	. "github.com/innotech/hydra-reverse-proxy/vendors/github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
)

var _ = Describe("FocusedFixture", func() {
	FDescribe("focused", func() {
		It("focused", func() {

		})
	})

	FContext("focused", func() {
		It("focused", func() {

		})
	})

	FIt("focused", func() {

	})

	Describe("not focused", func() {
		It("not focused", func() {

		})
	})

	Context("not focused", func() {
		It("not focused", func() {

		})
	})

	It("not focused", func() {

	})
})
