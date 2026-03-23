package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatClient_ChatCompletion_ok(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("path %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"role": "assistant", "content": " hello "}},
			},
		})
	}))
	defer srv.Close()

	c := &ChatClient{BaseURL: srv.URL, HTTP: srv.Client()}
	out, err := c.ChatCompletion(context.Background(), "m", "sys", "user", 0.2, 100)
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello" {
		t.Fatalf("got %q", out)
	}
}

func TestChatClient_ChatCompletion_apiError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"message":"bad model"}}`))
	}))
	defer srv.Close()

	c := &ChatClient{BaseURL: srv.URL, HTTP: srv.Client()}
	_, err := c.ChatCompletion(context.Background(), "m", "s", "u", 0, 10)
	if err == nil {
		t.Fatal("expected error")
	}
}
