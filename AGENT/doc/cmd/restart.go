package cmd

import (
	"github.com/spf13/cobra"
)

var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart zhongxin server by daemon/pid file",
	Run: func(cmd *cobra.Command, args []string) {
		stop()
		start()
	},
}

func init() {
	RootCmd.AddCommand(RestartCmd)
}
