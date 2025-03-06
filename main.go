package main

import (
	"flag"

	"github.com/Gmax76/urlcheck/pkg/crawler"
	"github.com/Gmax76/urlcheck/pkg/parser"
)

func main() {
	targetFlag := flag.String("target", "", "Target file, can either be local or located in s3")
	headersFlag := flag.String("headers", "", "headers, formatted in a \"Key:Value\" fashion, separated by commas")
	bucketRegionFlag := flag.String("bucketRegion", "", "")
	flag.Parse()

	crawler := crawler.NewCrawler(crawler.CrawlerParams{RawHeaders: *headersFlag})
	parser := parser.InitParser(parser.ParserParams{Target: *targetFlag, BucketRegion: *bucketRegionFlag}, crawler)
	parser.Parse()
}
