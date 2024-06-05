//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, wg *sync.WaitGroup) <-chan *Tweet {

	tweetChan := make(chan *Tweet)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				close(tweetChan)
				return
			} else {
				tweetChan <- tweet
			}
		}
	}()
	return tweetChan
}

func consumer(tweetChan <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		t, more := <-tweetChan
		if !more {
			return
		}
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	stream := GetMockStream()

	// Producer
	tweetChan := producer(stream, &wg)

	// Consumer
	wg.Add(1)
	go consumer(tweetChan, &wg)

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
