package broadcaster

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aaronland/go-image-encode"
	"github.com/aaronland/go-mastodon-api/client"
	"github.com/aaronland/go-mastodon-api/response"
	"github.com/sfomuseum/runtimevar"
	_ "image"
	"log"
	"net/url"
	"strconv"
)

func init() {
	ctx := context.Background()
	RegisterBroadcaster(ctx, "mastodon", NewMastodonBroadcaster)
}

type MastodonBroadcaster struct {
	Broadcaster
	mastodon_client client.Client
	testing         bool
	dryrun          bool
	encoder         encode.Encoder
	logger          *log.Logger
}

func NewMastodonBroadcaster(ctx context.Context, uri string) (Broadcaster, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	creds_uri := q.Get("credentials")

	if creds_uri == "" {
		return nil, fmt.Errorf("Missing ?credentials= parameter")
	}

	client_uri, err := runtimevar.StringVar(ctx, creds_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive URI from credentials, %w", err)
	}

	cl, err := client.NewClient(ctx, client_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new Mastodon client, %w", err)
	}

	enc, err := encode.NewEncoder(ctx, "png://")

	if err != nil {
		return nil, fmt.Errorf("Failed to create image encoder, %w", err)
	}

	testing := false
	dryrun := false

	str_testing := q.Get("testing")

	if str_testing != "" {

		t, err := strconv.ParseBool(str_testing)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?testing= parameter, %w", err)
		}

		testing = t
	}

	str_dryrun := q.Get("dryrun")

	if str_dryrun != "" {

		d, err := strconv.ParseBool(str_dryrun)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?dryrun= parameter, %w", err)
		}

		dryrun = d
	}

	logger := log.Default()

	br := &MastodonBroadcaster{
		mastodon_client: cl,
		testing:         testing,
		dryrun:          dryrun,
		encoder:         enc,
		logger:          logger,
	}

	return br, nil
}

func (b *MastodonBroadcaster) BroadcastMessage(ctx context.Context, msg *Message) error {

	status := msg.Body

	if b.testing {
		status = fmt.Sprintf("this is a test and there may be more / please disregard and apologies for the distraction / meanwhile: %s", status)
	}

	args := &url.Values{}

	args.Set("status", status)
	args.Set("visibility", "public")

	if msg.Image != nil {

		// but what if GIF...

		r := new(bytes.Buffer)

		err := b.encoder.Encode(ctx, msg.Image, r)

		if err != nil {
			return fmt.Errorf("Failed to encode image, %w", err)
		}

		if b.dryrun {
			args.Set("media_ids[]", "dryrun")
		} else {

			rsp, err := b.mastodon_client.UploadMedia(ctx, r, nil)

			if err != nil {
				return fmt.Errorf("Failed to upload image, %w", err)
			}

			media_id, err := response.Id(ctx, rsp)

			if err != nil {
				return fmt.Errorf("Failed to derive media ID from response, %w", err)
			}

			args.Set("media_ids[]", media_id)
		}

	}

	if b.dryrun {
		b.logger.Println(args)
		return nil
	}

	rsp, err := b.mastodon_client.ExecuteMethod(ctx, "POST", "/api/v1/statuses", args)

	if err != nil {
		return fmt.Errorf("Failed to post message, %w", err)
	}

	status_id, err := response.Id(ctx, rsp)

	if err != nil {
		return fmt.Errorf("Failed to derive status ID from response, %w", err)
	}

	b.logger.Printf("mastodon post %s", status_id)
	return nil
}

func (b *MastodonBroadcaster) SetLogger(ctx context.Context, logger *log.Logger) error {
	b.logger = logger
	return nil
}
