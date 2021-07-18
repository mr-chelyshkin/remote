package representer

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mr-chelyshkin/remote/pb"

	"github.com/pkg/errors"
)

/*
	Transporter service interface.
	Include method Transport() for represent proto request object.
*/
type Transporter interface {
	LocalPath() string
	Transport() *pb.TransportData
}

// TransportFiles recurse get all files as DataFile object pointer list from income path.
func TransportFiles(pathFrom, pathTo string) (*[]Transporter, error) {
	var dataList []Transporter

	_walkActions := func (from string, info os.FileInfo, err error) error {
		if err != nil { return err }

		if (!info.IsDir()) && info.Mode()&os.ModeSymlink == 0 {
			dirFrom, _ := filepath.Split(from)
			to := path.Clean(path.Join(pathTo, strings.Replace(dirFrom, pathFrom, "",1)))

			dataList = append(dataList, newDataFile(from, to, info))
		}
		return nil
	}

	// ->
	if err := filepath.Walk(pathFrom, _walkActions); err != nil {
		return nil, errors.Errorf("Recurse file walk error, got %s", err)
	}
	return &dataList, nil
}

// TransportSymLinks recurse get all symlinks as SymLink object pointer list from income path.
func TransportSymLinks(pathFrom, pathTo string) (*[]Transporter, error) {
	var dataList []Transporter

	_walkActions := func (from string, info os.FileInfo, err error) error {
		if err != nil { return err }

		if (!info.IsDir()) && info.Mode()&os.ModeSymlink != 0 {
			to := path.Clean(path.Join(pathTo, strings.Replace(from, pathFrom, "",1)))
			link, _ := os.Readlink(from)

			dataList = append(dataList, newSymLink(to, link))
		}
		return nil
	}

	// ->
	if err := filepath.Walk(pathFrom, _walkActions); err != nil {
		return nil, errors.Errorf("Recurse symlink walk error, got %s", err)
	}
	return &dataList, nil
}

// TransportDir recurse get all symlinks as DataDir object pointer list from income path.
func TransportDir(pathFrom, pathTo string) (*[]Transporter, error) {
	var dataList []Transporter

	_walkActions := func (from string, info os.FileInfo, err error) error {
		if err != nil { return err }

		if info.IsDir() {
			to := path.Clean(path.Join(pathTo, strings.Replace(from, pathFrom, "",1)))

			dataList = append(dataList, newDataDir(from, to))
		}
		return nil
	}

	// ->
	if err := filepath.Walk(pathFrom, _walkActions); err != nil {
		return nil, errors.Errorf("Recurse file walk error, got %s", err)
	}
	return &dataList, nil
}