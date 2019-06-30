package webcrawler_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/TomStuart92/web-crawler/webcrawler"
)

func TestNew(t *testing.T) {
	wc := webcrawler.New(1)
	if wc == nil {
		t.Error("New returned nil, instead of wc instance")
	}
}

func TestGenerateSiteMap(t *testing.T) {
	body := strings.NewReader(`<html><a/ href="http://test.com/about"></html>`)
	response := http.Response{Body: ioutil.NopCloser(body)}
	client := HTTPFake{&response, nil}
	wc := webcrawler.New(1)
	wc.Client = &client
	wc.GenerateSiteMap("http://test.com", true)
}

func TestPrintSiteMap(t *testing.T) {
	body := strings.NewReader(`<html><a/ href="http://test.com/about"></html>`)
	response := http.Response{Body: ioutil.NopCloser(body)}
	client := HTTPFake{&response, nil}
	wc := webcrawler.New(1)
	wc.Client = &client
	var buf bytes.Buffer
	wc.GenerateSiteMap("http://test.com", true)
	wc.PrintSiteMap(&buf)
	if buf.String() != "\nhttp://test.com ---> http://test.com/about\n" {
		t.Errorf("Expected: http://test.com ---> http://test.com/about, got: %s", buf.String())
	}
}
