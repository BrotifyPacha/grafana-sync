package fs

import (
	"github.com/brotifypacha/grafana-sync/internal/domain"
)

type FileSystemInterface interface {
	Save(grafanaFolder domain.GrafanaFolder, localFolder string) []error
}
