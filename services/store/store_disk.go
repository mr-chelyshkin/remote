package store

import (
	"bytes"
	"io"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

type DiskStore struct {
	directory string
	filename  string
	mode      string

	osFile *os.File
	mutex  sync.RWMutex
}

// NewDiskStore return DiskStore object as pointer.
func NewDiskStore(dir, name, mode string) (*DiskStore, error) {
	file, err := os.Create(path.Join(dir, name))
	if err != nil {
		return nil, errors.Errorf(
			"[DiskStore] error while create file: %s, got %s",
			path.Join(dir, name),
			err.Error(),
		)
	}

	return &DiskStore{
		directory: dir,
		filename:  name,
		osFile:    file,

		mode: mode,
	}, nil
}

// Implement "Store" interface methods

// Write method write content to 'DiskStore.osFile' object.
// Use io.Copy from income data chunks []byte as io.Reader to os.File (DiskStore.osFile) as io.Writer.
// This method allows you NOT to occupy memory with the transferred chunks before recording on disk.
func (ds *DiskStore) Write(chunk []byte) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	_, err := io.Copy(ds.osFile, bytes.NewReader(chunk))
	if err != nil {
		_ = ds.osFile.Close()

		return errors.Errorf(
			"[DiskStore] error while write data to file: %s, got %s",
			path.Join(ds.directory, ds.filename),
			err,
		)
	}

	return nil
}

// Close - changes the mode of the file to mode from DiskStore.mode (convert to os.Mode(uint)) and
// close rendering os.File (DiskStore.osFile) object.
func (ds *DiskStore) Close() error {
	defer ds.osFile.Close()

	uMode, err := strconv.ParseUint(ds.mode, 8, 32)
	if err != nil {
		return errors.Errorf(
			"[DiskStore] error while parse file mode: %s, got %s",
			path.Join(ds.directory, ds.filename),
			err,
		)
	}
	if err := ds.osFile.Chmod(os.FileMode(uMode)); err != nil {
		return errors.Errorf(
			"[DiskStore] error while set file permissions: %s, got %s",
			path.Join(ds.directory, ds.filename),
			err,
		)
	}

	return nil
}