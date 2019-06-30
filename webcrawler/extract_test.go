package webcrawler_test

import (
	"strings"
	"testing"

	"github.com/TomStuart92/web-crawler/webcrawler"
)

func TestExtractLinks(t *testing.T) {
	body := strings.NewReader(`<html><a href="http://test.com"/></html>`)
	page := webcrawler.Page{URL: "http://test.com", Body: body, Links: nil}
	in := make(chan webcrawler.Page)

	links := webcrawler.ExtractLinks(in)
	go func() { in <- page }()

	pageWithLinks := <-links
	if pageWithLinks.URL != page.URL || pageWithLinks.Body != page.Body {
		t.Error("Extract Links has mutated URL/Body")
		return
	}
	if len(pageWithLinks.Links) != 1 {
		t.Errorf("Expected 1 links, got %d", len(pageWithLinks.Links))
	}
	if pageWithLinks.Links[0] != "http://test.com" {
		t.Errorf("Expected http://test.com, got %s", pageWithLinks.Links[0])
	}
}

func TestExtractLinksNoLinks(t *testing.T) {
	body := strings.NewReader(`<html><body/></html>`)
	page := webcrawler.Page{URL: "http://test.com", Body: body, Links: nil}
	in := make(chan webcrawler.Page)

	links := webcrawler.ExtractLinks(in)
	go func() { in <- page }()

	pageWithLinks := <-links
	if pageWithLinks.URL != page.URL || pageWithLinks.Body != page.Body {
		t.Error("Extract Links has mutated URL/Body")
		return
	}
	if len(pageWithLinks.Links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(pageWithLinks.Links))
	}
}
