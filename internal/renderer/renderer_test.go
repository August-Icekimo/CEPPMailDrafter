package renderer

import (
	"strings"
	"testing"

	"github.com/icekimo/CEPPMailDrafter/internal/parser"
)

type mockDS struct {
	kv    map[string]string
	lists map[string][]string
}

func (m mockDS) Get(k string) (string, bool) {
	v, ok := m.kv[strings.ToUpper(k)]
	return v, ok
}

func (m mockDS) GetList(k string) ([]string, bool) {
	v, ok := m.lists[strings.ToUpper(k)]
	return v, ok
}

func TestRenderer(t *testing.T) {
	ds := mockDS{
		kv: map[string]string{
			"MONTH":      "三月",
			"ATTACHMENT": "true",
			"MISSING":    "", // will not be found in Get because empty value vs no key ok?
		},
		lists: map[string][]string{
			"ITEMS": {"item1", "item2"},
		},
	}
	delete(ds.kv, "MISSING") // remove entirely to simulate false ok

	r := New(ds)

	frontYAML := []byte("to: user@test.com\nsubject: 測試 {{MONTH}}\ncc:\n  - cc1@t.com\n  - cc2@t.com")
	tokens, err := parser.Tokenise("month: {{MONTH}}\n{{#ITEMS}}* {{ITEM}}\n{{/ITEMS}}{{#IF_ATTACHMENT}}has att{{/IF_ATTACHMENT}} {{MISSING}}")
	if err != nil {
		t.Fatalf("parser err: %v", err)
	}

	msg, err := r.Render(frontYAML, tokens)
	if err != nil {
		t.Fatalf("render err: %v", err)
	}

	if msg.Subject != "測試 三月" {
		t.Errorf("subject unexpanded: %s", msg.Subject)
	}
	if len(msg.To) != 1 || msg.To[0] != "user@test.com" {
		t.Errorf("to parsed wrong: %v", msg.To)
	}
	if len(msg.Cc) != 2 || msg.Cc[0] != "cc1@t.com" {
		t.Errorf("cc parsed wrong: %v", msg.Cc)
	}

	expectedBody := "month: 三月\n* item1\n* item2\nhas att {{MISSING}}"
	if msg.Body != expectedBody {
		t.Errorf("body wrong. expected:\n%s\ngot:\n%s", expectedBody, msg.Body)
	}
}
