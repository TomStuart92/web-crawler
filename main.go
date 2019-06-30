package main

import (
	"flag"
	"os"

	"github.com/TomStuart92/web-crawler/webcrawler"
)

func main() {
	url := flag.String("target", "https://www.jigsaw.xyz", "url to scrape")
	maxConcurrency := flag.Int("concurrency", 10, "maximum concurrent web connections")
	singleDomain := flag.Bool("singleDomain", true, "limit scrape to target base domain")
	flag.Parse()
	wc := webcrawler.New(*maxConcurrency)
	wc.GenerateSiteMap(*url, *singleDomain)
	wc.PrintSiteMap(os.Stdout)
}
