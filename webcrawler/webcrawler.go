package webcrawler

import (
	"log"
	"strings"
	"time"

	"github.com/TomStuart92/concurrent-web-crawler/graph"
)

// WebCrawler is an implementation of a concurrent webcrawler.
// It tracks seen pages in a graph data structure.
type WebCrawler struct {
	maxConcurrency int
	timeout        time.Duration
	pool           *Pool
	graph          *graph.Graph
	baseURL        string
	running        bool
}

func New(maxConcurrency int, timeout time.Duration) *WebCrawler {
	pool := InitializePool(maxConcurrency)
	pool.SetWorkFn(getPage)
	pool.SetTransformFn(extractLinks)
	return &WebCrawler{maxConcurrency, timeout, pool, graph.NewGraph(), "", false}
}

func isWithinBaseDomain(url string, baseURL string) bool {
	return strings.HasPrefix(url, baseURL)
}

func parseReturnedLink(base string, link string) string {
	var newLink string
	if strings.HasPrefix(link, "/") {
		newLink = base + link
	} else {
		newLink = link
	}

	if strings.HasSuffix(newLink, "/") {
		newLink = newLink[0 : len(newLink)-1]
	}
	return newLink
}

func (w *WebCrawler) Scrape(baseURL string, singleDomain bool) *WebCrawler {
	w.baseURL = baseURL
	workQueue := make(chan string)

	resultChan, errChan := w.pool.Process(workQueue)

	go func() {
		for err := range errChan {
			log.Print(err)
		}
	}()

	workQueue <- baseURL
	w.graph.AddNode(baseURL)
	w.running = true

	for w.running {
		select {
		case result := <-resultChan:
			for _, link := range result.Output {
				newLink := parseReturnedLink(result.Input, link)

				if singleDomain && !isWithinBaseDomain(newLink, baseURL) {
					log.Printf("%s is outside base domain %s, discarding...", newLink, baseURL)
					continue
				}
				if !w.graph.HasNode(newLink) {
					log.Printf("%s is a new link, scraping... ", newLink)
					w.graph.AddNode(newLink)
					w.graph.AddEdge(result.Input, newLink)
					go func() { workQueue <- newLink }()
					continue
				}
				log.Printf("%s has been seen before, discarding... ", newLink)
			}
		case <-time.After(w.timeout):
			log.Printf("Timeout Reached")
			if len(workQueue) == 0 {
				log.Printf("No New Urls To Scrape, Terminating...")
				w.Stop()
			}
		}
	}
	return w
}

func (w *WebCrawler) Stop() {
	log.Printf("WebCrawler Terminating...")
	w.pool.Stop()
	w.running = false
}

func (w *WebCrawler) PrintSiteMap() {
	w.graph.Print(w.baseURL)
}
