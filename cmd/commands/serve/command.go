package serve

import (
	"fmt"
	"net"
	"path"

	"github.com/mr-chelyshkin/remote/cmd/commands"
	"github.com/mr-chelyshkin/remote/pb"
	"github.com/mr-chelyshkin/remote/server"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

/*
	Command description.
		run application as daemon (server) part.
		listen and execute income requests.
*/

var (
	usage = "run as daemon"
	name  = "serve"
)

func Init(globalFlags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  name,
		Usage: usage,
		Flags: append(globalFlags, commandFlags()...),

		// actions
		Action: func(ctx *cli.Context) error { return action(ctx) },

		// add specific
		Hidden: true,
	}
}

// -- >
func action(ctx *cli.Context) error {
	if err := serveConfig(ctx.String(commands.FlagConfig)); err != nil {
		return err
	}

	// ->
	transportServer := server.NewTransportServer()

	address := fmt.Sprintf("0.0.0.0:%v", ctx.String(commands.FlagPort))
	listener, err := net.Listen("tcp", address)
	if err != nil { return err }

	serverOpts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(serverOpts...)

	pb.RegisterTransportServer(grpcServer, transportServer)
	return grpcServer.Serve(listener)
}

// -- >

func serveConfig(dir string) error {
	if dir != "" {
		configPath, configName := path.Split(dir)

		viper.SetConfigName(configName)
		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")

		return viper.ReadInConfig()
	} else {

		// default values

		return nil
	}
}