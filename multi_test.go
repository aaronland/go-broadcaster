package broadcaster

import (
	"context"
	"testing"
)

func TestMultiBroadcaster(t *testing.T) {

	ctx := context.Background()

	null_br, err := NewBroadcaster(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create null:// broadcaster, %v", err)
	}

	log_br, err := NewBroadcaster(ctx, "log://")

	if err != nil {
		t.Fatalf("Failed to create log:// broadcaster, %v", err)
	}

	mb, err := NewMultiBroadcaster(ctx, null_br, log_br)

	if err != nil {
		t.Fatalf("Failed to create multi broadcaster, %v", err)
	}

	msg := &Message{
		Body: "testing",
	}

	_, err = mb.BroadcastMessage(ctx, msg)

	if err != nil {
		t.Fatalf("Failed to broadcast message, %v", err)
	}

}

func TestMultiBroadcasterFromURIs(t *testing.T) {

	ctx := context.Background()

	uris := []string{
		"null://",
		"log://",
	}

	mb, err := NewMultiBroadcasterFromURIs(ctx, uris...)

	if err != nil {
		t.Fatalf("Failed to create multi broadcaster, %v", err)
	}

	msg := &Message{
		Body: "testing",
	}

	_, err = mb.BroadcastMessage(ctx, msg)

	if err != nil {
		t.Fatalf("Failed to broadcast message, %v", err)
	}

}
