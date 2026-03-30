package writer

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EMLWriter outputs the MailMessage to a local .eml file.
type EMLWriter struct {
	outputDir string
}

func NewEMLWriter(dir string) *EMLWriter {
	return &EMLWriter{outputDir: dir}
}

func (w *EMLWriter) Write(msg MailMessage, month string) error {
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	outPath := filepath.Join(w.outputDir, month+".eml")
	if _, err := os.Stat(outPath); err == nil {
		fmt.Fprintf(os.Stderr, "[WARN] Output file %s already exists, overwriting\n", outPath)
	}

	var buf bytes.Buffer

	// Headers
	buf.WriteString("From: " + msg.From + "\r\n")
	if len(msg.To) > 0 {
		buf.WriteString("To: " + strings.Join(msg.To, ", ") + "\r\n")
	}
	if len(msg.Cc) > 0 {
		buf.WriteString("Cc: " + strings.Join(msg.Cc, ", ") + "\r\n")
	}
	if len(msg.Bcc) > 0 {
		buf.WriteString("Bcc: " + strings.Join(msg.Bcc, ", ") + "\r\n")
	}
	buf.WriteString("Subject: " + mime.BEncoding.Encode("UTF-8", msg.Subject) + "\r\n")
	buf.WriteString("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("X-Unsent: 1\r\n")

	// Body / Attachments
	if len(msg.Attachments) == 0 {
		buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		buf.WriteString(msg.Body)
	} else {
		boundary := generateBoundary()
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundary))
		
		// preamble
		buf.WriteString("This is a multi-part message in MIME format.\r\n")

		// HTML body part
		buf.WriteString("\r\n--" + boundary + "\r\n")
		buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		buf.WriteString(msg.Body)

		// Attachments
		for _, att := range msg.Attachments {
			buf.WriteString("\r\n--" + boundary + "\r\n")
			buf.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", att.ContentType, att.Filename))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", att.Filename))
			
			// Base64 encode and split to 76 char lines
			encoded := base64.StdEncoding.EncodeToString(att.Data)
			for len(encoded) > 0 {
				chunk := 76
				if len(encoded) < chunk {
					chunk = len(encoded)
				}
				buf.WriteString(encoded[:chunk] + "\r\n")
				encoded = encoded[chunk:]
			}
		}

		// epilogue
		buf.WriteString("\r\n--" + boundary + "--\r\n")
	}

	return os.WriteFile(outPath, buf.Bytes(), 0644)
}

func generateBoundary() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("----=_Part_%x", b)
}

