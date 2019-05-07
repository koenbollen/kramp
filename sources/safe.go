package sources

import (
	"context"

	"go.uber.org/zap"
)

// Safe is a meta source that wraps around another source and allows
// it to fail silently and continue. Errors are logged by the given
// logger if not nil.
type Safe struct {
	Source Source
	Logger *zap.Logger
}

// Query calls the wrapped source and ignores/logs the potenital error.
func (s *Safe) Query(ctx context.Context, input string) ([]Result, error) {
	result, err := s.Source.Query(ctx, input)
	if err != nil && s.Logger != nil {
		s.Logger.Error("safely failed to execute query", zap.String("input", input), zap.Error(err))
	}
	return result, nil
}
