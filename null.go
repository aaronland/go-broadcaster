package broadcaster

import (
	"context"
	"log"
)

func init() {
	ctx := context.Background()
	RegisterBroadcaster(ctx, "null", NewNullBroadcaster)
}

type NullBroadcaster struct {
	Broadcaster
}

func NewNullBroadcaster(ctx context.Context, uri string) (Broadcaster, error) {

	b := NullBroadcaster{}
	return &b, nil
}

func (b *NullBroadcaster) BroadcastMessage(ctx context.Context, msg *Message) (string, error) {
	return "null", nil
}

func (b *NullBroadcaster) SetLogger(ctx context.Context, logger *log.Logger) error {
	return nil
}
