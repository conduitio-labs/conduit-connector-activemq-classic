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
	sourceCfg := map[string]string{
		SourceConfigUrl:      "localhost:61613",
		SourceConfigUser:     "admin",
		SourceConfigPassword: "admin",
	}
	destCfg := map[string]string{
		DestinationConfigUrl:      "localhost:61613",
		DestinationConfigUser:     "admin",
		DestinationConfigPassword: "admin",
	}

	driver := sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector:         Connector,
			SourceConfig:      sourceCfg,
			DestinationConfig: destCfg,
			BeforeTest: func(t *testing.T) {
				queueName := uniqueQueueName(t)
				sourceCfg[SourceConfigQueue] = queueName
				destCfg[DestinationConfigQueue] = queueName
			},
			WriteTimeout: 500 * time.Millisecond,
			ReadTimeout:  500 * time.Millisecond,
		},
	}

	sdk.AcceptanceTest(t, driver)
}

func TestAcceptanceTLS(t *testing.T) {
	sourceCfg := map[string]string{
		SourceConfigUrl:                   "localhost:61617",
		SourceConfigUser:                  "admin",
		SourceConfigPassword:              "admin",
		SourceConfigTlsEnabled:            "true",
		SourceConfigTlsClientKeyPath:      "./test/certs/client_key.pem",
		SourceConfigTlsClientCertPath:     "./test/certs/client_cert.pem",
		SourceConfigTlsCaCertPath:         "./test/certs/broker.pem",
		SourceConfigTlsInsecureSkipVerify: "true",
	}

	destCfg := map[string]string{
		DestinationConfigUrl:                   "localhost:61617",
		DestinationConfigUser:                  "admin",
		DestinationConfigPassword:              "admin",
		DestinationConfigTlsEnabled:            "true",
		DestinationConfigTlsClientKeyPath:      "./test/certs/client_key.pem",
		DestinationConfigTlsClientCertPath:     "./test/certs/client_cert.pem",
		DestinationConfigTlsCaCertPath:         "./test/certs/broker.pem",
		DestinationConfigTlsInsecureSkipVerify: "true",
	}

	driver := sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector:         Connector,
			SourceConfig:      sourceCfg,
			DestinationConfig: destCfg,
			BeforeTest: func(t *testing.T) {
				queueName := uniqueQueueName(t)
				sourceCfg[SourceConfigQueue] = queueName
				destCfg[DestinationConfigQueue] = queueName
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
