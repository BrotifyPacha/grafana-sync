package fs

import (
	"github.com/brotifypacha/grafana_searcher/internal/domain"
)

type FileSystemInterface interface {
	Save(grafanaFolder domain.GrafanaFolder, localFolder string) error
}