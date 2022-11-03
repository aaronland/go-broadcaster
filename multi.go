package broadcaster

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type MultiBroadcaster struct {
	Broadcaster
	broadcasters []Broadcaster
	logger       *log.Logger
	async        bool
}

func NewMultiBroadcasterFromURIs(ctx context.Context, broadcaster_uris ...string) (Broadcaster, error) {

	broadcasters := make([]Broadcaster, len(broadcaster_uris))

	for idx, br_uri := range broadcaster_uris {

		br, err := NewBroadcaster(ctx, br_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create broadcaster for '%s', %v", br_uri, err)
		}

		broadcasters[idx] = br
	}

	return NewMultiBroadcaster(ctx, broadcasters...)
}

func NewMultiBroadcaster(ctx context.Context, broadcasters ...Broadcaster) (Broadcaster, error) {

	logger := log.Default()

	async := true

	b := MultiBroadcaster{
		broadcasters: broadcasters,
		logger:       logger,
		async:        async,
	}

	return &b, nil
}

func (b *MultiBroadcaster) BroadcastMessage(ctx context.Context, msg *Message) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	th := b.newThrottle()

	wg := new(sync.WaitGroup)

	for _, bc := range b.broadcasters {

		<-th

		wg.Add(1)

		go func(bc Broadcaster, msg *Message) {

			defer func() {
				th <- true
				wg.Done()
			}()

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			err := bc.BroadcastMessage(ctx, msg)

			if err != nil {
				b.logger.Printf("[%T] Failed to broadcast message: %s\n", bc, err)
			}

		}(bc, msg)
	}

	wg.Wait()
	return nil
}

func (b *MultiBroadcaster) SetLogger(ctx context.Context, logger *log.Logger) error {

	b.logger = logger

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	th := b.newThrottle()

	wg := new(sync.WaitGroup)

	for _, bc := range b.broadcasters {

		<-th

		wg.Add(1)

		go func(bc Broadcaster, logger *log.Logger) {

			defer func() {
				th <- true
				wg.Done()
			}()

			err := bc.SetLogger(ctx, logger)

			if err != nil {
				b.logger.Printf("[%T] Failed to set logger: %v", bc, err)
			}

		}(bc, logger)
	}

	wg.Wait()
	return nil
}

func (b *MultiBroadcaster) newThrottle() chan bool {

	workers := len(b.broadcasters)

	if !b.async {
		workers = 1
	}

	throttle := make(chan bool, workers)

	for i := 0; i < workers; i++ {
		throttle <- true
	}

	return throttle
}
