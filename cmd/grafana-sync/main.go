package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/brotifypacha/grafana-sync/internal/domain"
	"github.com/brotifypacha/grafana-sync/internal/fs"
	"github.com/brotifypacha/grafana-sync/internal/fs/writer"
	"github.com/brotifypacha/grafana-sync/internal/grafana"
	"github.com/brotifypacha/grafana-sync/internal/grafana/miniGrafanaClient"
)

const (
	pathHelpText      = "`local path` to directory where dashboards' requests will be stored"
	hostHelpText      = "grafana HTTP api `host`"
	folderUidHelpText = "grafana folder UID to use for this command"

	rootUsageText = `Usage: grafana-sync command [command-flags...]

Commands:
	print     - prints grafana structure tree to stdout
	download  - downloads grafana structure to specified dir`
)

func main() {
	var grafanaHost string
	var folderUid string

	var printRecursive bool
	printCmd := flag.NewFlagSet("print", flag.ContinueOnError)
	printCmd.BoolVar(&printRecursive, "recursive", true, "print tree structure recursivly")
	setBaseFlags(printCmd, &grafanaHost, &folderUid)

	var localPath string
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	downloadCmd.StringVar(&localPath, "path", "", pathHelpText)
	setBaseFlags(downloadCmd, &grafanaHost, &folderUid)

	if len(os.Args) < 1 {
		fmt.Println(rootUsageText)
		return
	}

	switch os.Args[1] {
	case printCmd.Name():
		{
			err := printCmd.Parse(os.Args[2:])
			if errors.Is(err, flag.ErrHelp) || len(grafanaHost) == 0 {
				fmt.Printf("Usage: grafana-sync %s [flags...]\n\n", printCmd.Name())
				printCmd.PrintDefaults()
				return
			}
		}
	case downloadCmd.Name():
		{
			err := downloadCmd.Parse(os.Args[2:])
			if errors.Is(err, flag.ErrHelp) || len(grafanaHost) == 0 || len(localPath) == 0 {
				fmt.Printf("Usage: grafana-sync %s [flags...]\n\n", downloadCmd.Name())
				downloadCmd.PrintDefaults()
				return
			}
		}
	default:
		{
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println(rootUsageText)
			return
		}
	}

	client, err := miniGrafanaClient.NewClient(grafanaHost)
	ExitOnErrors(err)

	repo := grafana.NewRepository(client)
	tree, err := repo.GetTree()
	ExitOnErrors(err)

	pickedFolder, err := tree.FindFolderByUid(folderUid)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch printCmd.Name() {
	case printCmd.Name():
		{
			fmt.Println(pickedFolder.PrettySprint(printRecursive))
		}
	case downloadCmd.Name():
		{
			fs := fs.NewDeepFolderFs(repo, writer.NewLocalWriter())
			errs := fs.Save(*pickedFolder, localPath)
			ExitOnErrors(errs...)
		}
	}

}

func ExitOnErrors(errs ...error) {
	hasErrs := false
	for _, err := range errs {
		if err != nil {
			hasErrs = true
			fmt.Println(err)
		}
	}
	if hasErrs {
		os.Exit(1)
	}
}

func setBaseFlags(flagSet *flag.FlagSet, host *string, folderUid *string) {
	flagSet.SetOutput(os.Stdout)
	flagSet.StringVar(host, "host", "", hostHelpText)
	flagSet.StringVar(folderUid, "folder-uid", domain.RootFolderUid, folderUidHelpText)
	flagSet.Usage = func() {}
}
