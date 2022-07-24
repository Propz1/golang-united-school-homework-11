package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

var mx sync.Mutex

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {

	res = make([]user, 0, n)

	errG, _ := errgroup.WithContext(context.Background())
	errG.SetLimit(int(pool))

	for i := 0; i < int(n); i++ {
		id := int64(i)
		errG.Go(func() error {
			user := getOne(id)
			mx.Lock()
			res = append(res, user)
			mx.Unlock()
			return nil
		})
	}

	err := errG.Wait()

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return res
}
