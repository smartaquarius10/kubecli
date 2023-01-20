package deployments

import (
	"fmt"

	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

var selectorCmd = &cobra.Command{
	Use:   "selector",
	Short: "Return traverse node for specified attribute",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		attrib, _ := cmd.Flags().GetString("attrib")
		getSelector(attrib, namespace)
	},
}

func SelectorCmd() *cobra.Command {
	selectorCmd.Flags().StringP("namespace", "n", "", "Pass namespace having deployments")
	selectorCmd.Flags().StringP("attrib", "a", "", "Pass attribute name")
	selectorCmd.MarkFlagRequired("namespace")
	selectorCmd.MarkFlagRequired("attrib")
	return selectorCmd
}

func getSelector(attribute string, namespace string) {
	query := cmd.GetJsonNodePath(objects.GetObjectJson(objects.SelectObject(namespace, "deployments", "deployment.apps/"), namespace, "deployments"), attribute)
	fmt.Println(query)

}
