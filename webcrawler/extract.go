// Package webcrawler exports a non-thread safe webcrawler.
package webcrawler

import "github.com/TomStuart92/web-crawler/parser"

// ExtractLinks reads pages from a channel and spins out new go routines
// to handle the tokenisation and extraction of <a /> links from the html
func ExtractLinks(in chan Page) chan Page {
	resultsChan := make(chan Page)
	go func() {
		for page := range in {
			go func(p Page) {
				if p.Body == nil {
					return
				}
				p.Links = parser.ExtractLinksFromHTML(p.Body)
				resultsChan <- p
			}(page)
		}
	}()
	return resultsChan
}
