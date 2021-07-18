package representer

import (
	"github.com/mr-chelyshkin/remote/pb"
)

type SymLink struct {
	oldPath  string
	newPath  string
	linkPath string
	name     string
}

// newSymLink create SymLink object as pointer.
func newSymLink(original, symlink string) *SymLink {
	return &SymLink{
		newPath:  original,
		linkPath: symlink,
	}
}

// Implement "Transporter" interface methods

func (sl SymLink) LocalPath() string {
	return sl.oldPath
}

func (sl SymLink) Transport() *pb.TransportData {
	return &pb.TransportData{
		RemotePath: sl.newPath,
		LinkPath:   sl.linkPath,
	}
}