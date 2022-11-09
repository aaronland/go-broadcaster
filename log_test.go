package broadcaster

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
)

func TestLogBroadcaster(t *testing.T) {

	ctx := context.Background()

	br, err := NewBroadcaster(ctx, "log://")

	if err != nil {
		t.Fatalf("Failed to create log:// broadcaster, %v", err)
	}

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	logger := log.New(wr, "", 0)

	err = br.SetLogger(ctx, logger)

	if err != nil {
		t.Fatalf("Failed to set logger, %v", err)
	}

	msg := &Message{
		Title: "Testing",
		Body:  "Hello world",
	}

	_, err = br.BroadcastMessage(ctx, msg)

	if err != nil {
		t.Fatalf("Failed to broadcast message, %v", err)
	}

	wr.Flush()

	if strings.TrimSpace(buf.String()) != "Testing Hello world" {
		t.Fatalf("Unexpected output '%s'", buf.String())
	}

}
