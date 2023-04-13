package fs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
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

const (
	downloadWorkers = 8
)

func (fs *DeepFS) Save(grafanaFolder domain.GrafanaFolder, localFolder string) []error {

	errCh := make(chan error, 5)
	go func() {
		fs.saveInternal(grafanaFolder, localFolder, errCh)
		close(errCh)
	}()

	errors := []error{}
	for err := range errCh {
		errors = append(errors, err)
	}

	return errors
}

func (fs *DeepFS) saveInternal(grafanaFolder domain.GrafanaFolder, path string, errCh chan error) {

	folderPath := filepath.Join(path, grafanaFolder.Title)
	fs.writer.CreateDir(folderPath)

	wg := sync.WaitGroup{}
	sem := make(chan int, downloadWorkers)

	for _, folder := range grafanaFolder.FolderItems {
		wg.Add(1)
		sem <- 1
		go func(folder *domain.GrafanaFolder, path string) {
			fs.saveInternal(*folder, folderPath, errCh)

			wg.Done()
			<-sem

		}(folder, folderPath)
	}
	for _, dashboard := range grafanaFolder.DashboardItems {
		err := fs.saveDashboard(dashboard, folderPath)
		if err != nil {
			err = fmt.Errorf("DeepFS: error saving dashboard: %w", err)
			errCh <- err
		}
	}
	wg.Wait()
}

func (fs *DeepFS) saveDashboard(dashboard *domain.GrafanaDashboard, path string) error {
	dashboardFolder := filepath.Join(path, dashboard.Title)
	fs.writer.CreateDir(dashboardFolder)

	content, err := fs.client.GetDashboard(dashboard.Uid)
	if err != nil {
		return err
	}
	parsed := map[string]interface{}{}
	err = json.Unmarshal(content, &parsed)
	if err != nil {
		return err
	}
	queries := getRawQueries(parsed)
	for title, query := range queries {
		fs.writer.CreateFile(filepath.Join(dashboardFolder, title+".sql"), query)
	}

	content, err = formatJson(parsed)
	if err != nil {
		return err
	}

	dashboardFile := filepath.Join(dashboardFolder, "dashboard-data.json")
	fs.writer.CreateFile(dashboardFile, content)
	return nil
}

func getRawQueries(m map[string]interface{}) map[string][]byte {
	result := map[string][]byte{}
	dashboard, ok := m["dashboard"].(map[string]interface{})
	if !ok {
		return result
	}
	panels, ok := dashboard["panels"].([]interface{})
	if !ok {
		return result
	}
	return getQueriesPanels(panels)
}

func getQueriesPanels(panels []interface{}) map[string][]byte {
	result := map[string][]byte{}
	for _, panelI := range panels {
		panel, ok := panelI.(map[string]interface{})
		if !ok {
			continue
		}
		if panel["type"] == "row" {
			panels, ok = panel["panels"].([]interface{})
			if !ok {
				continue
			}
			queries := getQueriesPanels(panels)
			for k, v := range queries {
				result[k] = v
			}
		}
		targets, ok := panel["targets"].([]interface{})
		if !ok {
			continue
		}
		title := ""
		if panel["title"] != nil {
			title = panel["title"].(string)
		}
		for i, targetI := range targets {
			target, ok := targetI.(map[string]interface{})
			if !ok {
				continue
			}
			rawSql, ok := target["rawSql"]
			if ok {
				queryName := title + " #" + strconv.Itoa(i)
				result[queryName] = []byte(rawSql.(string))
			}
		}
	}
	return result
}

func formatJson(m map[string]interface{}) ([]byte, error) {
	buff := bytes.Buffer{}
	encoder := json.NewEncoder(&buff)
	encoder.SetIndent("", "  ")
	encoder.Encode(m)
	return buff.Bytes(), nil
}
