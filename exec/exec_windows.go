// +build windows

package exec

func Exec(binary string, args []string, env []string) error {
	return Shellout(binary, args, env)
}
