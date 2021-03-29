package e2e

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/docker/compose-switch/redirect"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/icmd"
)

// E2eCLI is used to wrap the CLI for end to end testing
// nolint stutter
type E2eCLI struct {
	BinDir    string
	ConfigDir string
	Test      *testing.T
}

func getBinExt() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

var dockerCompose = "docker-compose" + getBinExt()

// NewE2eCLI returns a configured TestE2eCLI
func NewE2eCLI(t *testing.T, binDir string) *E2eCLI {
	configDir, err := ioutil.TempDir("", "")
	assert.Check(t, is.Nil(err))

	t.Parallel()

	t.Cleanup(func() {
		if t.Failed() {
			conf, _ := ioutil.ReadFile(filepath.Join(configDir, "features.json"))
			t.Errorf("features.json: %s\n", string(conf))
		}
		_ = os.RemoveAll(configDir)
	})

	return &E2eCLI{binDir, configDir, t}
}

// SetupExistingCLI copies the existing CLI in a temporary directory so that the
// new CLI can be configured to use it
func SetupExistingCLI() (string, func(), error) {
	binDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", nil, err
	}

	bin, err := findExecutable(dockerCompose, []string{"../../bin"})
	if err != nil {
		return "", nil, err
	}

	if err := CopyFile(bin, filepath.Join(binDir, dockerCompose)); err != nil {
		return "", nil, err
	}

	echostub, err := findExecutable("echostub"+getBinExt(), []string{"../../bin/test"})
	if err != nil {
		return "", nil, err
	}

	if err := CopyFile(echostub, filepath.Join(binDir, redirect.ComposeV1Binary)); err != nil {
		return "", nil, err
	}

	if err := CopyFile(echostub, filepath.Join(binDir, redirect.DockerBinary)); err != nil {
		return "", nil, err
	}

	cleanup := func() {
		_ = os.RemoveAll(binDir)
	}

	return binDir, cleanup, nil
}

func findExecutable(executableName string, paths []string) (string, error) {
	for _, p := range paths {
		bin, err := filepath.Abs(path.Join(p, executableName))
		if err != nil {
			return "", err
		}

		if _, err := os.Stat(bin); os.IsNotExist(err) {
			continue
		}

		return bin, nil
	}

	return "", fmt.Errorf("executable not found %q in %s", executableName, strings.Join(paths, ", "))
}

// CopyFile copies a file from a sourceFile to a destinationFile setting permissions to 0755
func CopyFile(sourceFile string, destinationFile string) error {
	src, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	// nolint: errcheck
	defer src.Close()

	dst, err := os.OpenFile(destinationFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	// nolint: errcheck
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return err
}

// NewCmd creates a cmd object configured with the test environment set
func (c *E2eCLI) NewCmd(command string, args ...string) icmd.Cmd {
	testPath := c.PathEnvVar()
	env := append(os.Environ(),
		"DOCKER_CONFIG="+c.ConfigDir,
		"TEST_DESKTOP_SOCKET="+c.DesktopSocket(),
		"PATH="+testPath,
	)
	return icmd.Cmd{
		Command: append([]string{command}, args...),
		Env:     env,
	}
}

// RunNewComposeCmd runs a docker-compose cmd
func (c *E2eCLI) RunNewComposeCmd(args ...string) *icmd.Result {
	fmt.Printf("	[%s] %s %s\n", c.Test.Name(), filepath.Join(c.BinDir, dockerCompose), strings.Join(args, " "))
	return icmd.RunCmd(c.NewComposeCmd(args...))
}

// NewComposeCmd creates a docker-compose cmd without running it
func (c *E2eCLI) NewComposeCmd(args ...string) icmd.Cmd {
	return c.NewCmd(filepath.Join(c.BinDir, dockerCompose), args...)
}

// RunCmd runs a command, expects no error and returns a result
func (c *E2eCLI) RunCmd(args ...string) *icmd.Result {
	fmt.Printf("	[%s] %s\n", c.Test.Name(), strings.Join(args, " "))
	assert.Assert(c.Test, len(args) >= 1, "require at least one command in parameters")
	res := icmd.RunCmd(c.NewCmd(args[0], args[1:]...))
	res.Assert(c.Test, icmd.Success)
	return res
}

// DesktopSocket get the path where test metrics will be sent
func (c *E2eCLI) DesktopSocket() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\` + filepath.Base(c.ConfigDir) + "_metrics"
	}
	return filepath.Join(c.ConfigDir, "./docker-cli.sock")
}

// PathEnvVar returns path (os sensitive) for running test
func (c *E2eCLI) PathEnvVar() string {
	path := c.BinDir + ":" + os.Getenv("PATH")
	if runtime.GOOS == "windows" {
		path = c.BinDir + ";" + os.Getenv("PATH")
	}
	return path
}
