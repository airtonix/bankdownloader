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

type Cwd *string

func GetCwd() Cwd {
	cwd, err := os.Getwd()
	AssertErrorToNilf("could not get cwd: %w", err)
	return &cwd
}

// ResolveFileArg resolves the filename to use for a config file
// based on the following rules:
// 1. envvar runtime override
// 2. args filename override
// 3. config file in current directory
// 4. config file in XDG directory
func ResolveFileArg(
	argFilename string,
	envvarKey string,
	defaultFilename string,
) string {
	xdgFilepath := GetUserFilePath(defaultFilename)
	envFilepath := os.Getenv(envvarKey)

	// envvar runtime override
	if envFilepath != "" {
		return envFilepath
	}

	// args filename override
	if argFilename != "" {
		return argFilename
	}

	// config file in current directory
	if FileExists(defaultFilename) {
		return defaultFilename
	}

	// config file in XDG directory
	return xdgFilepath
}
