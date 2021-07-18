package requests

import (
	"context"
	"github.com/gammazero/workerpool"
	"github.com/mr-chelyshkin/remote/pb"
	"github.com/mr-chelyshkin/remote/services/representer"
)

// AsyncWithLimit is decorator for command job functions (async execution with goroutines limits).
// break executing if error.
// use: https://github.com/gammazero/workerpool
func AsyncWithLimit(transportDataList *[]representer.Transporter, client pb.TransportClient, limitCount int,
	job func(data representer.Transporter, client pb.TransportClient) error,
) error {
	// -- >

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp := workerpool.New(limitCount)
	ch := make(chan error, len(*transportDataList))

	for _, data := range *transportDataList {
		_data := data

		wp.Submit(func() {
			if err := job(_data, client); err != nil {
				ch <- err
				cancel()
			}
		})
	}

	ctx.Done()
	wp.StopWait()
	close(ch)

	for r := range ch {
		if r != nil {
			return r
		}
	}
	return nil
}