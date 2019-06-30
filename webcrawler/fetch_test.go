package webcrawler_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/TomStuart92/web-crawler/webcrawler"
)

type HTTPFake struct {
	Response *http.Response
	Error    error
}

func (h *HTTPFake) Get(url string) (resp *http.Response, err error) {
	if h.Error != nil {
		return nil, h.Error
	}
	return h.Response, nil
}

func TestFetch(t *testing.T) {
	body := strings.NewReader(`<html><body/></html>`)
	response := http.Response{Body: ioutil.NopCloser(body)}
	client := HTTPFake{&response, nil}
	concurrency := 1
	workChan := make(chan webcrawler.Page)

	outChan, _ := webcrawler.FetchURLs(&client, concurrency, workChan)

	page := webcrawler.Page{URL: "http://test.com"}

	go func() { workChan <- page }()

	outPage := <-outChan

	if outPage.URL != page.URL {
		t.Error("Fetch Mutated page URL")
		return
	}

	if outPage.Body == nil {
		t.Error("Fetch Failed to set body")
		return
	}
}

func TestFetchWithError(t *testing.T) {
	err := errors.New("Fake Error")
	client := HTTPFake{nil, err}
	concurrency := 1
	workChan := make(chan webcrawler.Page)

	_, errChan := webcrawler.FetchURLs(&client, concurrency, workChan)

	page := webcrawler.Page{URL: "http://test.com"}

	go func() { workChan <- page }()

	outErr := <-errChan

	if outErr == nil {
		t.Error("Fetch Failed to return Error")
		return
	}

	if outErr.Error() != err.Error() {
		t.Error("Fetch Returned Wrong Error")
		return
	}
}
