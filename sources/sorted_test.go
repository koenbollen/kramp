package sources_test

import (
	"context"
	"testing"

	"github.com/koenbollen/kramp/sources"
)

func TestSortedByTitle(t *testing.T) {
	t.Parallel()

	a := &single{"a", nil}
	b := &single{"b", nil}

	comp := &sources.SortedByTitle{
		&sources.Composite{
			Sources: []sources.Source{b, a},
		},
	}

	got, err := comp.Query(context.Background(), "q")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got[0].Title != "aq" {
		t.Errorf("expected result 0 to be 'aq', got %v", got[0])
	}
	if got[1].Title != "bq" {
		t.Errorf("expected result 1 to be 'bq', got %v", got[1])
	}
}
