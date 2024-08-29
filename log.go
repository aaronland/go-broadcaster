package broadcaster

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/aaronland/go-uid"
)

func init() {
	ctx := context.Background()
	RegisterBroadcaster(ctx, "log", NewLogBroadcaster)
}

// LogBroadcaster implements the `Broadcaster` interface to broadcast messages
// to a `log.Logger` instance.
type LogBroadcaster struct {
	Broadcaster
	ascii bool
}

// NewLogBroadcaster returns a new `LogBroadcaster` configured by 'uri' which is expected to
// take the form of:
//
//	log://?{PARAMETERS}
//
// Where `{PARAMETERS}` may be one or more of the following:
// * ?ascii={BOOLEAN} – If true then ASCII-representations of each image in a message will be logged. Default is "false".
// By default `LogBroadcaster` instances are configured to broadcast messages to a `log/slog.Default`
// instance with an `INFO` level.
func NewLogBroadcaster(ctx context.Context, uri string) (Broadcaster, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	ascii := false

	if q.Has("ascii") {

		v, err := strconv.ParseBool(q.Get("ascii"))

		if err != nil {
			return nil, fmt.Errorf("Invalid ?ascii= parameter, %w", err)
		}

		ascii = v
	}

	b := LogBroadcaster{
		ascii: ascii,
	}

	return &b, nil
}

// BroadcastMessage broadcast the title and body properties of 'msg' to the `log.Logger` instance
// associated with 'b'. If the `?ascii=` parameter in the constructor URI (see `NewLogBroadcaster`)
// is true then each image will be converted in to an ASCII representation and logged. It returns
// the value of the Unix timestamp that the log message was broadcast.
func (b *LogBroadcaster) BroadcastMessage(ctx context.Context, msg *Message) (uid.UID, error) {

	log_msg := fmt.Sprintf("%s %s", msg.Title, msg.Body)

	logger := slog.Default()
	logger.Info(log_msg)

	if b.ascii {

		for _, im := range msg.Images {
			im_scaled, w, h := scaleImage(im, 80)
			im_ascii := convert2Ascii(im_scaled, w, h)
			logger.Info(fmt.Sprintf("\n%s\n", string(im_ascii)))
		}
	}

	now := time.Now()
	ts := now.Unix()

	return uid.NewInt64UID(ctx, ts)
}
