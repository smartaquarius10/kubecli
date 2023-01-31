package deployments

import (
	"fmt"

	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart deployment",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		restartObject(namespace)
	},
}

func RestartCmd() *cobra.Command {
	restartCmd.Flags().StringP("namespace", "n", "", "Pass namespace having deployments")
	restartCmd.MarkFlagRequired("namespace")
	return restartCmd
}

func restartObject(namespace string) {

	deployment := objects.SelectObject(namespace, "deployments", "deployment.apps/", "name")
	stdout := cmd.ExecuteCommand("rollout", "restart", "deployment", deployment, "-n", namespace)
	fmt.Println(string(stdout))
}
