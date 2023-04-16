package domain

import (
	"reflect"
	"testing"
)

var FolderStructure = &GrafanaFolder{
	Id:       0,
	Title:    "root",
	FolderId: -1,
	FolderItems: []*GrafanaFolder{
		&GrafanaFolder{
			Id:    10,
			Title: "Main",
			FolderItems: []*GrafanaFolder{
				&GrafanaFolder{
					Id:             40,
					Title:          "Main sub 1",
					FolderItems:    []*GrafanaFolder{},
					DashboardItems: []*GrafanaDashboard{},
				},
			},
			DashboardItems: []*GrafanaDashboard{},
		},
		&GrafanaFolder{
			Id:    20,
			Title: "Secondary",
			FolderItems: []*GrafanaFolder{
				&GrafanaFolder{
					Id:             25,
					Title:          "Secondary sub 1",
					FolderItems:    []*GrafanaFolder{},
					DashboardItems: []*GrafanaDashboard{},
				},
				&GrafanaFolder{
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

func TestGrafanaFolder_FindFolderById(t *testing.T) {
	tests := []struct {
		name    string
		folder  *GrafanaFolder
		findId  int
		want    *GrafanaFolder
		wantErr bool
	}{
		{
			name:    "find success 1",
			folder:  FolderStructure,
			findId:  0,
			want:    FolderStructure,
			wantErr: false,
		},
		{
			name:    "find success 2",
			folder:  FolderStructure,
			findId:  40,
			want:    FolderStructure.FolderItems[0].FolderItems[0],
			wantErr: false,
		},
		{
			name:    "find success 3",
			folder:  FolderStructure,
			findId:  30,
			want:    FolderStructure.FolderItems[1].FolderItems[1],
			wantErr: false,
		},
		{
			name:    "find failure",
			folder:  FolderStructure,
			findId:  1000000,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.folder.FindFolderById(tt.findId)
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
