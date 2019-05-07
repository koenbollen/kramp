package sources

import (
	"context"
	"encoding/json"
	"expvar"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/koenbollen/warmed"
	"github.com/pkg/errors"
)

const googleURL = `https://www.googleapis.com/books/v1/volumes`

var (
	google_call = expvar.NewInt("last_google_call")
)

// Books is the source that fetches books search on the Google Books API
type Books struct {
	URL   string
	Limit int

	client *warmed.Client
}

// Query returns a list of books based on the input given.
func (b *Books) Query(ctx context.Context, input string) ([]Result, error) {
	if b.URL == "" {
		b.URL = googleURL
	}
	if b.client == nil {
		b.client = warmed.New(b.URL)
	}

	params := url.Values{}
	params.Set("q", input)
	params.Set("printType", "books")
	params.Set("maxResults", strconv.Itoa(b.Limit))
	req, _ := http.NewRequest(http.MethodGet, b.URL+"?"+params.Encode(), nil)
	req = req.WithContext(ctx)

	start := time.Now()
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request books from google")
	}
	took := time.Since(start)
	google_call.Set(int64(took))

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status from google: %v", resp.Status)
	}

	body := struct {
		TotalItems int
		Items      []struct {
			VolumeInfo struct {
				Title   string
				Authors []string
			}
		}
	}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse json from google")
	}

	var result []Result
	for _, i := range body.Items {
		result = append(result, Result{
			Title:   i.VolumeInfo.Title,
			Authors: i.VolumeInfo.Authors,
			Kind:    KindBook,
		})
	}

	return result, nil
}
