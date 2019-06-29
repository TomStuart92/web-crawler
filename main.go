package main

import (
	"flag"
	"time"

	"github.com/TomStuart92/concurrent-web-crawler/webcrawler"
)

func main() {
	url := flag.String("target", "https://www.monzo.com", "url to scrape")
	maxConcurrency := flag.Int("concurrency", 10, "maximum concurrent web connections")
	timeout := flag.Int("timeout", 5, "timeout in seconds")
	singleDomain := flag.Bool("singleDomain", true, "limit scrape to target base domain")
	flag.Parse()
	t := time.Duration(*timeout)
	wc := webcrawler.New(*maxConcurrency, t*time.Second)
	wc.Scrape(*url, *singleDomain).PrintSiteMap()
}
