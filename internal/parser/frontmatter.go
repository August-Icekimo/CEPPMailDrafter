package parser

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoFrontMatter = errors.New("template has no front matter")
	ErrMissingKey    = errors.New("front matter is missing required key(s)")
)

// SplitFrontMatter separates the YAML front matter from the markdown body.
// It returns the raw YAML bytes, the body string, and an error if validation fails.
func SplitFrontMatter(content string) ([]byte, string, error) {
	const delimiter = "---"

	// Must start with "---"
	if !strings.HasPrefix(content, delimiter) {
		return nil, "", ErrNoFrontMatter
	}

	// Find the end of the first "---" line
	newlineIdx := strings.Index(content, "\n")
	if newlineIdx == -1 {
		return nil, "", ErrNoFrontMatter
	}

	// Find the closing "---"
	closingIdx := strings.Index(content[newlineIdx+1:], "\n"+delimiter)
	if closingIdx == -1 {
		// Could be exactly at the end without newline, but standard is \n---
		return nil, "", ErrNoFrontMatter
	}
	
	closingIdx += newlineIdx + 1

	yamlContent := content[newlineIdx+1 : closingIdx]
	
	// The rest of the content after the closing "---" and its trailing newline
	bodyStart := closingIdx + len("\n"+delimiter)
	var body string
	if bodyStart < len(content) {
		if content[bodyStart] == '\n' {
			bodyStart++ // skip the newline right after ---
		}
		body = content[bodyStart:]
	}

	rawYAML := []byte(yamlContent)

	// Validate required keys
	var parsed map[string]any
	if err := yaml.Unmarshal(rawYAML, &parsed); err != nil {
		return nil, "", err
	}

	if _, ok := parsed["to"]; !ok {
		return nil, "", ErrMissingKey
	}
	if _, ok := parsed["subject"]; !ok {
		return nil, "", ErrMissingKey
	}

	return rawYAML, body, nil
}
