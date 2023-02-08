package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/brotifypacha/grafana_searcher/internal/domain"
	"github.com/brotifypacha/grafana_searcher/internal/grafana"
)

type SingleFolderFS struct {
	client grafana.Repository
	writer Writer
}

func NewSingleFolderFS(
	client grafana.Repository,
	writer Writer,
) *SingleFolderFS {
	return &SingleFolderFS{
		client: client,
		writer: writer,
	}
}

func (fs *SingleFolderFS) Save(folder *domain.GrafanaFolder, path string) error {
	folderPath := smartJoin(folder.FolderId, path, replaceSpecials(folder.Title))
	if folder.FolderId == -1 {
		err := fs.writer.CreateDir(folderPath)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("Couldn't create dir '%s': %w", folderPath, err)
		}
	}
	for _, subFolder := range folder.FolderItems {
		err := fs.Save(subFolder, folderPath)
		if err != nil {
			return fmt.Errorf("error saving folder: %w", err)
		}
	}
	for _, dashboard := range folder.DashboardItems {
		bytes, err := fs.client.GetDashboard(dashboard.Uid)
		if err != nil {
			return fmt.Errorf("error getting dashboard: %w", err)
		}
		filePath := smartJoin(dashboard.FolderId, folderPath, replaceSpecials(dashboard.Title)) + ".json"
		err = fs.writer.CreateFile(filePath, bytes)
		if err != nil {
			return fmt.Errorf("couldn't write file '%s': %w", filePath, err)
		}
	}
	return nil
}

func smartJoin(folderId int, path string, child string) string {
	if folderId == -1 || folderId == 0 {
		return filepath.Join(path, child)
	} else {
		return fmt.Sprintf("%s__%s", path, child)
	}
}

func replaceSpecials(str string) string {
	return strings.NewReplacer(
		" ", "_",
		"/", "_",
	).Replace(str)
}
