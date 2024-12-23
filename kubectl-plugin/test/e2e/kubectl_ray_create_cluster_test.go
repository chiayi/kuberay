package e2e

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

var _ = Describe("Calling ray plugin `create cluster` command", Ordered, func() {
	It("succeed in creating RayCluster", func() {
		rayClusterName := "test-rayCluster"
		cmd := exec.Command("kubectl", "ray", "create", "cluster", rayClusterName)
		output, err := cmd.CombinedOutput()

		// Wait for cluster creation
		time.Sleep(15 * time.Second)

		// Check that the cluster is created
		Expect(err).NotTo(HaveOccurred())
		expectedOutput := fmt.Sprintf("Created Ray Cluster: %s", rayClusterName)
		Expect(strings.TrimSpace(string(output))).To(Equal(expectedOutput))

		cmd = exec.Command("kubectl", "ray", "get", "cluster", rayClusterName, "--namespace", "default")
		output, err = cmd.CombinedOutput()

		Expect(err).NotTo(HaveOccurred())
		expectedOutputTablePrinter := printers.NewTablePrinter(printers.PrintOptions{})
		expectedTestResultTable := &v1.Table{
			ColumnDefinitions: []v1.TableColumnDefinition{
				{Name: "Name", Type: "string"},
				{Name: "Namespace", Type: "string"},
				{Name: "Desired Workers", Type: "string"},
				{Name: "Available Workers", Type: "string"},
				{Name: "CPUs", Type: "string"},
				{Name: "GPUs", Type: "string"},
				{Name: "TPUs", Type: "string"},
				{Name: "Memory", Type: "string"},
				{Name: "Age", Type: "string"},
			},
		}

		expectedTestResultTable.Rows = append(expectedTestResultTable.Rows, v1.TableRow{
			Cells: []interface{}{
				rayClusterName,
				"default",
				"1",
				"1",
				"4",
				"0",
				"0",
				"8G",
			},
		})

		var resbuffer bytes.Buffer
		bufferr := expectedOutputTablePrinter.PrintObj(expectedTestResultTable, &resbuffer)
		Expect(bufferr).NotTo(HaveOccurred())

		Expect(err).NotTo(HaveOccurred())
		Expect(strings.TrimSpace(string(output))).To(ContainSubstring(strings.TrimSpace(resbuffer.String())))

		// Delete raycluster
		cmd = exec.Command("kubectl", "ray", "delete", "cluster", rayClusterName)
		cmdStdin, err := cmd.StdinPipe()
		Expect(err).NotTo(HaveOccurred())

		go func() {
			defer cmdStdin.Close()
			_, err := io.WriteString(cmdStdin, "yes\n")
			Expect(err).NotTo(HaveOccurred())
		}()

		output, err = cmd.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Delete raycluster %s", rayClusterName))

		// Wait for cluster delete to process
		time.Sleep(15 * time.Second)

		cmd = exec.Command("kubectl", "ray", "get", "cluster", rayClusterName, "--namespace", "default")
		output, err = cmd.CombinedOutput()

		Expect(err).To(HaveOccurred())
		Expect(output).ToNot(ContainElements(rayClusterName))
	})
})
