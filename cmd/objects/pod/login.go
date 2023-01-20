package pod

import (
	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "This command login into pod",
	Long:  "Pass namespace & attribute. For eg. qa1 and nodetype",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		container, _ := cmd.Flags().GetString("container")
		loginPod(namespace, container)
	},
}

func LoginCmd() *cobra.Command {
	loginCmd.Flags().StringP("namespace", "n", "", "Pass namespace having deployments")
	loginCmd.Flags().StringP("container", "c", "", "Pass container name")
	loginCmd.MarkFlagRequired("namespace")
	return loginCmd
}
func loginPod(namespace string, container string) {
	podname := objects.SelectObject(namespace, "pods", "pod/")
	if container == "" {
		cmd.ExecuteSessionCommand("exec", podname, "-n", namespace, "-it", "--", "sh")
	} else {
		cmd.ExecuteSessionCommand("exec", podname, "-c", container, "-n", namespace, "-it", "--", "sh")
	}

}
