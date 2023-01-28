package fs

import (
	"io"
	"os"
)

type LocalWriter struct {}

func NewLocalWriter() *LocalWriter {
	return new(LocalWriter)
}

func (*LocalWriter) CreateDir(path string) error {
	return os.Mkdir(path, 0777)
}

func (*LocalWriter) CreateFile(path string) (io.Writer, error) {
	return os.Create(path)
}
