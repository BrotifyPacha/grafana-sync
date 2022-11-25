package dashboards

const (
	ITEM_TYPE_FOLDER    = "dash-folder"
	ITEM_TYPE_DASHBORAD = "dash-db"
)

type RawGrafanaApiItem struct {
	Id          int
	Uid         string
	Title       string
	Uri         string
	Url         string
	Slug        string
	Type        string
	Tags        []string
	IsStarred   bool
	FolderId    int
	FolderUid   string
	FolderTitle string
	FolderUrl   string
	SortMeta    int
}

type GrafanaDashboard struct {
	Id       int
	Title    string
	FolderId int
}
