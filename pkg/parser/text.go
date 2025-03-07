package parser

import (
	"bufio"
	"log/slog"
	"os"
	"strings"

	"github.com/Gmax76/urlcheck/pkg/crawler"
)

type textParser struct {
	filepath string
	results  string
	crawler  crawler.Crawler
}

func NewTextParser(filepath string, crawler crawler.Crawler) Parser {
	return &textParser{
		filepath: filepath,
		results:  "",
		crawler:  crawler,
	}
}

func (p *textParser) Parse() []crawler.CrawlerTarget {
	results := []crawler.CrawlerTarget{}
	target, err := os.Open(p.filepath)
	if err != nil {
		slog.Error("Error while opening file", "error", err)
		os.Exit(1)
	}
	defer target.Close()
	scanner := bufio.NewScanner(target)
	slog.Debug("Parsing file for content to crawl")
	for scanner.Scan() {
		curLine := scanner.Text()
		args := strings.Split(curLine, " ")
		slog.Info("Fetching", "url", args)
		target := crawler.CrawlerTarget{Method: args[0], Url: args[1]}
		p.crawler.Fetch(&target)
		slog.Info("Fetched target", "url", target.Url, "status", target.Status)
		results = append(results, target)
	}
	return results
}
