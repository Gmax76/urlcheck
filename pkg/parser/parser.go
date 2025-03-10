package parser

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Gmax76/urlcheck/pkg/crawler"
	"github.com/Gmax76/urlcheck/pkg/file"
	"github.com/Gmax76/urlcheck/pkg/reporter"
)

type Parser interface {
	Parse() []crawler.CrawlerTarget
}

type ParserParams struct {
	Target       string
	Filename     string
	BucketRegion string
}

func InitParser(p ParserParams, crawler *crawler.Crawler, reporter *reporter.Reporter) Parser {
	if p.Target == "" {
		slog.Error("Target not defined, please specify a targets file")
		os.Exit(1)
	}
	slog.Info("Specified target file", "target", p.Target)

	if strings.HasPrefix(p.Target, "s3://") {
		slog.Info("File is located in s3, attempting download")
		s3Controller := file.NewS3Controller(p.BucketRegion)
		p.Filename = s3Controller.Get(p.Target)
	} else {
		slog.Info("Assuming file is local")
		p.Filename = p.Target
	}
	// Leaving room here to implement json parser later
	parser := NewTextParser(p.Filename, crawler, reporter)
	return parser
}
