package fs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/brotifypacha/grafana_searcher/internal/domain"
	"github.com/brotifypacha/grafana_searcher/internal/fs/writer"
	"github.com/brotifypacha/grafana_searcher/internal/grafana"
)

type DeepFS struct {
	client grafana.Repository
	writer writer.Writer
}

func NewDeepFolderFs(
	client grafana.Repository,
	writer writer.Writer,
) FileSystemInterface {
	return &DeepFS{
		client: client,
		writer: writer,
	}
}

func (fs *DeepFS) Save(grafanaFolder domain.GrafanaFolder, localFolder string) error {

	folderPath := filepath.Join(localFolder, grafanaFolder.Title)
	fs.writer.CreateDir(folderPath)

	wg := sync.WaitGroup{}
	sem := make(chan int, 8)
	for _, folder := range grafanaFolder.FolderItems {
		wg.Add(1)
		sem <- 1
		go func(folder *domain.GrafanaFolder, path string) {
			fs.Save(*folder, folderPath)
			wg.Done()
			<-sem
		}(folder, folderPath)
	}

	for _, dashboard := range grafanaFolder.DashboardItems {
		err := fs.saveDashboard(dashboard, folderPath)
		if err != nil {
			fmt.Println("DeepFS: error writing dashboard: ", err.Error())
			os.Exit(1)
			continue
		}
		fmt.Println("DeepFS: written:", folderPath+"/"+dashboard.Title)
	}

	wg.Wait()
	return nil
}

func (fs *DeepFS) saveDashboard(dashboard *domain.GrafanaDashboard, path string) error {
	dashboardFolder := filepath.Join(path, dashboard.Title)
	fs.writer.CreateDir(dashboardFolder)

	content, err := fs.client.GetDashboard(dashboard.Uid)
	if err != nil {
		return err
	}
	content, err = formatJson(content)
	if err != nil {
		return err
	}

	dashboardFile := filepath.Join(dashboardFolder, "dashboard-data.json")
	fs.writer.CreateFile(dashboardFile, content)
	return nil
}

func formatJson(b []byte) ([]byte, error) {
	parsed := map[string]interface{}{}
	err := json.Unmarshal(b, &parsed)
	if err != nil {
		return nil, err
	}

	buff := bytes.Buffer{}
	encoder := json.NewEncoder(&buff)
	encoder.SetIndent("", "  ")
	encoder.Encode(parsed)
	return buff.Bytes(), nil
}
