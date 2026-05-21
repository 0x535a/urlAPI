package cmd

import (
	"fmt"
	"os"
	"zhongxin/cmd/flags"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "zhongxin",
	Short: "zhongxin",
	Long:  "The backend of an AI-powered zhongxinee-assistant software",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
}

func init() {
	exePath, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	RootCmd.PersistentFlags().StringVar(&flags.DataDir, "data", exePath+"/data", "data folder")
	RootCmd.PersistentFlags().BoolVar(&flags.Dev, "dev", false, "start with dev mode")
	RootCmd.PersistentFlags().BoolVar(&flags.Beta, "beta", false, "start with beta mode")
	RootCmd.PersistentFlags().BoolVar(&flags.NoRemote, "no-remote", false, "start without remote db")
	RootCmd.PersistentFlags().BoolVar(&flags.Verbose, "verbose", false, "show more logs")
	AdminCmd.PersistentFlags().StringVar(&name, "name", "", "admin name")
	AdminCmd.PersistentFlags().StringVar(&phone, "phone", "", "phone number")
	AdminCmd.PersistentFlags().StringVar(&wxid, "wxid", "", "wxid")
}
