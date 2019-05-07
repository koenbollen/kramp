package sources_test

import (
	"context"
	"errors"
	"testing"

	"github.com/koenbollen/kramp/sources"
)

func TestComposite(t *testing.T) {
	t.Parallel()

	a := &single{"a", nil}
	b := &single{"b", nil}

	comp := &sources.Composite{
		Sources: []sources.Source{a, b},
	}

	got, err := comp.Query(context.Background(), "q")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}

	if got[0].Title != "aq" {
		t.Errorf("expected result 0 to be 'aq', got %v", got[0])
	}
	if got[1].Title != "bq" {
		t.Errorf("expected result 1 to be 'bq', got %v", got[1])
	}
}

func TestComposite_Error(t *testing.T) {
	t.Parallel()

	want := errors.New("hi")
	a := &single{"a", want}

	comp := &sources.Composite{
		Sources: []sources.Source{a},
	}

	got, err := comp.Query(context.Background(), "q")
	if err != want {
		t.Errorf("expected an error: %v", err)
	}
	if got != nil {
		t.Errorf("unexpected results: %v", got)
	}
}

type single struct {
	prefix string
	err    error
}

func (s *single) Query(ctx context.Context, input string) ([]sources.Result, error) {
	if s.err != nil {
		return nil, s.err
	}
	return []sources.Result{{Title: s.prefix + input}}, nil
}
