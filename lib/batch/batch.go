package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	//channel for data
	ch := make(chan user, int(n))
	
	//determine wait group for semaphore logic
	var wg sync.WaitGroup
	//semaphore to limit the number of goroutines
	sem := make(chan struct{}, pool)
	
	//start getting users in cycle with goroutines
	var i int64
	for i = 0; i < n; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(j int64) {
			u := getOne(j)
			ch <- u
			<-sem
			wg.Done()
		}(i)
	}
	//wait until all wait groups are done
	wg.Wait()
	
	//close a channel
	close(ch)

	//read users from channel
	users := make([]user, 0, n)
	for u := range ch {
		users = append(users, u)
	}
	
	return users
}