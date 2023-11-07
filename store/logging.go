package store

import (
	"os"
	"path/filepath"

	"github.com/airtonix/bank-downloaders/core"
	"github.com/airtonix/bank-downloaders/meta"
)

func EnsureLogFilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if core.AssertErrorToNilf("could not get user cache directory: %w", err) {
		return "", err
	}

	path := filepath.Join(cacheDir, meta.Name, meta.Name+".log")

	err = os.MkdirAll(filepath.Dir(path), 0750)
	if core.AssertErrorToNilf("could not create log file directory: %w", err) {
		return "", err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if core.AssertErrorToNilf("could not create log file: %w", err) {
		return "", err
	}

	err = file.Close()
	if core.AssertErrorToNilf("could not close log file: %w", err) {
		return "", err
	}

	return path, err
}
