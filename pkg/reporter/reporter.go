package reporter

import (
	"log/slog"

	"github.com/Gmax76/urlcheck/pkg/crawler"
)

type (
	Reporter struct {
		results []*crawler.CrawlerTarget
	}
)

func NewReporter() Reporter {
	return Reporter{[]*crawler.CrawlerTarget{}}
}

func (r *Reporter) AppendResult(result *crawler.CrawlerTarget) {
	r.results = append(r.results, result)
}

func (r *Reporter) ProduceReport() {
	slog.Info("Report", "results", r.results)
}
