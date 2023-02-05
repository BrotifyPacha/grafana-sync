package fs

type FakeWriter struct {
	WrittenEntities []string
}

func NewFakeWriter() *FakeWriter {
	return new(FakeWriter)
}

func (f *FakeWriter) CreateDir(path string) error {
	f.WrittenEntities = append(f.WrittenEntities, path)
	return nil
}

func (f *FakeWriter) CreateFile(path string, content []byte) error {
	f.WrittenEntities = append(f.WrittenEntities, path)
	return nil
}
