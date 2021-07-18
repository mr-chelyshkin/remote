package serve

import (
	"github.com/mr-chelyshkin/remote/cmd/commands"

	"github.com/urfave/cli/v2"
)

func commandFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  commands.FlagPort,
			Usage: "specify port for listening",
			Value: "50053",
		},
		&cli.StringFlag{
			Name:  commands.FlagConfig,
			Usage: "path to serve configs",
		},
	}
}