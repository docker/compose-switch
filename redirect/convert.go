package redirect

import (
	"fmt"
	"os"
)

var (
	boolflags = []string{
		"--debug", "-D",
		"--verbose",
		"--tls",
		"--tlsverify",
	}

	stringflags = []string{
		"--tlscacert",
		"--tlscert",
		"--tlskey",
		"--host", "-H",
		"--log-level", "-l",
		"--context",
	}
)

func Convert(args []string) []string {
	var rootFlags []string
	command := []string{"compose", "--compatibility"}
	l := len(args)
	for i := 0; i < l; i++ {
		arg := args[i]
		if arg[0] != '-' {
			// not a top-level flag anymore, keep the rest of the command unmodified
			command = append(command, args[i:]...)
			break
		}
		if arg == "--verbose" {
			arg = "--debug"
		}
		if arg == "-h" {
			// docker cli has deprecated -h to avoid ambiguity with -H, while docker-compose still support it
			arg = "--help"
		}
		if arg == "--version" || arg == "-v" {
			// redirect --version pseudo-command to actual command
			arg = "version"
		}
		if contains(boolflags, arg) {
			rootFlags = append(rootFlags, arg)
			continue
		}
		if contains(stringflags, arg) {
			i++
			if i >= l {
				fmt.Fprintf(os.Stderr, "flag needs an argument: '%s'\n", arg)
				os.Exit(1)
			}
			rootFlags = append(rootFlags, arg, args[i])
			continue
		}
		command = append(command, arg)
	}
	return append(rootFlags, command...)
}

func contains(array []string, needle string) bool {
	for _, val := range array {
		if val == needle {
			return true
		}
	}
	return false
}
