package grafana

import "github.com/brotifypacha/grafana_searcher/internal/domain"

//go:generate mockgen -source repository.go -package grafana -destination mockRepository.go
type Repository interface {
	GetTree() (*domain.GrafanaFolder, error)
	GetDashboard(uid string) ([]byte, error)
}
