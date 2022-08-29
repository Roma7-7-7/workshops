// Hint 1: time.Ticker can be used to cancel function
// Hint 2: to calculate time-diff for Advanced lvl use:
//  start := time.Now()
//	// your work
//	t := time.Now()
//	elapsed := t.Sub(start) // 1s or whatever time has passed

package main

import "time"

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool  // can be used for 2nd level task. Premium users won't have 10 seconds limit.
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	processed := make(chan bool)
	defer close(processed)

	started := time.Now()
	go func() {
		process()
		processed <- true
	}()

	select {
	case <-processed:
		u.TimeUsed = u.TimeUsed + int64(time.Now().Sub(started).Seconds())
		return true
	case <-time.NewTicker(time.Duration(10-u.TimeUsed) * time.Second).C:
		return false
	}
}

func main() {
	RunMockServer()
}
