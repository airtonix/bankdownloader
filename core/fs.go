package core

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/adrg/xdg"
	"github.com/airtonix/bank-downloaders/meta"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
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
	return err == nil
}

func Dos2Unix(str string) string {
	return strings.ReplaceAll(str, "\r\n", "\n")
}

func GetCwd() string {
	cwd, err := os.Getwd()
	AssertErrorToNilf("could not get cwd: %w", err)
	return cwd
}

func Slugify(thing string) string {
	return slug.Make(thing)
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
	logrus.Debug("envFilepath: ", envFilepath)
	logrus.Debug("argFilename: ", argFilename)

	// envvar runtime override
	if envFilepath != "" {
		logrus.Debug("using envFilepath")
		return envFilepath
	}

	// args filename override
	if argFilename != "" {
		logrus.Debug("using argFilename")
		return argFilename
	}

	// config file in current directory
	if FileExists(defaultFilename) {
		logrus.Debug("using defaultFilename")
		return defaultFilename
	}

	logrus.Debug("using xdgFilepath")
	// config file in XDG directory
	return xdgFilepath
}
