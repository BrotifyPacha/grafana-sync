package fs

import (
	"errors"
	"testing"

	"github.com/brotifypacha/grafana-sync/internal/domain"
	"github.com/brotifypacha/grafana-sync/internal/fs/writer"
	"github.com/brotifypacha/grafana-sync/internal/grafana"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeepFS_Save(t *testing.T) {
	type args struct {
		grafanaFolder domain.GrafanaFolder
		localFolder   string
	}
	tests := []struct {
		name              string
		args              args
		expectedFilepaths []string
	}{
		{
			name: "1",
			args: args{
				grafanaFolder: domain.GrafanaFolder{
					Id:       0,
					Title:    "root",
					FolderId: -1,
					FolderItems: []*domain.GrafanaFolder{
						&domain.GrafanaFolder{
							Title: "PM Metrics",
						},
					},
					DashboardItems: []*domain.GrafanaDashboard{
						&domain.GrafanaDashboard{
							Uid:   "1",
							Title: "dashboard #1",
						},
						&domain.GrafanaDashboard{
							Uid:   "2",
							Title: "dashboard #2",
						},
						&domain.GrafanaDashboard{
							Uid:   "should-error",
							Title: "dashboard #3",
						},
					},
				},
				localFolder: "",
			},
			expectedFilepaths: []string{
				"root",
				"root/PM Metrics",
				"root/dashboard #1",
				"root/dashboard #1/dashboard-data.json",
				"root/dashboard #2",
				"root/dashboard #2/dashboard-data.json",
				"root/dashboard #3",
			},
		},
	}

	ctrl := gomock.NewController(t)
	repo := grafana.NewMockRepository(ctrl)
	repo.EXPECT().GetDashboard(gomock.Any()).AnyTimes().DoAndReturn(func(uid string) (*domain.GrafanaDashboardDetails, error) {
		if uid == "should-error" {
			return nil, errors.New("should error")
		}
		return &domain.GrafanaDashboardDetails{
			Title:   "dashboard #" + uid,
			Panels:  []domain.GrafanaPanel{},
			RawData: []byte{},
		}, nil
	})

	fakeWriter := writer.NewFakeWriter()
	fs := NewDeepFolderFs(
		repo,
		fakeWriter,
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := fs.Save(tt.args.grafanaFolder, tt.args.localFolder)
			assert.Subset(t, tt.expectedFilepaths, fakeWriter.WrittenEntities)
			assert.Len(t, errors, 1)
		})
	}
}

func TestDeepFS_saveDashboard(t *testing.T) {
	type args struct {
		dashboard *domain.GrafanaDashboard
		path      string
	}
	tests := []struct {
		name              string
		args              args
		expectedFilepaths []string
	}{
		{
			name: "dashboard with panels",
			args: args{
				dashboard: &domain.GrafanaDashboard{
					Uid:      "aabbcc",
					Title:    "dashboard #1",
					FolderId: 0,
				},
				path: "./root/",
			},
			expectedFilepaths: []string{
				"root",
				"root/dashboard #1",
				"root/dashboard #1/dashboard-data.json",
				"root/dashboard #1/Panel (uid=uid_of_panel)",
				"root/dashboard #1/Panel (uid=uid_of_panel)/A.sql",
			},
		},
	}

	ctrl := gomock.NewController(t)
	repo := grafana.NewMockRepository(ctrl)
	repo.EXPECT().GetDashboard(gomock.Any()).AnyTimes().DoAndReturn(func(uid string) (*domain.GrafanaDashboardDetails, error) {
		return &domain.GrafanaDashboardDetails{
			Title: "dashboard #1",
			Panels: []domain.GrafanaPanel{
				{
					Uid:   "uid_of_panel",
					Title: "Panel",
					Targets: []domain.Target{
						{
							RefId:    "A",
							Expr:     "",
							RawSql:   "select * from table",
							RawQuery: true,
						},
					},
				},
			},
			RawData: []byte("{}"),
		}, nil
	})

	fakeWriter := writer.NewFakeWriter()
	fs := NewDeepFolderFs(repo, fakeWriter).(*DeepFS)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs.saveDashboard(tt.args.dashboard, tt.args.path)
			assert.Subset(t, tt.expectedFilepaths, fakeWriter.WrittenEntities)
		})
	}
}
