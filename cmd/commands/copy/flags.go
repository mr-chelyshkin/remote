package copy

import (
	"github.com/mr-chelyshkin/remote/cmd/commands"

	"github.com/urfave/cli/v2"
)

func commandFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     commands.FlagFrom,
			Required: true,
			Usage:    "specify file or directory",
		},
		&cli.StringFlag{
			Name:     commands.FlagTo,
			Required: true,
			Usage:    "path where to copy",
		},
		&cli.StringFlag{
			Name:     commands.FlagHost,
			Required: true,
			Usage:    "specify target hostname",
		},
		&cli.StringFlag{
			Name:     commands.FlagPort,
			Usage:    "specify port on server",
			Value:    "50053",
		},
	}
}