package activemq

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/go-stomp/stomp"
)

type Source struct {
	sdk.UnimplementedSource
	config SourceConfig

	conn         *stomp.Conn
	subscription *stomp.Subscription
}

func NewSource() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{}, sdk.DefaultSourceMiddleware()...)
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	err := sdk.Util.ParseConfig(cfg, &s.config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	sdk.Logger(ctx).Debug().Any("config", s.config).Msg("configured source")

	return nil
}

func (s *Source) Open(ctx context.Context, sdkPos sdk.Position) (err error) {
	s.conn, err = stomp.Dial("tcp", s.config.URL)
	if err != nil {
		return fmt.Errorf("failed to dial to ActiveMQ: %w", err)
	}

	s.subscription, err = s.conn.Subscribe(s.config.Queue, stomp.AckClientIndividual)
	if err != nil {
		return fmt.Errorf("failed to subscribe to queue: %w", err)
	}

	sdk.Logger(ctx).Debug().Msg("opened source")

	return nil
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	var rec sdk.Record

	select {
	case <-ctx.Done():
		return rec, ctx.Err()
	case msg, ok := <-s.subscription.C:
		if !ok {
			return rec, errors.New("source message channel closed")
		}

		if err := msg.Err; err != nil {
			return rec, fmt.Errorf("source message error: %w", err)
		}

		var (
			pos = Position{
				Queue: s.config.Queue,
			}
			sdkPos   = pos.ToSdkPosition()
			metadata = sdk.Metadata(nil)
			key      = sdk.RawData(nil)
			payload  = sdk.RawData(msg.Body)
		)

		rec = sdk.Util.Source.NewRecordCreate(sdkPos, metadata, key, payload)

		sdk.Logger(ctx).Trace().
			Str("queue", s.config.Queue).
			Msgf("read message")

		return rec, nil
	}
}

func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
	pos, err := parseSDKPosition(position)
	if err != nil {
		return fmt.Errorf("failed to parse position: %w", err)
	}

	// The go-stomp library doesn't provide another way to ack a message than
	// to give the message itself, so we try to recreate the message from the
	// position and ack it.

	fakeMsg := pos.toMsg(s)
	if err := s.conn.Ack(fakeMsg); err != nil {
		return fmt.Errorf("failed to ack message: %w", err)
	}

	sdk.Logger(ctx).Trace().Msg("acked message")
	return nil
}

func (s *Source) Teardown(ctx context.Context) error {
	if err := s.subscription.Unsubscribe(); err != nil {
		return fmt.Errorf("failed to unsubscribe from queue: %w", err)
	}

	if err := s.conn.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect from ActiveMQ: %w", err)
	}

	sdk.Logger(ctx).Debug().Msg("teardown source")
	return nil
}
