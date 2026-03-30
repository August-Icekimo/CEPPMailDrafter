package renderer

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/icekimo/CEPPMailDrafter/internal/data"
	"github.com/icekimo/CEPPMailDrafter/internal/parser"
	"github.com/icekimo/CEPPMailDrafter/internal/writer"
	"gopkg.in/yaml.v3"
)

// Renderer holds data source and logic for processing a template.
type Renderer struct {
	source data.DataSource
}

// New creates a Renderer with the provided data sources.
func New(sources ...data.DataSource) *Renderer {
	if len(sources) == 1 {
		return &Renderer{source: sources[0]}
	}
	return &Renderer{source: data.NewChainedSource(sources...)}
}

// Render processes front matter and body tokens into a MailMessage.
func (r *Renderer) Render(frontYAML []byte, tokens []parser.Token) (writer.MailMessage, error) {
	var fm map[string]any
	if err := yaml.Unmarshal(frontYAML, &fm); err != nil {
		return writer.MailMessage{}, fmt.Errorf("parse front matter yaml: %w", err)
	}

	from, _ := fm["from"].(string)

	to := parseList(fm["to"])
	cc := parseList(fm["cc"])
	bcc := parseList(fm["bcc"])

	subjectRaw, _ := fm["subject"].(string)
	subjectTokens, err := parser.Tokenise(subjectRaw)
	if err != nil {
		return writer.MailMessage{}, fmt.Errorf("parse subject tag: %w", err)
	}
	subjectStr := r.renderTokensString(subjectTokens)

	var attachments []writer.Attachment
	if attsRaw, ok := fm["attachments"]; ok {
		attPaths := parseList(attsRaw)
		for _, p := range attPaths {
			b, err := os.ReadFile(p)
			if err != nil {
				return writer.MailMessage{}, fmt.Errorf("read attachment %s: %w", p, err)
			}
			contentType := http.DetectContentType(b)
			attachments = append(attachments, writer.Attachment{
				Filename:    filepath.Base(p),
				ContentType: contentType,
				Data:        b,
			})
		}
	}

	bodyStr := r.renderTokensString(tokens)

	return writer.MailMessage{
		From:        from,
		To:          to,
		Cc:          cc,
		Bcc:         bcc,
		Subject:     subjectStr,
		Body:        bodyStr,
		Attachments: attachments,
	}, nil
}

func parseList(val any) []string {
	if val == nil {
		return nil
	}
	if s, ok := val.(string); ok {
		return []string{s}
	}
	if slice, ok := val.([]any); ok {
		var res []string
		for _, v := range slice {
			res = append(res, fmt.Sprintf("%v", v))
		}
		return res
	}
	return nil
}

func (r *Renderer) renderTokensString(tokens []parser.Token) string {
	var buf bytes.Buffer
	r.renderTokens(&buf, tokens, data.NewEmptySource())
	return buf.String()
}

func (r *Renderer) renderTokens(w *bytes.Buffer, tokens []parser.Token, loopItemCtx data.DataSource) {
	// loopItemCtx handles {{ITEM}} logic overriding r.source selectively.
	for _, tok := range tokens {
		switch tok.Type {
		case parser.TokenText:
			w.WriteString(tok.Raw)
		case parser.TokenVar:
			if loopItemCtx != nil {
				if v, ok := loopItemCtx.Get(tok.Key); ok {
					w.WriteString(v)
					continue
				}
			}
			if v, ok := r.source.Get(tok.Key); ok {
				w.WriteString(v)
			} else {
				fmt.Fprintf(os.Stderr, "[WARN] missing value for tag %s\n", tok.Raw)
				w.WriteString(tok.Raw) // keep original
			}
		case parser.TokenLoop:
			list, ok := r.source.GetList(tok.Key)
			if !ok {
				// empty list silently does nothing
				continue
			}
			for _, item := range list {
				// We need a tiny DataSource just for {{ITEM}}
				itemSource := &singleItemSource{item: item}
				r.renderTokens(w, tok.Body, itemSource)
			}
		case parser.TokenCond:
			v, ok := r.source.Get(tok.Key)
			if ok && isTruthy(v) {
				r.renderTokens(w, tok.Body, loopItemCtx)
			}
		}
	}
}

type singleItemSource struct {
	item string
}

func (s *singleItemSource) Get(key string) (string, bool) {
	if strings.ToUpper(key) == "ITEM" {
		return s.item, true
	}
	return "", false
}
func (s *singleItemSource) GetList(string) ([]string, bool) { return nil, false }

func isTruthy(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "true" || v == "1" || v == "yes"
}
