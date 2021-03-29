package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/compose-switch/redirect"
	. "github.com/docker/compose-switch/tests"
	"gotest.tools/v3/icmd"
)

var binDir string

const testUUID = `ABCDEFAB-ABCD-1234-5678-ABCDEFABCDEF`

func TestMain(m *testing.M) {
	p, cleanup, err := SetupExistingCLI()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	binDir = p
	exitCode := m.Run()
	cleanup()
	os.Exit(exitCode)
}

func TestComposeRedirect(t *testing.T) {
	c := NewE2eCLI(t, binDir)

	t.Run("disable-v2", func(t *testing.T) {
		c.RunNewComposeCmd("disable-v2")
		c.RunNewComposeCmd("--version").Assert(t, icmd.Expected{Out: redirect.ComposeV1Binary + " --version"})
		c.RunCmd("cat", filepath.Join(c.ConfigDir, "features.json")).Assert(t, icmd.Expected{Out: `"composeV2": "disabled"`})
	})

	t.Run("enable-v2", func(t *testing.T) {
		c.RunNewComposeCmd("enable-v2")
		c.RunNewComposeCmd("--version").Assert(t, icmd.Expected{Out: redirect.DockerBinary + " compose version"})
		c.RunCmd("cat", filepath.Join(c.ConfigDir, "features.json")).Assert(t, icmd.Expected{Out: `"composeV2": "enabled"`})
	})
}
