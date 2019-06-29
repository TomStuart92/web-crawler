package webcrawler

import (
	"io"

	"golang.org/x/net/html"
)

func extractLinks(body io.Reader) []string {
	foundLinks := make([]string, 0)
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			return foundLinks
		}

		if tt == html.StartTagToken {
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						link := a.Val
						foundLinks = append(foundLinks, link)
						break
					}
				}
			}
		}
	}
	return foundLinks
}
