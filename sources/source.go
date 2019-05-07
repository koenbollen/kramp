package sources

import (
	"context"
	"os"
	"strconv"

	"go.uber.org/zap"
)

// ResultKind indicated what type of entity a result entry is.
type ResultKind string

const (
	// KindBook is a book fetched from Google's API.
	KindBook ResultKind = "book"

	// KindAlbum is an album from the iTunes API
	KindAlbum ResultKind = "album"
)

// Result is one object that a source can return multiple of.
type Result struct {
	Title   string     `json:"title"`
	Authors []string   `json:"authors"`
	Kind    ResultKind `json:"kind"`
}

// Source is a queriable source that can returns results on a given
// query.
type Source interface {
	Query(context.Context, string) ([]Result, error)
}

// NewFromEnvironment returns the default composite source
// from with config from the environment.
func NewFromEnvironment(logger *zap.Logger) Source {
	limitBooks := 5
	if l := os.Getenv("LIMIT_BOOKS"); l != "" {
		limitBooks, _ = strconv.Atoi(l)
	}

	limitAlbums := 5
	if l := os.Getenv("LIMIT_ALBUMS"); l != "" {
		limitAlbums, _ = strconv.Atoi(l)
	}

	return &SortedByTitle{
		&Composite{
			[]Source{
				&Safe{
					Source: &Albums{Limit: limitAlbums},
					Logger: logger.With(zap.String("source", "albums")),
				},
				&Safe{
					Source: &Books{Limit: limitBooks},
					Logger: logger.With(zap.String("source", "books")),
				},
			},
		},
	}
}
