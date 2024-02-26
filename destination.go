package activemq

import (
	"context"
	"fmt"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/go-stomp/stomp"
)

type Destination struct {
	sdk.UnimplementedDestination
	config DestinationConfig

	conn *stomp.Conn
}

func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

func (d *Destination) Parameters() map[string]sdk.Parameter {
	return d.config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) (err error) {
	err = sdk.Util.ParseConfig(cfg, &d.config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	sdk.Logger(ctx).Debug().Any("config", d.config).Msg("configured destination")

	return nil
}

func (d *Destination) Open(ctx context.Context) (err error) {
	d.conn, err = stomp.Dial("tcp", d.config.URL)
	if err != nil {
		return fmt.Errorf("failed to dial to ActiveMQ: %w", err)
	}

	sdk.Logger(ctx).Debug().Msg("opened destination")
	return nil
}

func (d *Destination) Write(ctx context.Context, records []sdk.Record) (int, error) {
	for _, rec := range records {
		err := d.conn.Send(d.config.Queue, d.config.ContentType, rec.Bytes())
		if err != nil {
			return 0, fmt.Errorf("failed to send message: %w", err)
		}
	}

	return len(records), nil
}

func (d *Destination) Teardown(ctx context.Context) error {
	if err := d.conn.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect from ActiveMQ: %w", err)
	}

	sdk.Logger(ctx).Debug().Msg("teardown destination")
	return nil
}
