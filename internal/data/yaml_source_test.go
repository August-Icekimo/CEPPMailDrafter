package data

import (
	"os"
	"path/filepath"
	"testing"
)

func TestYAMLSource(t *testing.T) {
	tempDir := t.TempDir()
	
	month := "2025-01"
	yamlContent := `
month: "一月"
ITEMS:
  - "項目A"
  - "項目B"
`
	err := os.WriteFile(filepath.Join(tempDir, month+".yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		src, err := NewYAMLSource(month, tempDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		v, ok := src.Get("month")
		if !ok || v != "一月" {
			t.Errorf("expected '一月', got %v", v)
		}

		v2, ok := src.Get("MONTH") // case insensitive check
		if !ok || v2 != "一月" {
			t.Errorf("expected '一月' from uppercase, got %v", v2)
		}

		list, ok := src.GetList("items")
		if !ok || len(list) != 2 || list[0] != "項目A" {
			t.Errorf("expected list ['項目A', '項目B'], got %v", list)
		}
	})

	t.Run("missing fallback", func(t *testing.T) {
		src, err := NewYAMLSource("missing", tempDir)
		if err != nil {
			t.Fatalf("expected no error for missing file, got %v", err)
		}
		
		_, ok := src.Get("anything")
		if ok {
			t.Errorf("expected nothing from missing source")
		}
	})
}
