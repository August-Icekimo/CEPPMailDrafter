package loader

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	ErrDirMissing  = errors.New("template directory is missing")
	ErrFileMissing = errors.New("template file does not exist")
)

// TemplateLoader defines how templates are loaded.
type TemplateLoader interface {
	Load(month string) (string, error)
}

// FileLoader implements TemplateLoader using the local filesystem.
type FileLoader struct {
	dir string
}

// NewFileLoader creates a new FileLoader with the specified directory.
func NewFileLoader(dir string) *FileLoader {
	return &FileLoader{dir: dir}
}

// Load reads the template file for the given month from the configured directory.
func (fl *FileLoader) Load(month string) (string, error) {
	info, err := os.Stat(fl.dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", ErrDirMissing
		}
		return "", fmt.Errorf("stat template dir: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("template path is not a directory: %s", fl.dir)
	}

	ext := ".md"
	filePath := filepath.Join(fl.dir, month+ext)

	content, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("%w: expected path %s", ErrFileMissing, filePath)
		}
		return "", fmt.Errorf("read template file: %w", err)
	}

	return string(content), nil
}
