package resources

import (
	"fmt"

	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

var tsCmd = &cobra.Command{
	Use:   "time",
	Short: "This command return last updated timestamp of resource",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		objectType, _ := cmd.Flags().GetString("objectType")
		objectName, _ := cmd.Flags().GetString("objectName")
		getObjectLastUpdatedTimeStamp(namespace, objectType, objectName)
	},
}

func TsCmd() *cobra.Command {
	tsCmd.Flags().StringP("namespace", "n", "default", "Pass namespace having pods. If blank then default")
	tsCmd.Flags().StringP("objectType", "t", "", "Pass type of object. Foeg. secrets, deployment")
	tsCmd.Flags().StringP("objectName", "o", "", "Pass name of object")
	tsCmd.MarkFlagRequired("objectType")
	tsCmd.MarkFlagRequired("objectName")
	return tsCmd
}

func getObjectLastUpdatedTimeStamp(namespace string, objectType string, objectName string) {
	fmt.Println(objects.GetObjectLastUpdatedTimeStamp(objectType, objectName, namespace))
}
