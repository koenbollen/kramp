package sources

import "context"

// Composite is a meta sources that combines multiple sources into one
// result set. All sources are query in parallel and are return in order
// of the sources themselves. If one query fails that error is returned
// directly.
type Composite struct {
	Sources []Source
}

type indexed struct {
	index  int
	result []Result
}

// Query calls all wrapped sources in parallel and returns one result set.
func (c *Composite) Query(ctx context.Context, input string) ([]Result, error) {
	l := len(c.Sources)
	resultChannel := make(chan indexed, l)
	errChannel := make(chan error, l)

	for i, s := range c.Sources {
		index := i
		source := s
		go func() {
			var err error
			result, err := source.Query(ctx, input)
			if err != nil {
				errChannel <- err
			}
			resultChannel <- indexed{index, result}
		}()
	}

	results := make([][]Result, l)
	for {
		select {
		case r := <-resultChannel:
			results[r.index] = r.result
			l--
		case e := <-errChannel:
			return nil, e
		}
		if l == 0 {
			break
		}
	}

	var result []Result
	for _, r := range results {
		result = append(result, r...)
	}
	return result, nil
}
