package e2e

import (
	"testing"

	. "github.com/shipengqi/example.v1/cli/test/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e")
}

var t *TestFramework

var _ = BeforeSuite(func() {
	t = NewTestFramework("")
})

var _ = AfterSuite(func() {
	t.Teardown()
})
