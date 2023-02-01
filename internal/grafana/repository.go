package grafana

import "github.com/brotifypacha/grafana_searcher/internal/domain"

type Repository interface {
	GetTree() (*domain.GrafanaFolder, error)
	GetDashboard(uid string) ([]byte, error)
}
