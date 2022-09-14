package main

import (
	"fmt"
	"time"
)

func producer(stream Stream) <-chan *Tweet {
	ch := make(chan *Tweet)
	go func() {
		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				break
			}
			ch <- tweet
		}
		close(ch)
	}()

	return ch
}

func consumer(tweets <-chan *Tweet) {
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
			continue
		}

		fmt.Println(t.Username, "\tdoes not tweet about golang")
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// Producer
	tweets := producer(stream)
	// Consumer
	consumer(tweets)

	fmt.Printf("Process took %s\n", time.Since(start))
}
