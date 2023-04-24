package domain

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var FolderStructure = &GrafanaFolder{
	Uid:      RootFolderUid,
	Id:       RootFolderId,
	Title:    "root",
	FolderId: -1,
	FolderItems: []*GrafanaFolder{
		&GrafanaFolder{
			Uid:   "uid-10",
			Id:    10,
			Title: "Main",
			FolderItems: []*GrafanaFolder{
				&GrafanaFolder{
					Uid:            "uid-40",
					Id:             40,
					Title:          "Main sub 1",
					FolderItems:    []*GrafanaFolder{},
					DashboardItems: []*GrafanaDashboard{},
				},
			},
			DashboardItems: []*GrafanaDashboard{},
		},
		&GrafanaFolder{
			Uid:   "uid-20",
			Id:    20,
			Title: "Secondary",
			FolderItems: []*GrafanaFolder{
				&GrafanaFolder{
					Uid:         "uid-25",
					Id:          25,
					Title:       "Secondary sub 1",
					FolderItems: []*GrafanaFolder{},
					DashboardItems: []*GrafanaDashboard{
						&GrafanaDashboard{
							Uid:      "uid-1",
							Title:    "Dashboard #1",
							FolderId: 25,
						},
						&GrafanaDashboard{
							Uid:      "uid-2",
							Title:    "Dashboard #2",
							FolderId: 25,
						},
					},
				},
				&GrafanaFolder{
					Uid:            "uid-30",
					Id:             30,
					Title:          "Secondary sub 2",
					FolderItems:    []*GrafanaFolder{},
					DashboardItems: []*GrafanaDashboard{},
				},
			},
		},
	},
	DashboardItems: []*GrafanaDashboard{},
}

func TestGrafanaFolder_FindFolderByUid(t *testing.T) {
	tests := []struct {
		name    string
		folder  *GrafanaFolder
		findUid string
		want    *GrafanaFolder
		wantErr bool
	}{
		{
			name:    "find success 1",
			folder:  FolderStructure,
			findUid: RootFolderUid,
			want:    FolderStructure,
			wantErr: false,
		},
		{
			name:    "find success 2",
			folder:  FolderStructure,
			findUid: "uid-40",
			want:    FolderStructure.FolderItems[0].FolderItems[0],
			wantErr: false,
		},
		{
			name:    "find success 3",
			folder:  FolderStructure,
			findUid: "uid-30",
			want:    FolderStructure.FolderItems[1].FolderItems[1],
			wantErr: false,
		},
		{
			name:    "find failure",
			folder:  FolderStructure,
			findUid: "not-valid-uid",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.folder.FindFolderByUid(tt.findUid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrafanaFolder.FindFolderById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrafanaFolder.FindFolderById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prettySprint(t *testing.T) {
	type args struct {
		folder    *GrafanaFolder
		indent    string
		recursive bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "recursive",
			args: args{
				folder:    FolderStructure,
				indent:    "",
				recursive: true,
			},
			want: `root [0]
╠═ Main [uid-10]
║  ╚═ Main sub 1 [uid-40]
╚═ Secondary [uid-20]
   ╠═ Secondary sub 1 [uid-25]
   ║  ╟─ Dashboard #1 [uid-1]
   ║  ╙─ Dashboard #2 [uid-2]
   ╚═ Secondary sub 2 [uid-30]
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prettySprint(tt.args.folder, tt.args.indent, tt.args.recursive)
			assert.Equal(t, tt.want, got)
		})
	}
}
