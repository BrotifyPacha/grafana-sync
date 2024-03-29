package grafana

import (
	"encoding/json"
	"fmt"

	"github.com/brotifypacha/grafana-sync/internal/domain"
	"github.com/brotifypacha/grafana-sync/internal/grafana/miniGrafanaClient"
)

type WebRepository struct {
	client *miniGrafanaClient.Client
}

func NewRepository(client *miniGrafanaClient.Client) *WebRepository {
	return &WebRepository{
		client: client,
	}
}

func (d *WebRepository) GetTree() (dashboards *domain.GrafanaFolder, err error) {
	bytes, err := d.client.Get("api/search")
	if err != nil {
		return
	}
	rawApiItems := []*RawGrafanaApiItem{}
	err = json.Unmarshal(bytes, &rawApiItems)
	if err != nil {
		return
	}
	rootFolder := buildTree(rawApiItems)
	return rootFolder, nil
}

func buildTree(grafanaItems []*RawGrafanaApiItem) *domain.GrafanaFolder {
	rootFolder := &domain.GrafanaFolder{
		Uid:            domain.RootFolderUid,
		Id:             domain.RootFolderId,
		Title:          "Root folder",
		FolderId:       -1,
		FolderItems:    make([]*domain.GrafanaFolder, 0),
		DashboardItems: make([]*domain.GrafanaDashboard, 0),
	}

	folders := make(map[int]*domain.GrafanaFolder, 0)
	dashboards := make(map[int]*domain.GrafanaDashboard, 0)

	folders[0] = rootFolder

	for i := range grafanaItems {
		item := grafanaItems[i]
		switch grafanaItems[i].Type {
		case ITEM_TYPE_FOLDER:
			{
				folder := &domain.GrafanaFolder{
					Uid:            item.Uid,
					Id:             item.Id,
					Title:          item.Title,
					FolderId:       item.FolderId,
					FolderItems:    make([]*domain.GrafanaFolder, 0),
					DashboardItems: make([]*domain.GrafanaDashboard, 0),
				}
				folders[item.Id] = folder
			}
		case ITEM_TYPE_DASHBORAD:
			{
				dashboard := &domain.GrafanaDashboard{
					Uid:      item.Uid,
					Title:    item.Title,
					FolderId: item.FolderId,
				}
				dashboards[item.Id] = dashboard
			}
		}
	}

	for i := range dashboards {
		_, ok := folders[dashboards[i].FolderId]
		if !ok {
			fmt.Printf("Dashboard is referring to inexisting folder id, %#v", dashboards[i])
			continue
		}
		folders[dashboards[i].FolderId].DashboardItems = append(folders[dashboards[i].FolderId].DashboardItems, dashboards[i])
	}

	for i := range folders {
		if folders[i].Id == domain.RootFolderId {
			continue
		}
		_, ok := folders[folders[i].FolderId]
		if !ok {
			fmt.Printf("Folder is referring to inexisting folder id, %#v", folders[i])
			continue
		}
		folders[folders[i].FolderId].FolderItems = append(folders[folders[i].FolderId].FolderItems, folders[i])
	}
	return rootFolder
}

func (d *WebRepository) GetDashboard(uid string) (*domain.GrafanaDashboardDetails, error) {
	bytes, err := d.client.Get(fmt.Sprintf("/api/dashboards/uid/%s", uid))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	db := &domain.GrafanaDashboardDetails{}
	tmp := struct {
		Dashboard *domain.GrafanaDashboardDetails
	}{
		Dashboard: db,
	}
	err = json.Unmarshal(bytes, &tmp)
	db.RawData = bytes
	return db, err
}
