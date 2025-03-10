package crawler

import (
	"log"
	"log/slog"
	"net/http"
)

type Crawler struct {
	headers http.Header
	client  *http.Client
}

type CrawlerTarget struct {
	Method string
	Url    string
	Status string
}

type CrawlerParams struct {
	Headers http.Header
}

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func NewCrawler(p CrawlerParams) Crawler {
	client := &http.Client{
		CheckRedirect: defaultCheckRedirect,
	}

	return Crawler{
		client:  client,
		headers: p.Headers,
	}
}

func (c *Crawler) Fetch(t *CrawlerTarget) {
	req, err := http.NewRequest(t.Method, t.Url, nil)
	if err != nil {
		log.Printf("Error forging request: %v", err)
	}
	req.Header = c.headers
	resp, err := c.client.Do(req)
	if err != nil {
		slog.Error("Error fetching target", "target", t, "error", err)
	}
	t.Status = resp.Status
}
