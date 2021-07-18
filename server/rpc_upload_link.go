package server

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"sync"

	"github.com/mr-chelyshkin/remote/pb"
)

/*
	Server gPRC method.
		Create symlink from client request data.
		Send response message.
*/
func (ts *TransportServer) UploadLink(_ context.Context, data *pb.TransportData) (*pb.Response, error) {
	mu := sync.Mutex{}

	if err := mutexUploadLink(data.LinkPath, data.RemotePath, &mu); err != nil {
		msg := errors.Errorf(
			"[TransportServer] err while exec 'UploadLink', got %s",
			err,
		)

		return &pb.Response{
			Message: msg.Error(),
			Ok:      false,
		}, msg
	}
	return &pb.Response{
		Message: fmt.Sprintf("Symlink: %s -> OK", data.RemotePath),
		Ok:      true,
	}, nil
}

// -- >

func mutexUploadLink(oldName, newName string, mu *sync.Mutex) error {
	mu.Lock()
	defer mu.Unlock()

	os.Remove(newName)
	return os.Symlink(oldName, newName)
}