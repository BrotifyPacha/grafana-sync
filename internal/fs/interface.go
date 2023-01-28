package fs

import (
	"io"
	"github.com/brotifypacha/grafana_searcher/internal/domain"
)

type FileSystemInterface interface {
	Save(grafanaFolder domain.GrafanaFolder, localFolder string) error
}

type FileSystemWriter interface {
	CreateDir(path string) error
	CreateFile(path string) (io.Writer, error)
}
