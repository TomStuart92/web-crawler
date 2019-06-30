// Package parser package provides functionality to parse a http response as html,
// identify all link elements (<a/>), and return an array of their href values
package parser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinksFromHTML takes a io.Reader containing a html string
// and iterates through it to find all <a elements>, returning
// an array of all the href attributes from the links.
func ExtractLinksFromHTML(body io.Reader) []string {
	foundLinks := make([]string, 0)

	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		// we've reached end of input/an error has occured
		if tt == html.ErrorToken {
			return foundLinks
		}

		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			t := z.Token()
			// check if token is link
			if t.Data == "a" {

				// if it is, itterate attributes to find ref
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
}

// ParseReturnedLink checks if a returned link is a relative path,
// and returns the absolute link if so without any trailing slashes.
func ParseReturnedLink(base string, link string) string {
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
