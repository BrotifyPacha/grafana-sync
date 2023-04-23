package grafana

import (
	"testing"

	"github.com/brotifypacha/grafana-sync/internal/domain"
	"github.com/stretchr/testify/assert"
)

func Test_buildTree(t *testing.T) {
	type args struct {
		grafanaItems []*RawGrafanaApiItem
	}
	tests := []struct {
		name string
		args args
		want *domain.GrafanaFolder
	}{
		{
			name: "",
			args: args{
				grafanaItems: []*RawGrafanaApiItem{
					{
						Id:       1,
						Type:     ITEM_TYPE_FOLDER,
						Title:    "main folder",
						FolderId: 0,
					},
					{
						Id:       2,
						Type:     ITEM_TYPE_DASHBORAD,
						Uid:      "dashboard1",
						Title:    "dashboard",
						FolderId: 1,
					},
					{
						// dashboard is placed in unexisting folder
						// and so it should not be linked to existing structure
						Id:       3,
						Type:     ITEM_TYPE_DASHBORAD,
						Title:    "dashboard",
						FolderId: -100,
					},
				},
			},
			want: &domain.GrafanaFolder{
				Title:    "Root folder",
				Uid:      domain.RootFolderUid,
				Id:       domain.RootFolderId,
				FolderId: -1,
				FolderItems: []*domain.GrafanaFolder{
					{
						Id:          1,
						Title:       "main folder",
						FolderId:    domain.RootFolderId,
						FolderItems: make([]*domain.GrafanaFolder, 0),
						DashboardItems: []*domain.GrafanaDashboard{
							{
								Uid:      "dashboard1",
								Title:    "dashboard",
								FolderId: 1,
							},
						},
					},
				},
				DashboardItems: make([]*domain.GrafanaDashboard, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildTree(tt.args.grafanaItems); !assert.Equal(t, tt.want, got) {
				t.Errorf("buildTree() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
