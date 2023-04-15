package fs

import (
	"path/filepath"
	"strings"
)

func joinEscaping(mainPath string, toJoin string) string {
	return filepath.Join(mainPath, strings.ReplaceAll(toJoin, `/`, `-`))
}
