package parser

import (
	"errors"
	"strings"
)

var (
	ErrUnclosedTag        = errors.New("unclosed block tag")
	ErrNestedNotSupported = errors.New("nested block tags are not supported")
)

type TokenType int

const (
	TokenText TokenType = iota
	TokenVar
	TokenLoop
	TokenCond
)

// Token represents a parsed chunk of the template.
type Token struct {
	Type TokenType
	Key  string
	Body []Token // For Loop and Cond
	Raw  string  // For Text, or the original tag string if needed
}

// Tokenise parses the body text into a list of Tokens.
func Tokenise(text string) ([]Token, error) {
	return parseBlock(text)
}

func parseBlock(text string) ([]Token, error) {
	var tokens []Token

	for len(text) > 0 {
		startIdx := strings.Index(text, "{{")
		if startIdx == -1 {
			tokens = append(tokens, Token{Type: TokenText, Raw: text})
			break
		}

		if startIdx > 0 {
			tokens = append(tokens, Token{Type: TokenText, Raw: text[:startIdx]})
		}

		text = text[startIdx:]
		endIdx := strings.Index(text, "}}")
		if endIdx == -1 {
			// Unclosed {{, treat as text
			tokens = append(tokens, Token{Type: TokenText, Raw: text})
			break
		}

		tagContent := strings.TrimSpace(text[2:endIdx])
		fullTag := text[:endIdx+2]

		if strings.HasPrefix(tagContent, "#") {
			// Block tag: {{#LOOP}} or {{#IF_X}}
			key := tagContent[1:]
			isCond := strings.HasPrefix(key, "IF_")

			closeTag := "{{/" + key + "}}"
			closeIdx := strings.Index(text, closeTag)
			if closeIdx == -1 {
				return nil, ErrUnclosedTag
			}

			innerBody := text[len(fullTag):closeIdx]
			
			// Nested tags?
			if strings.Contains(innerBody, "{{#") {
				return nil, ErrNestedNotSupported
			}

			innerTokens, err := parseBlock(innerBody)
			if err != nil {
				return nil, err
			}

			tok := Token{
				Key:  key,
				Body: innerTokens,
				Raw:  fullTag,
			}
			if isCond {
				tok.Type = TokenCond
				tok.Key = key[3:] // strip IF_
			} else {
				tok.Type = TokenLoop
				tok.Key = key
			}
			tokens = append(tokens, tok)

			text = text[closeIdx+len(closeTag):]
		} else if strings.HasPrefix(tagContent, "/") {
			// Closing tag outside of an open block shouldn't happen unless malformed,
			// or treated as text
			tokens = append(tokens, Token{Type: TokenText, Raw: fullTag})
			text = text[endIdx+2:]
		} else {
			// Simple variable: {{VAR}}
			tokens = append(tokens, Token{Type: TokenVar, Key: tagContent, Raw: fullTag})
			text = text[endIdx+2:]
		}
	}

	return tokens, nil
}
