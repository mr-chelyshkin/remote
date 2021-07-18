package copy

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"sync"

	"github.com/mr-chelyshkin/remote/cmd/commands"
	"github.com/mr-chelyshkin/remote/pb"
	"github.com/mr-chelyshkin/remote/requests"
	"github.com/mr-chelyshkin/remote/services/representer"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

/*
	Command description.
		Make copy data from local directory to remote.
*/

var (
	usage = "copying data to another server"
	name  = "cp"
)

func Init(globalFlags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  name,
		Usage: usage,
		Flags: append(globalFlags, commandFlags()...),

		// actions
		Before: func(ctx *cli.Context) error { return commands.Before(ctx) },
		Action: func(ctx *cli.Context) error { return action(ctx) },
		After:  func(ctx *cli.Context) error { return commands.After(ctx) },
	}
}

// -- >
func action(ctx *cli.Context) error {
	filesList, err := representer.TransportFiles(ctx.String(commands.FlagFrom), ctx.String(commands.FlagTo))
	if err != nil { return err }

	linksList, err := representer.TransportSymLinks(ctx.String(commands.FlagFrom), ctx.String(commands.FlagTo))
	if err != nil { return err }

	dirsList, err := representer.TransportDir(ctx.String(commands.FlagFrom), ctx.String(commands.FlagTo))
	if err != nil { return err }

	// ->
	conn := ctx.App.Metadata["conn"].(*grpc.ClientConn)
	client := pb.NewTransportClient(conn)

	if err := requests.AsyncWithLimit(dirsList, client, ctx.Int(commands.FlagLimit), execActionDir); err != nil {
		return err
	}
	if err := requests.AsyncWithLimit(filesList, client, ctx.Int(commands.FlagLimit), execActionFile); err != nil {
		return err
	}
	if err := requests.AsyncWithLimit(linksList, client, ctx.Int(commands.FlagLimit), execActionLink); err != nil {
		return err
	}
	return nil
}

// -- >

var execActionFile = func (data representer.Transporter, client pb.TransportClient) error {
	mutex := sync.Mutex{}
	defer mutex.Unlock()
	mutex.Lock()

	file, err := os.Open(data.LocalPath())
	if err != nil {
		return errors.Errorf(
			"[ExecAction] err while open %s, got %s",
			data.LocalPath(),
			err,
		)
	}
	defer file.Close()

	// ->
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.UploadFile(ctx)
	if err != nil { return err }

	if err := stream.Send(data.Transport()); err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF { break }
		if err != nil { return err }

		if err := stream.Send(&pb.TransportData{Chunk: buffer[:n]}); err != nil {
			return err
		}
	}
	// ->

	response, err := stream.CloseAndRecv()
	if err != nil { return err }

	if !response.Ok {
		return errors.Errorf(response.Message)
	}

	log.Println(response.Message)
	return nil
}

var execActionLink = func (data representer.Transporter, client pb.TransportClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.UploadLink(ctx, data.Transport())
	if err != nil { return err }

	if !response.Ok {
		return errors.Errorf(response.Message)
	}

	log.Println(response.Message)
	return nil
}

var execActionDir = func (data representer.Transporter, client pb.TransportClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.UploadDir(ctx, data.Transport())
	if err != nil { return err }

	if !response.Ok {
		return errors.Errorf(response.Message)
	}

	log.Println(response.Message)
	return nil
}