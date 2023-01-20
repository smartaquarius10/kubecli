package deployments

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale replicas",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		filter, _ := cmd.Flags().GetStringArray("filter")
		count, _ := cmd.Flags().GetString("count")
		scaleobjects(namespace, filter, count)
	},
}

func ScaleCmd() *cobra.Command {
	scaleCmd.Flags().StringP("namespace", "n", "", "Pass namespace having deployments")
	scaleCmd.Flags().StringArrayP("filter", "f", []string{"*"}, "Pass search characters of deployment")
	scaleCmd.Flags().StringP("count", "c", "", "Pass replica count")
	scaleCmd.MarkFlagRequired("namespace")
	scaleCmd.MarkFlagRequired("count")
	return scaleCmd
}
func scaleobjects(namespace string, filter []string, count string) {
	stdout := objects.GetKubernetesObject(namespace, "deployments")
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))

	for scanner.Scan() {
		for _, f := range filter {
			if f == "*" || strings.Contains(scanner.Text(), f) {
				name := strings.ReplaceAll(scanner.Text(), ".apps", "")
				stdout = cmd.ExecuteCommand("scale", "--replicas="+count, name, "-n", namespace)
				fmt.Println(string(stdout))
			}
		}
	}
}
