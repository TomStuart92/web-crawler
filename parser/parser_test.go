package parser_test

import (
	"strings"
	"testing"

	"github.com/TomStuart92/web-crawler/parser"
)

func TestExtractLinksEmptyReader(t *testing.T) {
	reader := strings.NewReader("")
	links := parser.ExtractLinks(reader)
	if len(links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(links))
	}
}

func TestExtractLinksInvalidHTML(t *testing.T) {
	reader := strings.NewReader("<html><body")
	links := parser.ExtractLinks(reader)
	if len(links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(links))
	}
}

func TestExtractLinksLinkWithNoHREF(t *testing.T) {
	reader := strings.NewReader("<html><a/></html>")
	links := parser.ExtractLinks(reader)
	if len(links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(links))
	}
}

func TestExtractLinksLinkWithHREF(t *testing.T) {
	reader := strings.NewReader(`<html><a href="http://test.com"></a></html>`)
	links := parser.ExtractLinks(reader)
	if len(links) != 1 {
		t.Errorf("Expected 1 links, got %d", len(links))
		return
	}
	if links[0] != "http://test.com" {
		t.Errorf("Expected http://test.com, got %s", links[0])
	}
}

func TestExtractLinksLinkWithSelfClosingHREF(t *testing.T) {
	reader := strings.NewReader(`<html><a href="http://test.com"/></html>`)
	links := parser.ExtractLinks(reader)
	if len(links) != 1 {
		t.Errorf("Expected 1 links, got %d", len(links))
		return
	}
	if links[0] != "http://test.com" {
		t.Errorf("Expected http://test.com, got %s", links[0])
	}
}

func TestExtractLinksLinkWithNestedLink(t *testing.T) {
	reader := strings.NewReader(`<html><body><h1><a href="http://test.com"></a></h1></body></html>`)
	links := parser.ExtractLinks(reader)
	if len(links) != 1 {
		t.Errorf("Expected 1 links, got %d", len(links))
		return
	}
	if links[0] != "http://test.com" {
		t.Errorf("Expected http://test.com, got %s", links[0])
	}
}
