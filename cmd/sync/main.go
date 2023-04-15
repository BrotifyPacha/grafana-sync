package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/brotifypacha/grafana-sync/internal/fs"
	"github.com/brotifypacha/grafana-sync/internal/fs/writer"
	"github.com/brotifypacha/grafana-sync/internal/grafana"
	"github.com/brotifypacha/grafana-sync/internal/grafana/miniGrafanaClient"
)

const (
	pathHelpText = "`local path` to directory where dashboards' requests will be stored"
	hostHelpText = "grafana HTTP api `host`"
)

func main() {
	var localPath string
	var grafanaHost string
	var folderId int

	set := flag.NewFlagSet("sync", flag.ContinueOnError)
	set.StringVar(&localPath, "l", "", pathHelpText)
	set.StringVar(&grafanaHost, "h", "", hostHelpText)
	set.IntVar(&folderId, "folderId", grafana.RootRepositoryId, "")
	set.SetOutput(os.Stdout)

	err := set.Parse(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}

	if len(localPath) == 0 || len(grafanaHost) == 0 {
		set.PrintDefaults()
		return
	}

	client, err := miniGrafanaClient.NewClient(grafanaHost)
	if err != nil {
		fmt.Println(err)
		return
	}
	repo := grafana.NewRepository(client)

	tree, err := repo.GetTree()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	folderToSync, err := tree.FindFolderById(folderId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileWriter := writer.NewLocalWriter()
	fs := fs.NewDeepFolderFs(repo, fileWriter)
	errs := fs.Save(*folderToSync, localPath)
	if len(errs) != 0 {
		for _, err = range errs {
			if err != nil {
				fmt.Println(err)
			}
		}
		os.Exit(1)
	}
}
