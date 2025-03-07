package main

import (
	"github.com/Gmax76/urlcheck/pkg/config"
	"github.com/Gmax76/urlcheck/pkg/crawler"
	"github.com/Gmax76/urlcheck/pkg/parser"
)

func main() {
	config := config.NewConfig()
	crawler := crawler.NewCrawler(crawler.CrawlerParams{Headers: config.CrawlerHeaders})
	parser := parser.InitParser(parser.ParserParams{Target: config.ParserTargets, BucketRegion: config.ParserBucket}, crawler)
	parser.Parse()
}
