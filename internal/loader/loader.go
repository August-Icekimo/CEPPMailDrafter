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
	// Load reads <dir>/<month>.md
	Load(month string) (string, error)
	// LoadFile reads <dir>/<filename> (filename may include extension).
	LoadFile(dir, filename string) (string, error)
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
// The file path resolved is <dir>/<month>.md.
func (fl *FileLoader) Load(month string) (string, error) {
	return fl.LoadFile(fl.dir, month+".md")
}

// LoadFile reads an arbitrary template file given an explicit directory and filename.
// filename may include a path extension; if it has none, ".md" is NOT appended.
func (fl *FileLoader) LoadFile(dir, filename string) (string, error) {
	info, err := os.Stat(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", ErrDirMissing
		}
		return "", fmt.Errorf("stat template dir: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("template path is not a directory: %s", dir)
	}

	filePath := filepath.Join(dir, filename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("%w: expected path %s", ErrFileMissing, filePath)
		}
		return "", fmt.Errorf("read template file: %w", err)
	}

	return string(content), nil
}
