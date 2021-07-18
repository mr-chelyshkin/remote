package representer

import (
	"github.com/mr-chelyshkin/remote/pb"
)

type DataDir struct {
	pathFrom  string
	pathTo    string
}

// newDataDir create DataDir object as pointer.
func newDataDir(pathFrom, pathTo string) *DataDir {
	return &DataDir{
		pathFrom:  pathFrom,
		pathTo:    pathTo,
	}
}

// Implement "Transporter" interface methods

func (dd DataDir) LocalPath() string {
	return dd.pathFrom
}

func (dd DataDir) Transport() *pb.TransportData {
	return &pb.TransportData{
		RemotePath: dd.pathTo,
	}
}