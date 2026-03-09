package binding

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type bindFormPayload struct {
	Page int      `form:"page"`
	Tags []string `form:"tags"`
}

func TestBindForm(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?page=3&tags=go&tags=web", nil)

	var payload bindFormPayload
	if err := BindForm(req, &payload); err != nil {
		t.Fatalf("BindForm returned error: %v", err)
	}

	if payload.Page != 3 {
		t.Fatalf("expected page=3, got %d", payload.Page)
	}
	if len(payload.Tags) != 2 || payload.Tags[0] != "go" || payload.Tags[1] != "web" {
		t.Fatalf("unexpected tags: %#v", payload.Tags)
	}
}

func TestBindFormPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.PostForm = map[string][]string{
		"page": {"5"},
		"tags": {"api", "binding"},
	}

	var payload bindFormPayload
	if err := BindFormPost(req, &payload); err != nil {
		t.Fatalf("BindFormPost returned error: %v", err)
	}

	if payload.Page != 5 {
		t.Fatalf("expected page=5, got %d", payload.Page)
	}
	if len(payload.Tags) != 2 || payload.Tags[0] != "api" || payload.Tags[1] != "binding" {
		t.Fatalf("unexpected tags: %#v", payload.Tags)
	}
}
