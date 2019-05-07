package sources

import (
	"context"
	"sort"
)

// SortedByTitle is a meta source that wraps another source and sorts it's
// results before returning.
type SortedByTitle struct {
	Source Source
}

// Query sorts the wrapped source.
func (s *SortedByTitle) Query(ctx context.Context, input string) ([]Result, error) {
	result, err := s.Source.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Title < result[j].Title
	})

	return result, nil
}
