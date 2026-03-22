package tfidftags

import "testing"

func TestSuggest_TFIDF(t *testing.T) {
	corpus := []string{
		"payment checkout stripe billing",
		"payment refund billing support",
		"login oauth session security",
	}
	target := "improve payment checkout flow for billing"
	got := Suggest(corpus, target, 10, 2)
	if len(got) == 0 {
		t.Fatal("expected suggestions")
	}
	found := false
	for _, g := range got {
		if g.Token == "checkout" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected checkout in results, got %#v", got)
	}
}

func TestSuggest_FreqFallback(t *testing.T) {
	got := Suggest(nil, "alpha beta alpha gamma", 5, 2)
	if len(got) < 2 {
		t.Fatalf("got %#v", got)
	}
	if got[0].Token != "alpha" {
		t.Fatalf("want alpha first by count, got %#v", got)
	}
}

func TestSuggest_StopwordsFiltered(t *testing.T) {
	got := Suggest(nil, "the quick brown fox", 10, 2)
	for _, g := range got {
		if g.Token == "the" {
			t.Fatalf("stopword leaked: %#v", got)
		}
	}
}

func TestTokenize(t *testing.T) {
	got := Tokenize("Hello, world! 123")
	if len(got) != 3 || got[0] != "hello" || got[1] != "world" || got[2] != "123" {
		t.Fatalf("got %q", got)
	}
}
