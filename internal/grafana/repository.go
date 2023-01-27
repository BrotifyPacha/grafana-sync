package grafana

import (
	"encoding/json"
	"fmt"
	"github.com/brotifypacha/grafana_searcher/internal/grafana/miniGrafanaClient"
	"github.com/brotifypacha/grafana_searcher/internal/domain"
)

type DashboardRepoInterface interface {
	GetTree() (*domain.GrafanaFolder, error)
}

type Repository struct {
	client *miniGrafanaClient.Client
}

func NewRepository(client *miniGrafanaClient.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (d *Repository) GetTree() (dashboards *domain.GrafanaFolder, err error) {
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
		Id: 0,
		Title: "Root folder",
		FolderId: -1,
		FolderItems: make([]*domain.GrafanaFolder, 0),
		DashboardItems: make([]*domain.GrafanaDashboard, 0),
	}

	folders := make(map[int]*domain.GrafanaFolder, 0)
	dashboards := make(map[int]*domain.GrafanaDashboard, 0)

	folders[0] = rootFolder

	for i := range grafanaItems {
		item := grafanaItems[i]
		// fmt.Printf("%#v\n", item)
		switch grafanaItems[i].Type {
		case ITEM_TYPE_FOLDER:
			{
				folder := &domain.GrafanaFolder{
					Id: item.Id,
					Title: item.Title,
					FolderId: item.FolderId,
					FolderItems: make([]*domain.GrafanaFolder, 0),
					DashboardItems: make([]*domain.GrafanaDashboard, 0),
				}
				folders[item.Id] = folder
			}
		case ITEM_TYPE_DASHBORAD:
			{
				dashboard := &domain.GrafanaDashboard{
					Id: item.Id,
					Title: item.Title,
					FolderId: item.FolderId,
				}
				dashboards[item.Id] = dashboard
			}
		}
	}

	for i := range dashboards {
		folders[dashboards[i].FolderId].DashboardItems = append(folders[dashboards[i].FolderId].DashboardItems, dashboards[i])
	}

	for i := range folders {
		if folders[i].FolderId == -1 {
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

func printChildren(
	parent     *domain.GrafanaFolder,
	indent     string,
) {
	for i := range parent.DashboardItems {
		db := parent.DashboardItems[i]
		isLast := i == len(parent.DashboardItems) - 1 && len(parent.FolderItems) == 0
		if isLast {
			fmt.Printf("%s╙── %s [%d]\n", indent, db.Title, db.Id)
		} else {
			fmt.Printf("%s╟── %s [%d]\n", indent, db.Title, db.Id)
		}
	}
	for i := range parent.FolderItems {
		folder := parent.FolderItems[i]

		notLast := i != len(parent.FolderItems) - 1
		hasChildren := len(folder.FolderItems) != 0 || len(folder.DashboardItems) != 0
		if notLast {
			if hasChildren {
				fmt.Printf("%s╠═╦ %s [%d]\n", indent, folder.Title, folder.Id)
				printChildren(parent.FolderItems[i], indent + "║ ")
			} else {
				fmt.Printf("%s╠══ %s [%d]\n", indent, folder.Title, folder.Id)
			}
		} else {
			if hasChildren {
				fmt.Printf("%s╚═╦ %s [%d]\n", indent, folder.Title, folder.Id)
				printChildren(parent.FolderItems[i], indent + "  ")
			} else {
				fmt.Printf("%s╚══ %s [%d]\n", indent, folder.Title, folder.Id)
			}
		}
	}
}
