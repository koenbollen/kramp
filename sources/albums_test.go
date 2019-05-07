package sources_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koenbollen/kramp/sources"
)

const mockItunesResponse = `
{
  "resultCount": 4,
  "results": [
    {
      "collectionType": "Album",
      "artistName": "Minecraft",
      "collectionName": "Gameover - Single"
    },
    {
      "collectionType": "Album",
      "artistName": "Minecraft",
      "collectionName": "Take Me Down - Single"
    },
    {
      "collectionType": "Album",
      "artistName": "Minecraft",
      "collectionName": "The Remixes - Single"
    },
    {
      "collectionType": "Album",
      "artistName": "Minecraft",
      "collectionName": "The Best of Minecraft"
    }
  ]
}
`

func TestAlbums(t *testing.T) {
	t.Parallel()

	expectedURL := "/?entity=album&limit=5&term=Minecraft"
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
		w.Write([]byte(mockItunesResponse)) // nolint: errcheck
	}))
	defer server.Close()

	source := &sources.Albums{
		Limit: 5,
		URL:   server.URL,
	}

	result, err := source.Query(context.Background(), "Minecraft")
	if err != nil {
		t.Fatalf("unexpected error from query: %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("expected 4 result, got %d", len(result))
	}

	if result[0].Title != "Gameover - Single" {
		t.Fatalf("unexpected title #0: want 'Gameover - Single', got %v", result[0])
	}

	if !headCalled {
		t.Error("HEAD request not called")
	}
}
