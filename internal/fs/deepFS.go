package fs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/brotifypacha/grafana-sync/internal/domain"
	"github.com/brotifypacha/grafana-sync/internal/fs/writer"
	"github.com/brotifypacha/grafana-sync/internal/grafana"
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

	folderPath := joinEscaping(path, grafanaFolder.Title)
	err := fs.writer.CreateDir(folderPath)
	if err != nil {
		errCh <- err
	}

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
	dashboardFolder := joinEscaping(path, dashboard.Title)
	err := fs.writer.CreateDir(dashboardFolder)
	if err != nil {
		return err
	}

	content, err := fs.client.GetDashboard(dashboard.Uid)
	if err != nil {
		return err
	}
	parsed := map[string]interface{}{}
	err = json.Unmarshal(content, &parsed)
	if err != nil {
		return err
	}
	panels := getPanels(parsed)
	for _, panel := range panels {
		panelDir := joinEscaping(dashboardFolder, panel.Title)
		fs.writer.CreateDir(panelDir)
		for _, query := range panel.Queries {
			err = fs.writer.CreateFile(joinEscaping(panelDir, query.Title+".sql"), query.Data)
			if err != nil {
				return err
			}
		}
	}

	content, err = formatJson(parsed)
	if err != nil {
		return err
	}

	dashboardFile := joinEscaping(dashboardFolder, "dashboard-data.json")
	fs.writer.CreateFile(dashboardFile, content)
	return nil
}

type Panel struct {
	Title   string
	Queries []Query
}

type Query struct {
	Title string
	Data  []byte
}

func getPanels(m map[string]interface{}) (result []Panel) {
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

func getQueriesPanels(panels []interface{}) (result []Panel) {
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
			for _, v := range queries {
				result = append(result, v)
			}
		}
		targets, ok := panel["targets"].([]interface{})
		if !ok {
			continue
		}
		panelItem := Panel{}
		if panel["title"] != nil {
			title := panel["title"].(string)
			id := panel["id"].(float64)
			panelItem.Title = fmt.Sprintf("%s (id=%v)", title, id)
		}
		for _, targetI := range targets {
			target, ok := targetI.(map[string]interface{})
			if !ok {
				continue
			}
			rawSql, ok := target["rawSql"]
			refId := target["refId"].(string)
			if ok {
				panelItem.Queries = append(panelItem.Queries, Query{
					Title: refId,
					Data:  []byte(rawSql.(string)),
				})
			}
		}
		result = append(result, panelItem)
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
