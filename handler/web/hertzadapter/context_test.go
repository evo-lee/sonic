package hertzadapter

import (
	"context"
	"testing"

	hertzapp "github.com/cloudwego/hertz/pkg/app"
)

func TestContextPersistsRequestContextValuesAcrossWrappers(t *testing.T) {
	var reqCtx hertzapp.RequestContext

	first := NewContext(context.Background(), &reqCtx)
	first.Set("authorized_user", "litang")

	second := NewContext(context.Background(), &reqCtx)
	value := second.RequestContext().Value("authorized_user")
	if value != "litang" {
		t.Fatalf("expected persisted request context value, got %v", value)
	}

	stored, ok := second.Get("authorized_user")
	if !ok || stored != "litang" {
		t.Fatalf("expected context Get to return persisted value, got %v, %v", stored, ok)
	}
}
