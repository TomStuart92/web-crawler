package parser_test

import (
	"strings"
	"testing"

	"github.com/TomStuart92/web-crawler/parser"
)

func TestExtractLinksFromHTMLEmptyReader(t *testing.T) {
	reader := strings.NewReader("")
	links := parser.ExtractLinksFromHTML(reader)
	if len(links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(links))
	}
}

func TestExtractLinksFromHTMLInvalidHTML(t *testing.T) {
	reader := strings.NewReader("<html><body")
	links := parser.ExtractLinksFromHTML(reader)
	if len(links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(links))
	}
}

func TestExtractLinksFromHTMLLinkWithNoHREF(t *testing.T) {
	reader := strings.NewReader("<html><a/></html>")
	links := parser.ExtractLinksFromHTML(reader)
	if len(links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(links))
	}
}

func TestExtractLinksFromHTMLLinkWithHREF(t *testing.T) {
	reader := strings.NewReader(`<html><a href="http://test.com"></a></html>`)
	links := parser.ExtractLinksFromHTML(reader)
	if len(links) != 1 {
		t.Errorf("Expected 1 links, got %d", len(links))
		return
	}
	if links[0] != "http://test.com" {
		t.Errorf("Expected http://test.com, got %s", links[0])
	}
}

func TestExtractLinksFromHTMLLinkWithSelfClosingHREF(t *testing.T) {
	reader := strings.NewReader(`<html><a href="http://test.com"/></html>`)
	links := parser.ExtractLinksFromHTML(reader)
	if len(links) != 1 {
		t.Errorf("Expected 1 links, got %d", len(links))
		return
	}
	if links[0] != "http://test.com" {
		t.Errorf("Expected http://test.com, got %s", links[0])
	}
}

func TestExtractLinksFromHTMLLinkWithNestedLink(t *testing.T) {
	reader := strings.NewReader(`<html><body><h1><a href="http://test.com"></a></h1></body></html>`)
	links := parser.ExtractLinksFromHTML(reader)
	if len(links) != 1 {
		t.Errorf("Expected 1 links, got %d", len(links))
		return
	}
	if links[0] != "http://test.com" {
		t.Errorf("Expected http://test.com, got %s", links[0])
	}
}

func TestNewAbsoluteURLNoTrailing(t *testing.T) {
	base := "http://test.com"
	link := "http://otherdomain.com"
	parsed := parser.ParseReturnedLink(base, link)
	if parsed != link {
		t.Errorf("expected %s, got %s", link, parsed)
	}
}

func TestNewAbsoluteURLTrailingSlash(t *testing.T) {
	base := "http://test.com"
	link := "http://otherdomain.com/"
	expected := "http://otherdomain.com"
	parsed := parser.ParseReturnedLink(base, link)
	if parsed != expected {
		t.Errorf("expected %s, got %s", expected, parsed)
	}
}

func TestNewRelativeURLNoTrailingSlash(t *testing.T) {
	base := "http://test.com"
	link := "/otherpage"
	expected := "http://test.com/otherpage"
	parsed := parser.ParseReturnedLink(base, link)
	if parsed != expected {
		t.Errorf("expected %s, got %s", expected, parsed)
	}
}

func TestNewRelativeURLTrailingSlash(t *testing.T) {
	base := "http://test.com"
	link := "/otherpage/"
	expected := "http://test.com/otherpage"
	parsed := parser.ParseReturnedLink(base, link)
	if parsed != expected {
		t.Errorf("expected %s, got %s", expected, parsed)
	}
}
