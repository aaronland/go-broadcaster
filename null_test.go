package broadcaster

import (
	"context"
	"testing"
)

func TestNullBroadcaster(t *testing.T) {

	ctx := context.Background()

	br, err := NewBroadcaster(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create null:// broadcaster, %v", err)
	}

	msg := &Message{
		Body: "testing",
	}

	_, err = br.BroadcastMessage(ctx, msg)

	if err != nil {
		t.Fatalf("Failed to broadcast message, %v", err)
	}
}
