package pod

import (
	"fmt"

	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete pod in a namespace",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		deletepod(namespace)
	},
}

func DeleteCmd() *cobra.Command {
	deleteCmd.Flags().StringP("namespace", "n", "", "Pass namespace having deployments")
	deleteCmd.MarkFlagRequired("namespace")
	return deleteCmd
}

func deletepod(namespace string) {
	podname := objects.SelectObject(namespace, "pods", "pod/", "name")
	stdout := cmd.ExecuteCommand("delete", "pods", podname, "-n", namespace)
	fmt.Println(string(stdout))
}
