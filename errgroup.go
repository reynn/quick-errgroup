package quickerr

import (
	"context"
	"fmt"
	"sync"
)

type Group struct {
	wg *sync.WaitGroup

	errChn    chan error
	canceller context.CancelFunc
	ctx       context.Context
}

func New(parentCtx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)

	return &Group{
		wg:        &sync.WaitGroup{},
		errChn:    make(chan error, 1),
		canceller: cancel,
		ctx:       ctx,
	}, ctx
}

func (q *Group) Go(fn func() error) {
	q.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		if err := fn(); err != nil {
			q.errChn <- err
		}
	}(q.wg)
}

func (q *Group) Wait() error {
	waiter := make(chan bool, 1)
	go func() {
		q.wg.Wait()
		fmt.Println("All WaitGroup routines complete")
		waiter <- true
	}()

	select {
	case err := <-q.errChn:
		q.canceller()
		return err
	case <-q.ctx.Done():
		q.canceller()
		return context.Canceled
	case <-waiter:
		q.canceller()
		return nil
	}
}
