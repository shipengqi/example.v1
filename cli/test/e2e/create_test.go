package e2e

import (
	. "github.com/onsi/ginkgo"
)

// The actual test suite
var _ = t.Describe("create", func() {

	BeforeEach(func() {})

	AfterEach(func() {})

	It("should succeed with `create` subcommand", func() {
		t.ExecuteExpectSuccess("create", "")
	})

	It("should succeed with `--install` flag", func() {
		t.ExecuteExpectSuccess("--install", "")
	})
})
