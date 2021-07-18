package representer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mr-chelyshkin/remote/pb"
)

type DataFile struct {
	name      string
	extension string

	pathFrom  string
	pathTo    string

	mode os.FileMode
	size int64
}

// newDataFile create DataFile object as pointer.
func newDataFile(pathFrom, pathTo string, info os.FileInfo) *DataFile {
	return &DataFile{
		name:      info.Name(),
		extension: filepath.Ext(pathFrom),

		pathFrom:  pathFrom,
		pathTo:    pathTo,

		mode: info.Mode(),
		size: info.Size(),
	}
}

// Implement "Transporter" interface methods

func (df DataFile) LocalPath() string {
	return df.pathFrom
}

func (df DataFile) Transport() *pb.TransportData {
	return &pb.TransportData{
		FileName:   df.name,
		RemotePath: df.pathTo,
		Extension:  df.extension,

		Size:       df.size,
		Mode:       fmt.Sprintf("%04o", df.mode),
	}
}