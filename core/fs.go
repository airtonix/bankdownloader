package core

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/adrg/xdg"
	"github.com/airtonix/bank-downloaders/meta"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
)

func SortedKeys[K constraints.Ordered, V any](m map[K]V) []K {
	keys := lo.Keys(m)
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func GetUserFilePath(filename string) string {
	xdgDocumentDir := filepath.Join(xdg.UserDirs.Documents, meta.Name)
	xdgDocumentPath := filepath.Join(xdgDocumentDir, filename)

	return xdgDocumentPath
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func Dos2Unix(str string) string {
	return strings.ReplaceAll(str, "\r\n", "\n")
}
