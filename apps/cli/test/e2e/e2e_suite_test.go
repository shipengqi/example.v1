package e2e

import (
	. "github.com/shipengqi/example.v1/apps/cli/test/framework"
	"testing"

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
