package sources_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koenbollen/kramp/sources"
)

const mockGoogleResponse = `
{
  "kind": "books#volumes",
  "totalItems": 516,
  "items": [
    {
      "kind": "books#volume",
      "id": "4lFqCgAAQBAJ",
      "volumeInfo": {
        "title": "De ultieme gids voor Minecraft"
      }
    },
    {
      "kind": "books#volume",
      "id": "0MpyuQEACAAJ",
      "volumeInfo": {
        "title": "Minecraft"
      }
    },
    {
      "kind": "books#volume",
      "id": "GxCXwgEACAAJ",
      "volumeInfo": {
        "title": "Minecraft: The Search for Redstone: Een Onofficieel Minecraft Dungeon Diary",
        "authors": [
          "Vincent Verret"
        ]
      }
    }
  ]
}
`

func TestBooks(t *testing.T) {
	t.Parallel()

	expectedURL := "/?maxResults=5&printType=books&q=Minecraft"
	headCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			headCalled = true
			http.Error(w, "OK", http.StatusOK)
			return
		}
		if r.URL.String() != expectedURL {
			t.Errorf("invalid URL called, want %v got %v", expectedURL, r.URL.String())
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockGoogleResponse)) // nolint: errcheck
	}))
	defer server.Close()

	source := &sources.Books{
		Limit: 5,
		URL:   server.URL,
	}

	result, err := source.Query(context.Background(), "Minecraft")
	if err != nil {
		t.Fatalf("unexpected error from query: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 result, got %d", len(result))
	}

	if result[0].Title != "De ultieme gids voor Minecraft" {
		t.Fatalf("unexpected title #0: want 'De ultieme gids voor Minecraft', got %v", result[0])
	}

	if result[2].Authors[0] != "Vincent Verret" {
		t.Fatalf("unexpected author #2: want 'Vincent Verret', got %v", result[2])
	}

	if !headCalled {
		t.Error("HEAD request not called")
	}
}
