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
	"testing"
	"time"

	"github.com/rs/zerolog"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

func TestAcceptance(t *testing.T) {
	cfg := map[string]string{
		"url": "localhost:61613",
	}

	logger := zerolog.New(zerolog.NewConsoleWriter())
	ctx := logger.WithContext(context.Background())

	driver := sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Context:           ctx,
			Connector:         Connector,
			SourceConfig:      cfg,
			DestinationConfig: cfg,
			BeforeTest: func(t *testing.T) {
				// Ideally we could delete the queue before the test to ensure a clean
				// slate. However, I don't see a clear way to do this, so I'll assume that
				// the docker container was started from scratch.

				cfg["queue"] = t.Name()
			},
			Skip: []string{
				"TestSource_Configure_RequiredParams",
				"TestDestination_Configure_RequiredParams",
			},
			WriteTimeout: 500 * time.Millisecond,
			ReadTimeout:  500 * time.Millisecond,
		},
	}

	sdk.AcceptanceTest(t, driver)
}
