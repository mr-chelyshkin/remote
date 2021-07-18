package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mr-chelyshkin/remote/pb"
	"github.com/mr-chelyshkin/remote/services/store"
	"github.com/pkg/errors"
	"io"
	"path"
)

/*
	Server gPRC method.
		Create file object and write to it from client stream connection by chunks.
		Call DiskStore object for save data.
		Send response and close connection.
*/
func (ts *TransportServer) UploadFile(stream pb.Transport_UploadFileServer) error {
	fileBody := bytes.Buffer{}
	defer fileBody.Reset()

	fileMeta, err := stream.Recv()
	if err != nil {
		return errors.Errorf(
			"[TransportServer] error while stream 'UploadFile', got %s",
			err,
		)
	}

	diskObject, err := store.NewDiskStore(fileMeta.RemotePath, fileMeta.FileName, fileMeta.Mode)
	if err != nil {
		return err
	}
	defer diskObject.Close()

	// ->
	for {
		switch stream.Context().Err() {
		case context.Canceled:
			if err := closeStream(stream, false, "network: request is canceled"); err != nil {
				return errors.Errorf(
					"[TransportServer] error while close 'UploadFile' stream, got %s",
					err,
				)
			}
			return fmt.Errorf("network: request is canceled")
		case context.DeadlineExceeded:
			if err := closeStream(stream, false, "network: deadline is exceeded"); err != nil {
				return errors.Errorf(
					"[TransportServer] error while close 'UploadFile' stream, got %s",
					err,
				)
			}
			return fmt.Errorf("network: deadline is exceeded")
		default:
		}

		fileData, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			msg := errors.Errorf("[TransportServer] error while stream 'UploadFile', got %s", err)
			if err := closeStream(stream, false, msg.Error()); err != nil {
				return errors.Errorf(
					"[TransportServer] error while close 'UploadFile' stream, got %s",
					err,
				)
			}
			return msg
		}

		if err := diskObject.Write(fileData.GetChunk()); err != nil {
			return err
		}
	}
	// ->

	return closeStream(stream, true, fmt.Sprintf("File: %s -> OK", path.Join(fileMeta.RemotePath, fileMeta.FileName)))
}

// -- >

// closeStream close stream and send response message.
func closeStream(stream pb.Transport_UploadFileServer, status bool, message string) error {
	response := pb.Response{
		Message: message,
		Ok:      status,
	}
	return stream.SendAndClose(&response)
}

//TODO: defer remove files from cmd files requests (use in-memory db)