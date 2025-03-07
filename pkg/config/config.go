package config

import (
	"flag"
	"log/slog"
	"net/http"
	"strings"
)

type config struct {
	CrawlerHeaders http.Header
	ParserTargets  string
	ParserBucket   string
}

type ConfigParams struct {
	TargetsParam string
	BucketParam  string
	HeadersParam string
}

func NewConfig() config {
	cliParams := parseParams()
	headers, err := parseHeaders(cliParams.HeadersParam)
	if err != nil {
		slog.Error("Unable to parse headers")
	}

	return config{
		CrawlerHeaders: headers,
		ParserTargets:  cliParams.TargetsParam,
		ParserBucket:   cliParams.BucketParam,
	}
}

func parseParams() ConfigParams {
	targetFlag := flag.String("target", "", "Target file, can either be local or located in s3")
	headersFlag := flag.String("headers", "", "headers, formatted in a \"Key:Value\" fashion, separated by commas")
	bucketRegionFlag := flag.String("bucketRegion", "", "")
	flag.Parse()

	return ConfigParams{
		TargetsParam: *targetFlag,
		HeadersParam: *headersFlag,
		BucketParam:  *bucketRegionFlag,
	}
}

func parseHeaders(h string) (http.Header, error) {
	headers := http.Header{}
	if h != "" {
		slog.Debug("Parsing headers")
		headersRaw := strings.Split(h, ", ")
		for _, v := range headersRaw {
			hds := strings.SplitN(v, ":", 2)
			slog.Debug("Header discovered", "header", h[0])
			headers.Set(hds[0], hds[1])
		}
	}
	return headers, nil
}
