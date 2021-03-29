package redirect

import (
	"fmt"
	"os"
	"os/exec"

	composeexec "github.com/docker/compose-switch/exec"
)

func RunComposeV1(args []string) {
	execBinary, err := exec.LookPath(ComposeV1Binary)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Current PATH : "+os.Getenv("PATH"))
		os.Exit(1)
	}
	// syscall.Exec for compose v1 will not work if using compose V1 python self-extracting binary, as it belives it's packaged as docker-compose and not docker-compose-v1
	err = composeexec.Shellout(execBinary, append([]string{"docker-compose"}, args...), os.Environ())
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			os.Exit(exiterr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func RunComposeV2(args []string) {
	execBinary, err := exec.LookPath(DockerBinary)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Current PATH : "+os.Getenv("PATH"))
		os.Exit(1)
	}
	err = composeexec.Exec(execBinary, append([]string{"docker"}, args...), append(os.Environ(), "DOCKER_METRICS_SOURCE=cli-compose-v2redirect"))
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			os.Exit(exiterr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
