//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"context"
	"fmt"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	// Set a timer to process
	wt := make(chan int64)

	limitContext, cancel := context.WithCancel(context.Background())

	if !u.IsPremium {
		limitContext, cancel = context.WithTimeout(context.Background(), (10*time.Second)-time.Duration(u.TimeUsed))
	}

	// Execute a controlled goroutine
	defer cancel()

	go func(ctx context.Context, wastedTime chan<- int64) {
		defer ctx.Done()
		start := time.Now()
		process()
		end := time.Now()
		wastedTime <- end.Unix() - start.Unix()
	}(limitContext, wt)
	select {

	case elapsed := <-wt:
		fmt.Println("Consumed:", elapsed)
		u.TimeUsed += elapsed
		return true
	case <-limitContext.Done():
		fmt.Println("Error:", limitContext.Err())
		return false
	}
}

func main() {
	RunMockServer()
}
