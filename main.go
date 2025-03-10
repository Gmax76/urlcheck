package main

import (
	"log/slog"
	"os"

	"github.com/Gmax76/urlcheck/pkg/config"
	"github.com/Gmax76/urlcheck/pkg/crawler"
	"github.com/Gmax76/urlcheck/pkg/parser"
	"github.com/Gmax76/urlcheck/pkg/reporter"
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	config := config.NewConfig()
	reporter := reporter.NewReporter()
	crawler := crawler.NewCrawler(crawler.CrawlerParams{Headers: config.CrawlerHeaders})
	parser := parser.InitParser(parser.ParserParams{Target: config.ParserTargets, BucketRegion: config.ParserBucket}, &crawler, &reporter)
	parser.Parse()
	reporter.ProduceReport()
}
