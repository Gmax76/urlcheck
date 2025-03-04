package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func main() {
	headers := http.Header{}
	filename := ""
	targetFlag := flag.String("target", "", "Target file, can either be local or located in s3")
	headersFlag := flag.String("headers", "", "headers, formatted in a \"Key:Value\" fashion, separated by commas")
	bucketRegionFlag := flag.String("bucketRegion", "", "")
	flag.Parse()

	if *targetFlag == "" {
		log.Fatal("Target not defined, please specify a targets file")
	}
	log.Printf("Specified target file: %v", *targetFlag)

	if *headersFlag != "" {
		log.Print("Parsing headers")
		headersRaw := strings.Split(*headersFlag, ", ")
		for _, v := range headersRaw {
			h := strings.SplitN(v, ":", 2)
			log.Printf("Header %v  Value %v", h[0], h[1])
			headers.Set(h[0], h[1])
		}
	}

	if strings.HasPrefix(*targetFlag, "s3://") {
		log.Print("File is located in s3, attempting download")
		ctx := context.Background()
		cfg, _ := config.LoadDefaultConfig(ctx)
		if *bucketRegionFlag != "" {
			cfg.Region = *bucketRegionFlag
		}
		s3Client := s3.NewFromConfig(cfg)
		u, _ := url.Parse(*targetFlag)
		result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(u.Host),
			Key:    aws.String(u.Path[1:]),
		})
		if err != nil {
			log.Fatalf("Error fetching file from s3: %v", err)
		}
		f, err := os.CreateTemp("/tmp", "target")
		if err != nil {
			log.Fatalf("failed to create temporary file for s3 target, %v", err)
		}
		body, err := io.ReadAll(result.Body)
		if err != nil {
			log.Printf("Couldn't read object body from %v. Here's why: %v\n", u.Path[1:], err)
		}
		_, err = f.Write(body)
		if err != nil {
			log.Fatalf("Could not write to temporary file. Error: %v", err)
		}
		f.Close()
		filename = f.Name()
	} else {
		log.Print("Assuming file is local")
		filename = *targetFlag
	}

	target, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer target.Close()
	scanner := bufio.NewScanner(target)
	log.Print("Parsing file for content to crawl")
	client := &http.Client{
		CheckRedirect: defaultCheckRedirect,
	}
	for scanner.Scan() {
		curLine := scanner.Text()
		args := strings.Split(curLine, " ")
		req, err := http.NewRequest(args[0], args[1], nil)
		if err != nil {
			log.Printf("Error forging request: %v", err)
		}
		req.Header = headers
		resp, _ := client.Do(req)
		fmt.Printf("URL: %v, STATUS: %v \n", args[1], resp.Status)
	}
}
