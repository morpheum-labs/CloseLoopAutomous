package ai

import (
	"testing"
)

func TestParseIdeationJSON_plainArray(t *testing.T) {
	raw := `[{"title":"A","description":"D","impact":0.5,"feasibility":0.6,"category":"feature"}]`
	got, err := parseIdeationJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Title != "A" || got[0].Description != "D" {
		t.Fatalf("%+v", got)
	}
}

func TestParseIdeationJSON_codeFence(t *testing.T) {
	raw := "```json\n[{\"title\":\"X\",\"description\":\"Y\"}]\n```"
	got, err := parseIdeationJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Title != "X" {
		t.Fatalf("%+v", got)
	}
}

func TestParseIdeationJSON_skipsEmptyRows(t *testing.T) {
	raw := `[{"title":"","description":""},{"title":"ok","description":"body"}]`
	got, err := parseIdeationJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Title != "ok" {
		t.Fatalf("%+v", got)
	}
}

func TestParseIdeationJSON_emptyArrayErrors(t *testing.T) {
	_, err := parseIdeationJSON(`[]`)
	if err == nil {
		t.Fatal("expected error")
	}
}
