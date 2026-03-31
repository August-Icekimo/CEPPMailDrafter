package data

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// YAMLSource loads data from a yaml file. Keys are uppercase-normalized.
type YAMLSource struct {
	data map[string]any
}

// NewYAMLSource creates a new YAMLSource for a given month and data directory.
// It resolves the path as <dir>/<month>.yaml.
func NewYAMLSource(month string, dir string) (DataSource, error) {
	return NewYAMLSourceFromFile(dir, month+".yaml")
}

// NewYAMLSourceFromFile creates a new YAMLSource from an explicit directory and filename.
// filename should include the .yaml extension.
// If the file does not exist, it prints a warning to stderr and returns an empty source.
func NewYAMLSourceFromFile(dir, filename string) (DataSource, error) {
	path := filepath.Join(dir, filename)
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "[WARN] data file not found: %s — continuing with empty data\n", path)
			return NewEmptySource(), nil
		}
		return nil, fmt.Errorf("read yaml data: %w", err)
	}

	var raw map[string]any
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("parse yaml data %s: %w", path, err)
	}

	normalized := make(map[string]any)
	for k, v := range raw {
		normalized[strings.ToUpper(k)] = v
	}

	return &YAMLSource{data: normalized}, nil
}

func (y *YAMLSource) Get(key string) (string, bool) {
	val, ok := y.data[strings.ToUpper(key)]
	if !ok {
		return "", false
	}
	
	switch v := val.(type) {
	case string:
		return v, true
	default:
		return fmt.Sprintf("%v", v), true
	}
}

func (y *YAMLSource) GetList(key string) ([]string, bool) {
	val, ok := y.data[strings.ToUpper(key)]
	if !ok {
		return nil, false
	}

	slice, ok := val.([]any)
	if !ok {
		return nil, false
	}

	res := make([]string, 0, len(slice))
	for _, item := range slice {
		res = append(res, fmt.Sprintf("%v", item))
	}
	return res, true
}
