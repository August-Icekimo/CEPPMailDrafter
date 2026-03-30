package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEMLWriter(t *testing.T) {
	outDir := t.TempDir()
	w := NewEMLWriter(outDir)

	msg := MailMessage{
		From:    "sender@test.com",
		To:      []string{"a@test.com"},
		Subject: "中文測試",
		Body:    "<b>測試</b>",
	}

	err := w.Write(msg, "2025-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, err := os.ReadFile(filepath.Join(outDir, "2025-01.eml"))
	if err != nil {
		t.Fatalf("could not read written file: %v", err)
	}

	outStr := string(b)
	if !strings.Contains(outStr, "From: <sender@test.com>") && !strings.Contains(outStr, "From: sender@test.com") {
		t.Errorf("missing From header: %s", outStr)
	}
	if !strings.Contains(outStr, "Subject: =?UTF-8?") {
		t.Errorf("subject not b-encoded: %s", outStr)
	}
	if !strings.Contains(outStr, "Content-Type: text/html") {
		t.Errorf("missing inline body content type")
	}
	if !strings.Contains(outStr, "<b>測試</b>") {
		t.Errorf("missing body")
	}
}
