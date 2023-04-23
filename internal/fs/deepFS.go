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
	tmpStruct := struct {
		Dashboard struct {
			Title  string
			Panels []domain.GrafanaPanel
		}
	}{}

	err = json.Unmarshal(content, &tmpStruct)
	if err != nil {
		return err
	}
	panels := tmpStruct.Dashboard.Panels
	for _, panel := range panels {
		panelDir := joinEscaping(dashboardFolder, fmt.Sprintf("%s (id=%v)", panel.Title, panel.Uid))
		err = fs.writer.CreateDir(panelDir)
		for _, query := range panel.GetQueries() {
			err = fs.writer.CreateFile(
				joinEscaping(panelDir, fmt.Sprintf("%s.%s", query.Title, query.Type)),
				[]byte(query.Expression),
			)
			if err != nil {
				return err
			}
		}
	}

	// reformatting dashboard json and saving it
	parsed := map[string]interface{}{}
	err = json.Unmarshal(content, &parsed)
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

func formatJson(m map[string]interface{}) ([]byte, error) {
	buff := bytes.Buffer{}
	encoder := json.NewEncoder(&buff)
	encoder.SetIndent("", "  ")
	encoder.Encode(m)
	return buff.Bytes(), nil
}
