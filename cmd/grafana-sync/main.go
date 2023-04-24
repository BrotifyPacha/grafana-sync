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
	var localPath string
	var grafanaHost string
	var folderUid string

	var printRecursive bool
	printCmd := flag.NewFlagSet("print", flag.ContinueOnError)
	printCmd.BoolVar(&printRecursive, "recursive", true, "print tree structure recursivly")
	setBaseFlags(printCmd, &grafanaHost, &localPath, &folderUid)

	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	setBaseFlags(syncCmd, &grafanaHost, &localPath, &folderUid)

	if len(os.Args) < 1 {
		fmt.Println(rootUsageText)
		return
	}

	var cmd *flag.FlagSet
	switch os.Args[1] {
	case printCmd.Name():
		{
			cmd = printCmd
		}
	case syncCmd.Name():
		{
			cmd = syncCmd
		}
	default:
		{
			fmt.Println(rootUsageText)
			return
		}
	}

	err := cmd.Parse(os.Args[2:])
	if errors.Is(err, flag.ErrHelp) || len(localPath) == 0 || len(grafanaHost) == 0 {
		fmt.Printf("Usage: grafana-sync %s [flags...]\n\n", cmd.Name())
		cmd.PrintDefaults()
		return
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

	switch cmd.Name() {
	case printCmd.Name():
		{
			pickedFolder.PrettyPrint(printRecursive)
		}
	case syncCmd.Name():
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

func setBaseFlags(flagSet *flag.FlagSet, host *string, path *string, folderUid *string) {
	flagSet.SetOutput(os.Stdout)
	flagSet.StringVar(path, "path", "", pathHelpText)
	flagSet.StringVar(host, "host", "", hostHelpText)
	flagSet.StringVar(folderUid, "folder-uid", domain.RootFolderUid, folderUidHelpText)
	flagSet.Usage = func() {

	}
}
