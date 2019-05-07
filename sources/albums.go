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

const itunesURL = `https://itunes.apple.com/search`

var (
	itunes_call = expvar.NewInt("last_itunes_call")
)

// Albums is the source that fetches albums from the iTunes API
type Albums struct {
	URL   string
	Limit int

	client *warmed.Client
}

// Query returns album results based on the given input
func (a *Albums) Query(ctx context.Context, input string) ([]Result, error) {
	if a.URL == "" {
		a.URL = itunesURL
	}
	if a.client == nil {
		a.client = warmed.New(a.URL)
	}

	params := url.Values{}
	params.Set("term", input)
	params.Set("entity", "album")
	params.Set("limit", strconv.Itoa(a.Limit))
	req, _ := http.NewRequest(http.MethodGet, a.URL+"?"+params.Encode(), nil)
	req = req.WithContext(ctx)

	start := time.Now()
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request albums from itunes")
	}
	took := time.Since(start)
	itunes_call.Set(int64(took))

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status from itunes: %v", resp.Status)
	}

	body := struct {
		Results []struct {
			ArtistName     string
			CollectionName string
		}
	}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse json from itunes")
	}

	var result []Result
	for _, i := range body.Results {
		result = append(result, Result{
			Title:   i.CollectionName,
			Authors: []string{i.ArtistName},
			Kind:    KindAlbum,
		})
	}

	return result, nil
}
