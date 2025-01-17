package queue

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"getsturdy.com/api/pkg/queue/names"
)

var _ Queue = &Sync{}

type Sync struct {
	chansGuard *sync.RWMutex
	chans      map[names.IncompleteQueueName][]chan<- Message
}

func NewSync() *Sync {
	return &Sync{
		chansGuard: &sync.RWMutex{},
		chans:      make(map[names.IncompleteQueueName][]chan<- Message),
	}
}

func (q *Sync) Publish(ctx context.Context, name names.IncompleteQueueName, msg any) error {
	q.chansGuard.RLock()
	chs, ok := q.chans[name]
	q.chansGuard.RUnlock()
	if !ok {
		return fmt.Errorf("no subscriber for %s was found", name)
	}

	wg, _ := errgroup.WithContext(ctx)
	for _, ch := range chs {
		ch := ch
		wg.Go(func() error {
			m, err := newInmemoryMessage(msg)
			if err != nil {
				return fmt.Errorf("failed to create message: %w", err)
			}
			ch <- m
			m.AwaitAcked()
			return nil
		})
	}
	return wg.Wait()
}

func (q *Sync) Subscribe(ctx context.Context, name names.IncompleteQueueName, messages chan<- Message) error {
	q.chansGuard.Lock()
	q.chans[name] = append(q.chans[name], messages)
	q.chansGuard.Unlock()
	<-ctx.Done()
	return nil
}
