package config

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"regexp"
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
	config := config{}
	cliParams := parseParams()
	config.ParserBucket = cliParams.BucketParam
	config.ParserTargets = cliParams.TargetsParam

	headers, err := config.parseHeaders(cliParams.HeadersParam)
	if err != nil {
		slog.Error("Unable to parse headers", "error", err)
	}
	config.CrawlerHeaders = headers
	return config
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

func (c *config) parseHeaders(h string) (http.Header, error) {
	headers := http.Header{}
	headers.Set("User-Agent", "Gmax76/urlcheck")
	if h != "" {
		headersRaw := strings.Split(h, ", ")
		for _, v := range headersRaw {
			hds := strings.SplitN(v, ":", 2)
			headers.Set(hds[0], getTemplatedEnv(strings.TrimSpace(hds[1])))
			slog.Debug("Setting header", "header", hds)
		}
	}
	return headers, nil
}

func getTemplatedEnv(v string) string {
	value := v
	templateRegexp := `^{{(\w+)}}$`
	re := regexp.MustCompile(templateRegexp)
	matches := re.FindStringSubmatch(v)
	if len(matches) == 2 {
		val, found := os.LookupEnv(matches[1])
		if !found {
			slog.Warn("Tried to substitute env var, unset value", "env_var", matches[1])
		}
		value = val
	}
	return value
}
