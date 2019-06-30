// Package webcrawler exports a non-thread safe webcrawler.
package webcrawler

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/TomStuart92/web-crawler/graph"
	"github.com/TomStuart92/web-crawler/parser"
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
	Client         HTTPClient // export client to allow custom transport
	baseURL        string
	outstanding    int
}

// HTTPClient provides a simple abstraction
// to allow dependency injection.
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// New returns a new webcrawler
func New(maxConcurrency int) *WebCrawler {
	return &WebCrawler{
		maxConcurrency: maxConcurrency,
		graph:          graph.NewGraph(),
		Client:         &http.Client{},
	}
}

func isWithinBaseDomain(url string, baseURL string) bool {
	return strings.HasPrefix(url, baseURL)
}

// GenerateSiteMap handles concurrent generation of a sitemap based on  a given baseURL
func (w *WebCrawler) GenerateSiteMap(baseURL string, singleDomain bool) {
	w.baseURL = baseURL
	newLinksToScrape := make(chan Page)

	// Set up channels to fetch data
	fetchChannel, errChan := FetchURLs(w.Client, w.maxConcurrency, newLinksToScrape)

	// Set up go routine to parse html bodies
	resultChan := ExtractLinks(fetchChannel)

	go func() {
		for err := range errChan {
			log.Print(err)
		}
	}()

	newLinksToScrape <- Page{baseURL, nil, nil}
	w.graph.AddNode(baseURL)
	w.outstanding = 1

	// while there are outstanding pages to scape/we have results to process
	for w.outstanding > 0 && len(resultChan) == 0 {
		log.Printf("Outstanding Requests To Be Scraped-> %d", w.outstanding)
		result := <-resultChan
		w.outstanding--
		for _, link := range result.Links {
			newLink := parser.ParseReturnedLink(result.URL, link)

			// check if new link is within the same domain
			if singleDomain && !isWithinBaseDomain(newLink, baseURL) {
				log.Printf("%s is outside base domain %s, discarding...", newLink, baseURL)
				continue
			}

			// check if we've seen this page before, if not add nodes/edges
			// and schedule for scrape
			if !w.graph.HasNode(newLink) {
				log.Printf("%s is a new link, scraping... ", newLink)
				w.outstanding++
				w.graph.AddNode(newLink)
				w.graph.AddEdge(result.URL, newLink)
				go func() { newLinksToScrape <- Page{newLink, nil, nil} }()
				continue
			}

			// if we have seen it before, we just need to add edge between pages
			w.graph.AddEdge(result.URL, newLink)
			log.Printf("%s has been seen before, discarding... ", newLink)
		}
	}
	log.Printf("Scraping Complete, Returning...")
}

// PrintSiteMap handles the printing of a site map after it has been discovered
func (w *WebCrawler) PrintSiteMap(output io.Writer) {
	w.graph.Print(output, w.baseURL)
}
