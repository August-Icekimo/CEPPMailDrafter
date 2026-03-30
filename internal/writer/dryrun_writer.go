package writer

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// DryRunWriter implements Writer by printing a preview to an io.Writer.
type DryRunWriter struct {
	out io.Writer
}

func NewDryRunWriter(out io.Writer) *DryRunWriter {
	if out == nil {
		out = os.Stdout
	}
	return &DryRunWriter{out: out}
}

func (w *DryRunWriter) Write(msg MailMessage, month string) error {
	fmt.Fprintf(w.out, "[Dry Run Mode] Previewing output for: %s\n", month)
	fmt.Fprintf(w.out, strings.Repeat("=", 60) + "\n")
	
	fmt.Fprintf(w.out, "From: %s\n", msg.From)
	if len(msg.To) > 0 {
		fmt.Fprintf(w.out, "To: %s\n", strings.Join(msg.To, ", "))
	}
	if len(msg.Cc) > 0 {
		fmt.Fprintf(w.out, "Cc: %s\n", strings.Join(msg.Cc, ", "))
	}
	if len(msg.Bcc) > 0 {
		fmt.Fprintf(w.out, "Bcc: %s\n", strings.Join(msg.Bcc, ", "))
	}
	fmt.Fprintf(w.out, "Subject: %s\n", msg.Subject)
	
	if len(msg.Attachments) > 0 {
		fmt.Fprintf(w.out, "Attachments:\n")
		for _, a := range msg.Attachments {
			fmt.Fprintf(w.out, "  - %s (%s, %d bytes)\n", a.Filename, a.ContentType, len(a.Data))
		}
	}
	
	fmt.Fprintf(w.out, "\n--- BODY ---\n")
	fmt.Fprintf(w.out, "%s\n", msg.Body)
	fmt.Fprintf(w.out, strings.Repeat("=", 60) + "\n")
	
	return nil
}
