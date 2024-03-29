package writer

import (
	"os"
)

type LocalWriter struct{}

func NewLocalWriter() *LocalWriter {
	return new(LocalWriter)
}

func (*LocalWriter) CreateDir(path string) error {
	err := os.Mkdir(path, 0777)
	if os.IsExist(err) {
		return nil
	}
	return err
}

func (*LocalWriter) CreateFile(path string, content []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	return err
}
