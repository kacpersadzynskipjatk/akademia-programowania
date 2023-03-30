package fetcher

import (
	"context"
	"io"
)

type Response struct {
	Data struct {
		Children []struct {
			Data struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditFetcher interface {
	Fetch(context.Context) (Response, error)
	Save(io.Writer) error
}
