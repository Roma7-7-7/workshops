package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const stringToSearch = "concurrency"

var sites = []string{
	"https://google.com",
	"https://itc.ua/",
	"https://twitter.com/concurrencyinc",
	"https://twitter.com/",
	"http://localhost:8000",
	"https://github.com/bradtraversy/go_restapi/blob/master/main.go",
	"https://www.youtube.com/",
	"https://postman-echo.com/get",
	"https://en.wikipedia.org/wiki/Concurrency_(computer_science)#:~:text=In%20computer%20science%2C%20concurrency%20is,without%20affecting%20the%20final%20outcome.",
}

type SiteData struct {
	data []byte
	uri  string
}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	resultsCh := make(chan SiteData, len(sites))

	crawl(ctx, cancel, resultsCh)

	// give one second to validate if all other goroutines are closed
	time.Sleep(time.Second)
}

func crawl(ctx context.Context, cancel func(), results chan SiteData) {
	for _, uri := range sites {
		go request(ctx, uri, results)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case result := <-results:
			if strings.Contains(string(result.data), stringToSearch) {
				fmt.Printf("%s string is found in %s\n", stringToSearch, result.uri)
				cancel()
				return
			} else {
				fmt.Printf("Nothing found in %s\n", result.uri)
			}
		case <-time.NewTimer(time.Minute).C:
			fmt.Println("timeout")
			cancel()
			return
		}
	}
}

func request(ctx context.Context, uri string, results chan<- SiteData) {
	fmt.Printf("starting sending request to %s\n", uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	results <- SiteData{
		data: bodyBytes,
		uri:  uri,
	}
}
