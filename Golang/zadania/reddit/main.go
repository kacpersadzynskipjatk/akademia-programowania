package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reddit/fetcher"
	"sync"
	"time"
)

func main() {
	var f fetcher.RedditFetcher // do not change
	var w io.Writer             // do not change

	//Solution to 429 error https://www.reddit.com/r/redditdev/comments/t8e8hc/getting_nothing_but_429_responses_when_using_go/
	subreddits := []string{"golang", "docker", "kubernetes", "aws", "googlecloud"}
	for _, subreddit := range subreddits {
		wg.Add(1)
		go run(f, w, "https://www.reddit.com/r/"+subreddit+".json", "reddit_output_"+subreddit+".txt")
	}
	wg.Wait()
}

var wg sync.WaitGroup

func run(f fetcher.RedditFetcher, w io.Writer, hostUrl, fileName string) {
	defer wg.Done()
	f = NewFetcher(hostUrl, time.Second*3)
	_, err := f.Fetch(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w = file
	err = f.Save(w)
	if err != nil {
		log.Fatal(err)
	}
}

type Fetcher struct {
	c      *http.Client
	host   string
	output fetcher.Response
}

func NewFetcher(host string, t time.Duration) *Fetcher {
	return &Fetcher{
		host: host,
		c: &http.Client{
			Timeout: t,
			Transport: &http.Transport{
				TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
			},
		},
	}
}

func (e *Fetcher) Fetch(ctx context.Context) (fetcher.Response, error) {
	ctx = context.WithValue(ctx, "requestID", time.Now().Unix())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.host, http.NoBody)
	if err != nil {
		return fetcher.Response{}, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("User-Agent", "Custom Agent")
	resp, err := e.c.Do(req)
	if err != nil {
		return fetcher.Response{}, fmt.Errorf("cannot get data: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fetcher.Response{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var data fetcher.Response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fetcher.Response{}, fmt.Errorf("cannot unmarshal data: %w", err)
	}
	e.output = data
	return data, nil
}
func (e *Fetcher) Save(w io.Writer) error {
	//CREATE A TXT FILE
	for _, child := range e.output.Data.Children {
		d := fmt.Sprintf("%s\n%s\n\n", child.Data.Title, child.Data.URL)
		_, err := w.Write([]byte(d))
		if err != nil {
			return err
		}
	}
	return nil
}
