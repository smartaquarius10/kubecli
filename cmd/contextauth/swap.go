package contextauth

import (
	"fmt"

	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/spf13/cobra"
)

var SwapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Change the auth context",
	Run: func(cmd *cobra.Command, args []string) {
		swapContext()
	},
}

func swapContext() {
	stdout := cmd.ExecuteCommand("config", "use-context", getContext())
	fmt.Println(string(stdout))
}
