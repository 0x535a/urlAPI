package cmd

import (
	"fmt"
	"strconv"
	"zhongxin/internal/conf"
	"zhongxin/internal/op"
	"zhongxin/util"

	"github.com/spf13/cobra"
)

var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Run: func(cmd *cobra.Command, args []string) {
		testCmd(args)
	},
}

func testCmd(args []string) {
	Init()
	defer Release()
	if len(args) < 1 {
		fmt.Println("please provide action")
		return
	}
	switch args[0] {
	case "sync-machine-log":
		if len(args) < 2 {
			fmt.Println("please provide sync length")
			return
		}
		conf.MachineLogSyncLength, _ = strconv.Atoi(args[1])
		if err, _ := op.SyncMachineLog(); err != nil {
			panic(err)
		}
		fmt.Println("sync machine log at", util.TimeNow())
	case "test-remote-connection":
		switch {
		case len(args) < 2:
			fmt.Println("please provide machine ID")
			return
		case len(args) < 3:
			fmt.Println("please provide sync length")
			return
		}
		pasID, _ := strconv.Atoi(args[1])
		conf.MachineLogSyncLength, _ = strconv.Atoi(args[2])
		if err := op.TestRemoteConnection(pasID); err != nil {
			panic(err)
		}
	case "reset-zero-workperiods":
		if err := op.ResetInvalidWorkPeriods(); err != nil {
			panic(err)
		}
		fmt.Println("reset zero work periods")
	case "reset-all-workperiods":
		if err := op.ResetAllWorkPeriods(); err != nil {
			panic(err)
		}
		fmt.Println("reset all work periods")
	}
}

func init() {
	RootCmd.AddCommand(TestCmd)
}
