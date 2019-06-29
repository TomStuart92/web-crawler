package webcrawler

import (
	"io"
	"net/http"
)

func getPage(url string) (error, io.Reader) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	response, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	return nil, response.Body
}
