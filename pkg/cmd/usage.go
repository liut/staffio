package cmd

import (
	"github.com/spf13/cobra"

	config "github.com/liut/staffio/pkg/settings"
)

// usageCmd represents the usage command
var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Print usage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config.Usage()
	},
}

func init() {
	RootCmd.AddCommand(usageCmd)
}
