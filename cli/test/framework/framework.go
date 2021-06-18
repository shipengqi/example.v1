package framework

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type TestFramework struct {}

// NewTestFramework creates a new test framework instance
func NewTestFramework() *TestFramework {
	return &TestFramework{}
}

// Setup is the global initialization function which runs before each test
// suite
func (t *TestFramework) Setup(dir string) {
	// Global initialization for the whole framework goes in here
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(GinkgoWriter)
}

// Teardown is the global initialization function which runs after each test
// suite
func (t *TestFramework) Teardown() {
}

// Describe is a convenience wrapper around the `ginkgo.Describe` function
func (t *TestFramework) Describe(text string, body func()) bool {
	return Describe("cert-manager: "+text, body)
}

// Convenience method for command creation
func cmd(workDir, format string, args ...interface{}) *Session {
	c := strings.Split(fmt.Sprintf(format, args...), " ")
	command := exec.Command(c[0], c[1:]...)
	if workDir != "" {
		command.Dir = workDir
	}

	session, err := Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).To(BeNil())

	return session
}

// Convenience method for command creation in the current working directory
func lcmd(format string, args ...interface{}) *Session {
	return cmd("", format, args...)
}

// Run cert-manager and return the resulting session
func (t *TestFramework) CertManager(args string) *Session {
	return lcmd("cert-manager %s", args).Wait()
}

// Run cert-manager and expect success containing the specified output
func (t *TestFramework) CertManagerExpectSuccess(args, expectedOut string) {
	// When
	res := t.CertManager(args)

	// Then
	Expect(res).To(Exit(0))
	Expect(res.Out).To(Say(expectedOut))
	Expect(string(res.Err.Contents())).To(BeEmpty())
}

// Run cert-manager and expect error containing the specified outputs
func (t *TestFramework) CertManagerExpectFailure(
	args string, expectedOut, expectedErr string,
) {
	// When
	res := t.CertManager(args)

	// Then
	Expect(res).To(Exit(1))
	Expect(res.Out).To(Say(expectedOut))
	Expect(res.Err).To(Say(expectedErr))
}
