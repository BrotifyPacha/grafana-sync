package grafana

import "github.com/brotifypacha/grafana-sync/internal/domain"

var (
	RootRepositoryId = 0
)

//go:generate mockgen -source repository.go -package grafana -destination mockRepository.go
type Repository interface {
	GetTree() (*domain.GrafanaFolder, error)
	GetDashboard(uid string) ([]byte, error)
}
