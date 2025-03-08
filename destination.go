// Copyright © 2024 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package activemq

import (
	"context"
	"fmt"

	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/go-stomp/stomp/v3"
)

type DestinationConfig struct {
	sdk.DefaultDestinationMiddleware

	Config
}

type Destination struct {
	sdk.UnimplementedDestination
	config DestinationConfig

	conn *stomp.Conn
}

func (d *Destination) Config() sdk.DestinationConfig {
	return &d.config
}

func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{})
}

func (d *Destination) Open(ctx context.Context) (err error) {
	d.conn, err = connectDestination(ctx, d.config)
	if err != nil {
		return fmt.Errorf("failed to dial to ActiveMQ: %w", err)
	}
	sdk.Logger(ctx).Debug().Msg("opened destination")

	return nil
}

func (d *Destination) Write(ctx context.Context, records []opencdc.Record) (int, error) {
	for i, rec := range records {
		err := d.conn.Send(d.config.Queue, "application/json", rec.Bytes())
		if err != nil {
			return i, fmt.Errorf("failed to send message: %w", err)
		}
		sdk.Logger(ctx).Trace().Str("queue", d.config.Queue).Msg("wrote record")
	}

	return len(records), nil
}

func (d *Destination) Teardown(ctx context.Context) error {
	return teardown(ctx, nil, d.conn)
}
