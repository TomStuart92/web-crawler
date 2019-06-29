// Package webcrawler exports a non-thread safe webcrawler.
package webcrawler

import (
	"io"
	"log"
	"strings"

	"github.com/TomStuart92/web-crawler/graph"
)

// Page represents a single webpage to be scraped
type Page struct {
	URL   string
	Body  io.Reader
	Links []string
}

// WebCrawler is an implementation of a concurrent webcrawler.
// It tracks seen pages in a graph data structure.
type WebCrawler struct {
	maxConcurrency int
	graph          *graph.Graph
	baseURL        string
	outstanding    int
}

// New returns a new webcrawler
func New(maxConcurrency int) *WebCrawler {
	return &WebCrawler{maxConcurrency, graph.NewGraph(), "", 0}
}

func isWithinBaseDomain(url string, baseURL string) bool {
	return strings.HasPrefix(url, baseURL)
}

// GenerateSiteMap handles concurrent generation of a sitemap based on  a given baseURL
func (w *WebCrawler) GenerateSiteMap(baseURL string, singleDomain bool) *WebCrawler {
	w.baseURL = baseURL
	newLinksToScrape := make(chan Page)

	fetchChannel, errChan := FetchURLs(w.maxConcurrency, newLinksToScrape)
	resultChan := extractLinks(fetchChannel)

	go func() {
		for err := range errChan {
			log.Print(err)
		}
	}()

	newLinksToScrape <- Page{baseURL, nil, nil}
	w.graph.AddNode(baseURL)
	w.outstanding = 1

	for w.outstanding != 0 && len(resultChan) == 0 {
		log.Printf("Outstanding Requests To Be Scraped-> %d", w.outstanding-len(fetchChannel)-len(resultChan))
		result := <-resultChan
		w.outstanding--
		for _, link := range result.Links {
			newLink := parseReturnedLink(result.URL, link)

			if singleDomain && !isWithinBaseDomain(newLink, baseURL) {
				log.Printf("%s is outside base domain %s, discarding...", newLink, baseURL)
				continue
			}
			if !w.graph.HasNode(newLink) {
				log.Printf("%s is a new link, scraping... ", newLink)
				w.outstanding++
				w.graph.AddNode(newLink)
				w.graph.AddEdge(result.URL, newLink)
				go func() { newLinksToScrape <- Page{newLink, nil, nil} }()
				continue
			}
			log.Printf("%s has been seen before, discarding... ", newLink)
		}
	}
	log.Printf("Scraping Complete, Returning...")
	return w
}

// PrintSiteMap handles the printing of a site map after it has been discovered
func (w *WebCrawler) PrintSiteMap() {
	log.Print("Printing BFS graph...")
	w.graph.Print(w.baseURL)
}
