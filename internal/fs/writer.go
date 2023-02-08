package fs

type Writer interface {
	CreateDir(path string) error
	CreateFile(path string, content []byte) error
}
