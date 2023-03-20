package main

import (
	"bytes"
	"log"
	"net/http"
	"sync"
)

func loadTest() {
	var wg sync.WaitGroup

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				body := []byte(`{
					"name": "Macbook Pro 13 M1 32 GB 1 TB",
					"qtty": 10,
					"price": 73000000
				}`)
				request, errNewRequest := http.NewRequest(http.MethodPost, "http://localhost/api/v1/products", bytes.NewBuffer(body))
				if errNewRequest != nil {
					log.Printf("Error new request: %v", errNewRequest)
				}
				request.Header.Set("Content-Type", "application/json")

				response, errDo := client.Do(request)
				if errDo != nil {
					log.Printf("Error do: %v", errDo)
				}
				if response != nil {
					defer response.Body.Close()
					log.Printf("status: %v", response.StatusCode)
				}
			}
		}()
	}
	wg.Wait()
}

func main() {
	// if errPublishMessage := PublishMessage("purchases", []byte("tess gan")); errPublishMessage != nil {
	// 	log.Fatalf("errPublishMessage: %v", errPublishMessage)
	// }
	ReceiveMessage()
}
