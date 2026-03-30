package parser

import (
	"errors"
	"testing"
)

func TestSplitFrontMatter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		input := "---\nto: user@example.com\nsubject: hello\n---\nbody text"
		yamlBytes, body, err := SplitFrontMatter(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(yamlBytes) != "to: user@example.com\nsubject: hello" {
			t.Errorf("expected different yaml bytes: %q", yamlBytes)
		}
		if body != "body text" {
			t.Errorf("expected body to be 'body text', got %q", body)
		}
	})

	t.Run("missing to", func(t *testing.T) {
		input := "---\nsubject: hello\n---\nbody text"
		_, _, err := SplitFrontMatter(input)
		if !errors.Is(err, ErrMissingKey) {
			t.Errorf("expected ErrMissingKey, got %v", err)
		}
	})

	t.Run("no front matter", func(t *testing.T) {
		input := "just body text"
		_, _, err := SplitFrontMatter(input)
		if !errors.Is(err, ErrNoFrontMatter) {
			t.Errorf("expected ErrNoFrontMatter, got %v", err)
		}
	})
}

func TestTokenise(t *testing.T) {
	t.Run("simple substitutions", func(t *testing.T) {
		input := "Hello {{NAME}}!"
		tokens, err := Tokenise(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tokens) != 3 {
			t.Fatalf("expected 3 tokens, got %d", len(tokens))
		}
		if tokens[0].Type != TokenText || tokens[0].Raw != "Hello " {
			t.Errorf("token 0 invalid: %+v", tokens[0])
		}
		if tokens[1].Type != TokenVar || tokens[1].Key != "NAME" || tokens[1].Raw != "{{NAME}}" {
			t.Errorf("token 1 invalid: %+v", tokens[1])
		}
		if tokens[2].Type != TokenText || tokens[2].Raw != "!" {
			t.Errorf("token 2 invalid: %+v", tokens[2])
		}
	})

	t.Run("loop", func(t *testing.T) {
		input := "<ul>\n{{#ITEMS}}\n<li>{{ITEM}}</li>\n{{/ITEMS}}\n</ul>"
		tokens, err := Tokenise(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Expect text, loop, text
		if len(tokens) != 3 {
			t.Fatalf("expected 3 tokens, got %d", len(tokens))
		}
		if tokens[1].Type != TokenLoop || tokens[1].Key != "ITEMS" {
			t.Errorf("token 1 not loop items: %+v", tokens[1])
		}
		// inner loop tokens: \n<li>, {{ITEM}}, </li>\n
		inner := tokens[1].Body
		if len(inner) != 3 {
			t.Errorf("expected 3 inner tokens, got %d", len(inner))
		}
	})

	t.Run("unclosed loop", func(t *testing.T) {
		input := "{{#ITEMS}}hello"
		_, err := Tokenise(input)
		if !errors.Is(err, ErrUnclosedTag) {
			t.Errorf("expected ErrUnclosedTag, got %v", err)
		}
	})
}
