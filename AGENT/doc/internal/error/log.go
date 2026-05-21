package error

import (
	"fmt"
	"strings"
	"zhongxin/util"
)

func Print(err error) {
	if err != nil {
		util.Log.Errorln(err)
		stackTrace := fmt.Sprintf("%+v", err)
		lines := strings.Split(stackTrace, "\n")
		for _, line := range lines {
			if strings.Contains(line, "zhongxin") {
				util.Log.Errorln(line)
			}
		}
		fmt.Println()
	}
}
