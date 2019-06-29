package webcrawler

import (
	"strings"

	"github.com/TomStuart92/web-crawler/parser"
)

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

func extractLinks(in chan Page) chan Page {
	resultsChan := make(chan Page)
	go func() {
		for page := range in {
			go func(p Page) {
				if p.Body == nil {
					return
				}
				p.Links = parser.ExtractLinks(p.Body)
				resultsChan <- p
			}(page)
		}
	}()
	return resultsChan
}
