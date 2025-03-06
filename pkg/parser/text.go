package parser

import (
	"bufio"
	"log"
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

func (p *textParser) Parse() {
	target, err := os.Open(p.filepath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer target.Close()
	scanner := bufio.NewScanner(target)
	slog.Debug("Parsing file for content to crawl")
	for scanner.Scan() {
		curLine := scanner.Text()
		args := strings.Split(curLine, " ")
		slog.Info("Fetching", "url", args)
		target := p.crawler.Fetch(crawler.CrawlerTarget{Method: args[0], Url: args[1]})
		slog.Info("Fetched target", "url", target.Url, "status", target.Status)
	}
}
