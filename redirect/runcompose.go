package redirect

import (
	"fmt"
	"os"
	"os/exec"

	composeexec "github.com/docker/compose-switch/exec"
)

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
