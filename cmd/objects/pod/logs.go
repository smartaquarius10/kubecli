package pod

import (
	"fmt"

	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "This command return pods logs",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		getLogsPodPerContainer(namespace)
	},
}

func LogsCmd() *cobra.Command {
	logsCmd.Flags().StringP("namespace", "n", "default", "Pass namespace having deployments")
	return logsCmd
}
func getLogsPodPerContainer(namespace string) {
	podname := objects.SelectObject(namespace, "pods", "pod/", "name")
	container := objects.SelectContainer(namespace, podname, `jsonpath='{range .spec.containers[*]}{.name}{"\n"}{end}'`, "'")
	stdout := cmd.ExecuteCommand("logs", podname, "-c", container, "-n", namespace)
	fmt.Println(string(stdout))
}
