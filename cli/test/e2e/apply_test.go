package e2e

import (
	. "github.com/onsi/ginkgo"
)

// The actual test suite
var _ = t.Describe("apply", func() {

	BeforeEach(func() {})

	AfterEach(func() {})

	It("should succeed with `apply` subcommand", func() {
		t.ExecuteExpectSuccess("apply", "")
	})

	It("should succeed with `--apply` flag", func() {
		t.ExecuteExpectSuccess("--apply", "")
	})
})
