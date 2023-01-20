package contextauth

import (
	"fmt"

	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "authdelete",
	Short: "Delete the auth context",
	Run: func(cmd *cobra.Command, args []string) {
		deleteContext()
	},
}

func deleteContext() {
	stdout := cmd.ExecuteCommand("config", "delete-context", getContext())
	fmt.Println(string(stdout))
}
