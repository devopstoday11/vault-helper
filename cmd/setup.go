package cmd

import (
	"github.com/spf13/cobra"

	//"github.com/Sirupsen/logrus"
	"gitlab.jetstack.net/jetstack-experimental/vault-helper/pkg/kubernetes"
)

var MaxTTL string

// initCmd represents the init command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup kubernetes on a running vault server",
	Run: func(cmd *cobra.Command, args []string) {

		kubernetes.Run(cmd, args)

	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&MaxTTL, "MaxTTL", "", "Maxium Validity CA")
	RootCmd.AddCommand(setupCmd)
}
