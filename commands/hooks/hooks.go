package hooks

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/docker/cli/cli-plugins/hooks"
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli/command"

	"github.com/spf13/cobra"
)

func HookCommand(dockerCli command.Cli) *cobra.Command {
	hookCmd := &cobra.Command{
		Use: manager.HookSubcommandName,
		RunE: func(_ *cobra.Command, args []string) error {
			runHooks(dockerCli.Out(), []byte(args[0]))
			return nil
		},
		Args:   cobra.ExactArgs(1),
		Hidden: true,
	}

	return hookCmd
}

func runHooks(out io.Writer, input []byte) {
	var c manager.HookPluginData
	err := json.Unmarshal(input, &c)
	if err != nil {
		return
	}

	hint, shouldPrint := getHint(c.RootCmd, c.Flags, runtime.NumCPU())
	if !shouldPrint {
		return
	}

	returnType := hooks.HookMessage{
		Type:     hooks.NextSteps,
		Template: hint,
	}
	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "     ")
	_ = enc.Encode(returnType)
}

var (
	multiPlatformHint = fmt.Sprintf("Build multi-platform images faster with Docker Build Cloud: %s", dbcLink())
	reduceTimeHint    = fmt.Sprintf("Reduce build time with Docker Build Cloud: %s", dbcLink())
)

func getHint(cmd string, flags map[string]string, numCPUs int) (string, bool) {
	if cmd != "buildx build" {
		return "", false
	}

	if _, progress := flags["progress"]; progress {
		return "", false
	}
	if _, ok := os.LookupEnv("BUILDKIT_PROGRESS"); ok {
		return "", false
	}

	if _, platformFlag := flags["platform"]; platformFlag {
		return multiPlatformHint, true
	} else if numCPUs < 4 {
		return reduceTimeHint, true
	}
	return "", false
}

func dbcLink() string {
	link := "https://docs.docker.com/go/docker-build-cloud"
	// create an escape sequence using the OSC 8 format: https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, link)
}
