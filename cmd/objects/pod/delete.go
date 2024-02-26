package pod

import (
	"fmt"

	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete pods in a namespace",
	Long:  "Delete specific pods in a namespace or all pods by status",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		status, _ := cmd.Flags().GetString("status")
		deletepod(namespace, status)
	},
}

func DeleteCmd() *cobra.Command {
	deleteCmd.Flags().StringP("namespace", "n", "", "Pass namespace having pods.")
	deleteCmd.Flags().StringP("status", "s", "Failed", "Delete all pods in a namespace by specific status(Failed/Running). If not passed then ask for specific pods.")
	deleteCmd.MarkFlagRequired("namespace")
	return deleteCmd
}

func deletepod(namespace string, status string) {
	if status == "" {
		podname := objects.SelectObject(namespace, "pods", "pod/", "name")
		stdout := cmd.ExecuteCommand("delete", "pods", podname, "-n", namespace)
		fmt.Println(string(stdout))
	} else {
		stdout := cmd.ExecuteCommand("delete", "pods", "--field-selector", "status.phase="+status, "-n", namespace)
		fmt.Println(string(stdout))
	}
}
