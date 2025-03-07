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
			slog.Debug("Header discovered", "header", hds[0])

			headers.Set(hds[0], getTemplatedEnv(hds[1]))
		}
	}
	return headers, nil
}

func getTemplatedEnv(v string) string {
	value := v
	templateRegexp := `^\${(\w)+}$`
	re := regexp.MustCompile(templateRegexp)
	matches := re.FindStringSubmatch(v)
	if len(matches) == 1 {
		value = os.Getenv(matches[0])
		if value == "" {
			slog.Warn("Tried to substitute env var, empty value", "env_var", matches[0])
		}
	}
	return value
}
