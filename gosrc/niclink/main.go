package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/cloudflare/cloudflared/cmd/cloudflared/cliutil"
	"github.com/cloudflare/cloudflared/cmd/cloudflared/tunnel"
)

var (
	Version   = "2026.9.1"
	BuildTime = "unknown"
	BuildType = ""

	buildInfo *cliutil.BuildInfo
)

func main() {
	os.Setenv("QUIC_GO_DISABLE_ECN", "1")

	buildInfo = cliutil.GetBuildInfo(BuildType, Version)

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v", "V"},
		Usage:   "Print the version",
	}

	app := &cli.App{
		Name:            "niclink",
		Usage:           "Minimal link client",
		UsageText:       "niclink [global options] command [command options]",
		Version:         fmt.Sprintf("%s (built %s)", Version, BuildTime),
		Flags:           linkFlags(),
		HideHelpCommand: true,
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Run a link using a token or credentials",
				Action: cliutil.ConfiguredAction(runAction),
				Flags:  linkFlags(),
			},
		},
	}

	// No token/url-specified invocation runs quick mode through the same action.
	app.Action = cliutil.ConfiguredAction(runAction)

	tunnel.Init(buildInfo, make(chan struct{}))
	runApp(app)
}

func runApp(app *cli.App) {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s\n", c.App.Version)
	}
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	time.Sleep(0)
}
