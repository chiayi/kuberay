package e2e

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKubectlRayLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kubectl Ray Get")
}

var _ = Describe("Calling ray plugin `log` command on Ray Cluster", Ordered, func() {
	It("succeed in retrieving all ray cluster logs", func() {
		cmd := exec.Command("kubectl", "ray", "log", "raycluster-sample", "--node-type", "all")
		output, err := cmd.CombinedOutput()

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainElement("Testing 123"))
	})

	It("succeed in retrieving ray cluster head logs", func() {
		cmd := exec.Command("kubectl", "ray", "log", "raycluster-sample", "--node-type", "head")
		output, err := cmd.CombinedOutput()

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainElement("Testing 123"))
	})

	It("succeed in retrieving ray cluster worker logs", func() {
		cmd := exec.Command("kubectl", "ray", "log", "raycluster-sample", "--node-type", "worker")
		output, err := cmd.CombinedOutput()

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainElement("Testing 123"))
	})

	It("should not succeed", func() {
		cmd := exec.Command("kubectl", "ray", "log", "fakeclustername")
		output, err := cmd.CombinedOutput()

		Expect(err).To(HaveOccurred())
		Expect(output).ToNot(ContainElements("fakeclustername"))
	})
})
