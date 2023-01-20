package deployments

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "get object attributes value of deployment",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		filter, _ := cmd.Flags().GetStringArray("filter")
		attribute, _ := cmd.Flags().GetString("attrib")
		viewObject(namespace, attribute, filter)
	},
}

func ViewCmd() *cobra.Command {
	viewCmd.Flags().StringP("namespace", "n", "", "Pass namespace having deployments")
	viewCmd.Flags().StringArrayP("filter", "f", []string{"*"}, "Pass search characters or full name of deployment")
	viewCmd.Flags().StringP("attrib", "a", "", "Pass attribute name")
	viewCmd.MarkFlagRequired("namespace")
	viewCmd.MarkFlagRequired("attrib")
	return viewCmd
}

func viewObject(namespace string, attribute string, filter []string) {
	stdout := objects.GetKubernetesObject(namespace, "deployments")
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	var query string = ""
	for scanner.Scan() {
		for _, f := range filter {
			if f == "*" || strings.Contains(scanner.Text(), f) {
				name := objects.RemoveExtraChars(scanner.Text(), "deployment.apps/")
				if query == "" {
					stdoutj := objects.GetObjectJson(name, namespace, "deployments")
					query = cmd.GetJsonNodePath(stdoutj, attribute)
				}
				if query != "" {
					stdout = cmd.ExecuteCommand("get", "deployments", name, "-n", namespace, "-o", "jsonpath='{"+query+"}'")
					if len(stdout) > 2 {
						fmt.Println(name + " = " + strings.Trim(string(stdout), "'"))
					}
				}
			}
		}
	}
}
