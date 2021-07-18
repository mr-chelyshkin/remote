package commands

import (
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// After run only for CLIENT commands: close dial with server part.
func After(ctx *cli.Context) error {
	return ctx.App.Metadata["conn"].(*grpc.ClientConn).Close()
}