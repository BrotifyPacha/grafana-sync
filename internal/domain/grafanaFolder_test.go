package domain

import (
	"reflect"
	"testing"
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
					Uid:            "uid-25",
					Id:             25,
					Title:          "Secondary sub 1",
					FolderItems:    []*GrafanaFolder{},
					DashboardItems: []*GrafanaDashboard{},
				},
				&GrafanaFolder{
					Uid:            "uid-30",
					Id:             30,
					Title:          "Secondary sub 2",
					FolderItems:    []*GrafanaFolder{},
					DashboardItems: []*GrafanaDashboard{},
				},
			},
			DashboardItems: []*GrafanaDashboard{},
		},
	},
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
