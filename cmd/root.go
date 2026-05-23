package cmd

import (
	"fmt"
	"os"
)

var Port = "2233"

func Execute() {
	if err := Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run(args []string) error {
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "start", "server":
			return start()
		case "port":
			if index+1 >= len(args) {
				return fmt.Errorf("missing port value")
			}
			Port = args[index+1]
			index++
		case "admin":
			return admin(args[index+1:])
		case "repwd", "clear", "logout", "clear_ip_restriction":
			return admin(args[index:])
		default:
			return fmt.Errorf("unknown command %q", args[index])
		}
	}
	return start()
}
