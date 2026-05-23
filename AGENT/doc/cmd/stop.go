package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"syscall"
	"zhongxin/util"
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop zhongxin server by daemon/pid file",
	Run: func(cmd *cobra.Command, args []string) {
		stop()
	},
}

func stop() {
	initDaemon()
	if pid == -1 {
		util.Log.Info("Seems not have been started. Try use `zhongxin start` to start server.")
		return
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		util.Log.Errorf("failed to find process by pid: %d, reason: %v", pid, process)
		return
	}
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		util.Log.Errorf("failed to terminate process %d: %v", pid, err)
	} else {
		util.Log.Info("terminated process: ", pid)
	}
	err = os.Remove(pidFile)
	if err != nil {
		util.Log.Errorf("failed to remove pid file")
	}
	pid = -1
}

func init() {
	RootCmd.AddCommand(StopCmd)
}
