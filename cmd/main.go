package main

import (
	"log"
	"os"

	consts "github.com/mr-chelyshkin/remote/cmd/commands"
	"github.com/mr-chelyshkin/remote/cmd/commands/copy"
	"github.com/mr-chelyshkin/remote/cmd/commands/serve"

	"github.com/urfave/cli/v2"
)

var (
	usageText = "usageText"
	version   = "Version"
	usage     = "Usage"
)

// before is pre-define variable as function for execute "before command" actions
var before = func(ctx *cli.Context) error {
	return nil
}

// after is pre-define variable as function for execute "after commend" actions
var after = func(ctx *cli.Context) error {
	return nil
}

// flags returns list of flags for all app commands
func flags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:  consts.FlagLimit,
			Usage: "goroutine limits",
			Value: 100,
		},
	}
}

// commands returns app command list
func commands(flags []cli.Flag) []*cli.Command {
	return []*cli.Command {
		serve.Init(flags),
		copy.Init(flags),
	}
}

// -- >
func main() {
	app := cli.NewApp()

	app.UsageText = usageText
	app.Version   = version
	app.Usage     = usage

	app.Before    = before
	app.After     = after

	app.Flags     = flags()
	app.Commands  = commands(app.Flags)

	// -- >
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}