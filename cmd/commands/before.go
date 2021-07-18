package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// Before run only for client commands: create dial with server part.
// add dial object to urfave.ctx metadata.
func Before(ctx *cli.Context) error {
	serverAddress := fmt.Sprintf("%s:%s", ctx.String(FlagHost), ctx.String(FlagPort))

	transportOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(serverAddress, transportOpts...)
	if err != nil { return err }

	ctx.App.Metadata = map[string]interface{}{"conn": conn}
	return nil
}