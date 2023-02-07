package fs

import (
	"testing"

	"github.com/brotifypacha/grafana_searcher/internal/domain"
	"github.com/brotifypacha/grafana_searcher/internal/grafana"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type SaveTest struct {
	Folder        *domain.GrafanaFolder
	SavedEntities []string
}

func TestSingleFolderFS_Save(t *testing.T) {
	ctrl := gomock.NewController(t)

	dataTable := []SaveTest{
		{
			Folder: &domain.GrafanaFolder{
				Id:       0,
				Title:    "FolderName",
				FolderId: -1,
				FolderItems: []*domain.GrafanaFolder{
					{
						Id:       1,
						Title:    "Inner",
						FolderId: 0,
						DashboardItems: []*domain.GrafanaDashboard{
							{
								Title:    "dashboard_2",
								FolderId: 1,
							},
						},
					},
				},
				DashboardItems: []*domain.GrafanaDashboard{
					{
						Title:    "dashboard_1",
						FolderId: 0,
					},
				},
			},
			SavedEntities: []string{
				"FolderName",
				"FolderName/dashboard_1.json",
				"FolderName/Inner__dashboard_2.json",
			},
		},
		{
			Folder: &domain.GrafanaFolder{
				Id:          0,
				Title:       "FolderName",
				FolderId:    -1,
				FolderItems: []*domain.GrafanaFolder{},
				DashboardItems: []*domain.GrafanaDashboard{
					{
						Title:    "dashboard_1",
						FolderId: 0,
					},
					{
						Title:    "dashboard_2",
						FolderId: 0,
					},
				},
			},
			SavedEntities: []string{
				"FolderName",
				"FolderName/dashboard_1.json",
				"FolderName/dashboard_2.json",
			},
		},
	}

	repo := grafana.NewMockRepository(ctrl)
	repo.EXPECT().GetDashboard(gomock.Any()).AnyTimes()
	writer := NewFakeWriter()

	fileSystem := NewSingleFolderFS(repo, writer)

	for _, test := range dataTable {
		fileSystem.Save(test.Folder, ".")
		assert.Subset(t, test.SavedEntities, writer.WrittenEntities)
		writer.WrittenEntities = make([]string, 0)
	}

}
