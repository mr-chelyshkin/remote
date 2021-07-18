package server

import (
	"context"
	"fmt"
	"github.com/mr-chelyshkin/remote/server/hooks"
	"os"
	"sync"

	"github.com/mr-chelyshkin/remote/pb"

	"github.com/pkg/errors"
)

/*
	Server gPRC method.
		Create directory from client request data (data.RemotePath).
		Send response message.
*/
func (ts *TransportServer) UploadDir(_ context.Context, data *pb.TransportData) (*pb.Response, error) {
	mu := sync.Mutex{}

	if err := mutexUploadDir(data.RemotePath, &mu); err != nil {
		msg := errors.Errorf(
			"[TransportServer] err while exec 'uploadDir', got %s",
			err,
		)

		return &pb.Response{
			Message: msg.Error(),
			Ok:      false,
		}, msg
	}
	return &pb.Response{
		Message: fmt.Sprintf("Directory: %s -> OK", data.RemotePath),
		Ok:      true,
	}, nil
}

// -- >

func mutexUploadDir(path string, mu *sync.Mutex) error {
	mu.Lock()
	defer mu.Unlock()

	if err := hooks.ValidatePath(path); err != nil {
		return err
	}
	return os.MkdirAll(path, 0755)
}