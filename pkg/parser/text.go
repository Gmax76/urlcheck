package parser

import (
	"bufio"
	"log/slog"
	"os"
	"strings"

	"github.com/Gmax76/urlcheck/pkg/crawler"
	"github.com/Gmax76/urlcheck/pkg/reporter"
)

type textParser struct {
	filepath string
	results  string
	crawler  *crawler.Crawler
	reporter *reporter.Reporter
}

func NewTextParser(filepath string, crawler *crawler.Crawler, reporter *reporter.Reporter) Parser {
	return &textParser{
		filepath: filepath,
		results:  "",
		crawler:  crawler,
		reporter: reporter,
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
		crawlTarget := crawler.CrawlerTarget{Method: args[0], Url: args[1]}
		p.crawler.Fetch(&crawlTarget)
		p.reporter.AppendResult(&crawlTarget)
	}
	return results
}
