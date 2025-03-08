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
		"url":      "localhost:61613",
		"user":     "admin",
		"password": "admin",
	}
	destCfg := map[string]string{
		"url":      "localhost:61613",
		"user":     "admin",
		"password": "admin",
	}

	driver := sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector:         Connector,
			SourceConfig:      sourceCfg,
			DestinationConfig: destCfg,
			BeforeTest: func(t *testing.T) {
				queueName := uniqueQueueName(t)
				sourceCfg["queue"] = queueName
				destCfg["queue"] = queueName
			},
			WriteTimeout: 500 * time.Millisecond,
			ReadTimeout:  500 * time.Millisecond,
		},
	}

	sdk.AcceptanceTest(t, driver)
}

func TestAcceptanceTLS(t *testing.T) {
	sourceCfg := map[string]string{
		"url":                    "localhost:61617",
		"user":                   "admin",
		"password":               "admin",
		"tls.enabled":            "true",
		"tls.clientKeyPath":      "./test/certs/client_key.pem",
		"tls.clientCertPath":     "./test/certs/client_cert.pem",
		"tls.caCertPath":         "./test/certs/broker.pem",
		"tls.insecureSkipVerify": "true",
	}

	destCfg := map[string]string{
		"url":                    "localhost:61617",
		"user":                   "admin",
		"password":               "admin",
		"tls.enabled":            "true",
		"tls.clientKeyPath":      "./test/certs/client_key.pem",
		"tls.clientCertPath":     "./test/certs/client_cert.pem",
		"tls.caCertPath":         "./test/certs/broker.pem",
		"tls.insecureSkipVerify": "true",
	}

	driver := sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector:         Connector,
			SourceConfig:      sourceCfg,
			DestinationConfig: destCfg,
			BeforeTest: func(t *testing.T) {
				queueName := uniqueQueueName(t)
				sourceCfg["queue"] = queueName
				destCfg["queue"] = queueName
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
