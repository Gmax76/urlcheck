package crawler

import (
	"log"
	"log/slog"
	"net/http"
	"strings"
)

type Crawler interface {
	Fetch(CrawlerTarget) *CrawlerTarget
}

type crawler struct {
	headers http.Header
	client  *http.Client
}

type CrawlerTarget struct {
	Method string
	Url    string
	Status string
}

type CrawlerParams struct {
	RawHeaders string
}

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func NewCrawler(p CrawlerParams) Crawler {
	headers := http.Header{}
	client := &http.Client{
		CheckRedirect: defaultCheckRedirect,
	}

	if p.RawHeaders != "" {
		slog.Debug("Parsing headers")
		headersRaw := strings.Split(p.RawHeaders, ", ")
		for _, v := range headersRaw {
			h := strings.SplitN(v, ":", 2)
			slog.Debug("Header discovered", "header", h[0])
			headers.Set(h[0], h[1])
		}
	}

	return &crawler{
		client:  client,
		headers: headers,
	}
}

func (c *crawler) Fetch(t CrawlerTarget) *CrawlerTarget {
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
	return &t
}
