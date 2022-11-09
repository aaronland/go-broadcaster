package broadcaster

import (
	"context"
	"testing"
)

func TestRegisterBroadcaster(t *testing.T) {

	ctx := context.Background()

	err := RegisterBroadcaster(ctx, "null", NewNullBroadcaster)

	if err == nil {
		t.Fatalf("Expected error registering null:// scheme")
	}
}
