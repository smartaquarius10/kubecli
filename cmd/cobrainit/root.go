package cobrainit

import (
	"fmt"
	"os"

	"github.com/smartaquarius10/kubecli/cmd/contextauth"

	"github.com/smartaquarius10/kubecli/cmd/objects/deployments"
	"github.com/smartaquarius10/kubecli/cmd/objects/pod"
	"github.com/smartaquarius10/kubecli/cmd/objects/resources"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "kubectl command",
		Short: "Commands to manage kubernetes",
		Long:  `Kubectl is necessary on the machine`,
	}
)

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(contextauth.SwapCmd)
	rootCmd.AddCommand(contextauth.DeleteCmd)
	rootCmd.AddCommand(deployments.BackupCmd())
	rootCmd.AddCommand(deployments.ApplyCmd())
	rootCmd.AddCommand(deployments.RemoveCmd())
	rootCmd.AddCommand(deployments.ViewCmd())
	rootCmd.AddCommand(deployments.ScaleCmd())
	rootCmd.AddCommand(deployments.SelectorCmd())
	rootCmd.AddCommand(pod.DeleteCmd())
	rootCmd.AddCommand(deployments.RestartCmd())
	rootCmd.AddCommand(pod.LogsCmd())
	rootCmd.AddCommand(pod.TopCmd())
	rootCmd.AddCommand(resources.TsCmd())
}
