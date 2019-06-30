// Package webcrawler exports a non-thread safe webcrawler.
package webcrawler

import (
	"log"
)

// FetchURLs concurrently reads urls from a work queue,
// and makes a http get request to the url
func FetchURLs(client HTTPClient, concurency int, workChan chan Page) (chan Page, chan error) {
	resultsChan := make(chan Page, concurency)
	errChan := make(chan error)

	for i := 0; i < concurency; i++ {
		go func(workerID int) {
			for page := range workChan {
				log.Printf("Worker %d, processing %s", workerID, page.URL)
				result, err := client.Get(page.URL)
				if err != nil {
					errChan <- err
					continue
				}
				page.Body = result.Body
				resultsChan <- page
			}
		}(i)
	}
	return resultsChan, errChan
}
