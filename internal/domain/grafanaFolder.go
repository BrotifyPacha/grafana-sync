package domain

import (
	"fmt"
	"strings"
)

type GrafanaFolder struct {
	Id             int
	Title          string
	FolderId       int
	FolderItems    []*GrafanaFolder
	DashboardItems []*GrafanaDashboard
}

func (f *GrafanaFolder) PrettyPrint(recursive bool) {
	fmt.Println(prettySprint(f, "", recursive))
}

func prettySprint(folder *GrafanaFolder, indent string, recursive bool) string {
	out := strings.Builder{}

	selfIndent := ""
	if indent != "" {
		if mbsubstr(indent, -3, -1) == "║ " {
			selfIndent = strings.Replace(indent, "║  ", "╠═ ", 1)
		} else {
			selfIndent = strings.Replace(indent, "   ", "╚═ ", 1)
		}
	}
	out.WriteString(fmt.Sprintf("%s%s [%d]\n", selfIndent, folder.Title, folder.Id))

	for i := range folder.DashboardItems {
		db := folder.DashboardItems[i]
		isLast := i == len(folder.DashboardItems) - 1 && len(folder.FolderItems) == 0
		if isLast {
			out.WriteString(fmt.Sprintf("%s╙─ %s [%s]\n", indent, db.Title, db.Uid))
		} else {
			out.WriteString(fmt.Sprintf("%s╟─ %s [%s]\n", indent, db.Title, db.Uid))
		}
	}
	for i := range folder.FolderItems {

		notLast := i != len(folder.FolderItems) - 1

		if recursive {
			if notLast {
				out.WriteString(prettySprint(folder.FolderItems[i], indent + "║  ", recursive))
			} else {
				out.WriteString(prettySprint(folder.FolderItems[i], indent + "   ", recursive))
			}
		}
	}
	return out.String()
}

func mbsubstr(str string, from int, to int) string {
	runes := []rune(str)
	if from < 0 {
		from = len(runes) + from
	}
	if to < 0 {
		to = len(runes) + to
	}
	r := string(runes[from:to])
	return r
}
