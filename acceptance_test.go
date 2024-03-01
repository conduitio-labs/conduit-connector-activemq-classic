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
	"fmt"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

func TestAcceptance(t *testing.T) {
	cfg := map[string]string{
		"url":      "localhost:61613",
		"user":     "admin",
		"password": "admin",
	}

	driver := sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector:         Connector,
			SourceConfig:      cfg,
			DestinationConfig: cfg,
			BeforeTest: func(t *testing.T) {
				// Ideally we would delete the queue before the test to ensure
				// a clean slate. I don't see a clear way to do this,
				// so I'll assume that the docker container was started from
				// scratch.

				cfg["queue"] = uniqueQueueName(t)
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

func TestAcceptanceTLS(t *testing.T) {
	cfg := map[string]string{
		"url":                      "localhost:61617",
		"user":                     "admin",
		"password":                 "admin",
		"tlsConfig.useTLS":         "true",
		"tlsConfig.clientKeyPath":  "./test/certs/client_key.pem",
		"tlsConfig.clientCertPath": "./test/certs/client_cert.pem",
		"tlsConfig.caCertPath":     "./test/certs/broker.pem",
	}

	driver := sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector:         Connector,
			SourceConfig:      cfg,
			DestinationConfig: cfg,
			BeforeTest: func(t *testing.T) {
				cfg["queue"] = uniqueQueueName(t)
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

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func uniqueQueueName(t *testing.T) string {
	return fmt.Sprintf("%v_%v", t.Name(), randomString())
}
