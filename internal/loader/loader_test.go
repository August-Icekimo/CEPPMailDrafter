package loader

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileLoader_Load(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test template
	month := "2025-01"
	testContent := "test template content"
	err := os.WriteFile(filepath.Join(tempDir, month+".md"), []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		l := NewFileLoader(tempDir)
		got, err := l.Load(month)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != testContent {
			t.Errorf("expected %q, got %q", testContent, got)
		}
	})

	t.Run("missing directory", func(t *testing.T) {
		l := NewFileLoader(filepath.Join(tempDir, "nonexistent"))
		_, err := l.Load(month)
		if !errors.Is(err, ErrDirMissing) {
			t.Errorf("expected ErrDirMissing, got %v", err)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		l := NewFileLoader(tempDir)
		_, err := l.Load("non-existent-month")
		if !errors.Is(err, ErrFileMissing) {
			t.Errorf("expected ErrFileMissing, got %v", err)
		}
	})
}
