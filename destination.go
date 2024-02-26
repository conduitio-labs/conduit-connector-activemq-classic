package activemq

import (
	"context"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Destination struct {
	sdk.UnimplementedDestination
	config DestinationConfig
}

func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

func (d *Destination) Parameters() map[string]sdk.Parameter {
	return d.config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) (err error) {
}

func (d *Destination) Open(ctx context.Context) (err error) {
}

func (d *Destination) Write(ctx context.Context, records []sdk.Record) (int, error) {
}

func (d *Destination) Teardown(_ context.Context) error {
}
