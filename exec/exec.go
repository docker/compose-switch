// +build !windows

package exec

import "syscall"

func Exec(binary string, args []string, env []string) error {
	return syscall.Exec(binary, args, env)
}
