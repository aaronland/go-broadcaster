package broadcaster

import (
	"context"
	"log"
)

func init() {
	ctx := context.Background()
	RegisterBroadcaster(ctx, "log", NewLogBroadcaster)
}

type LogBroadcaster struct {
	Broadcaster
	logger *log.Logger
}

func NewLogBroadcaster(ctx context.Context, uri string) (Broadcaster, error) {
	logger := log.Default()
	b := LogBroadcaster{
		logger: logger,
	}
	return &b, nil
}

func (b *LogBroadcaster) BroadcastMessage(ctx context.Context, msg *Message) error {
	b.logger.Println(msg.Body)
	return nil
}

func (b *LogBroadcaster) SetLogger(ctx context.Context, logger *log.Logger) error {
	b.logger = logger
	return nil
}
