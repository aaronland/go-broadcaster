package broadcaster

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
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

	logger := slog.New(slog.NewTextHandler(wr, nil))
	slog.SetDefault(logger)

	msg := &Message{
		Title: "Testing",
		Body:  "Hello world",
	}

	_, err = br.BroadcastMessage(ctx, msg)

	if err != nil {
		t.Fatalf("Failed to broadcast message, %v", err)
	}

	wr.Flush()

	if strings.HasSuffix(buf.String(), `level=INFO msg="Testing Hello world"`) {
		t.Fatalf("Unexpected output '%s'", buf.String())
	}

}
