package activemq

import (
	"context"

	sdk "github.com/conduitio/conduit-connector-sdk"
)


type Source struct {
	sdk.UnimplementedSource
	config    SourceConfig
}

func NewSource() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{}, sdk.DefaultSourceMiddleware()...)
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
}

func (s *Source) Open(ctx context.Context, sdkPos sdk.Position) (err error) {
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
}

func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
}

func (s *Source) Teardown(ctx context.Context) error {
}
