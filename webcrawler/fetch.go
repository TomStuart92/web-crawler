package webcrawler

import (
	"log"
	"net/http"
)

type httpClient interface {
	Get(url string) (resp http.Response, err error)
}

// FetchURLs concurrently reads urls from a work queue,
// and makes a http get request to the url
func FetchURLs(concurency int, workChan chan Page) (chan Page, chan error) {
	resultsChan := make(chan Page, concurency)
	errChan := make(chan error)

	client := &http.Client{}

	for i := 0; i < concurency; i++ {
		go func(workerID int) {
			for page := range workChan {
				log.Printf("Worker %d, processing %s", workerID, page.URL)
				result, err := client.Get(page.URL)
				go func(p Page) {
					if err != nil {
						errChan <- err
					}
					p.Body = result.Body
					resultsChan <- p
				}(page)
			}
		}(i)
	}
	return resultsChan, errChan
}
