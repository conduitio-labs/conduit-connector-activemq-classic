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

import "time"

//go:generate paramgen -output=paramgen_src.go SourceConfig
//go:generate paramgen -output=paramgen_dest.go DestinationConfig

type Config struct {
	// URL is the URL of the ActiveMQ classic broker.
	URL string `json:"url" validate:"required"`

	// User is the username to use when connecting to the broker.
	User string `json:"user" validate:"required"`

	// Password is the password to use when connecting to the broker.
	Password string `json:"password" validate:"required"`

	// Queue is the name of the queue to write to.
	Queue string `json:"queue" validate:"required"`

	// SendTimeoutHeartbeat specifies the maximum amount of time between the
	// client sending heartbeat notifications to the server
	SendTimeoutHeartbeat time.Duration `json:"sendTimeoutHeartbeat" default:"2s"`

	// RecvTimeoutHeartbeat specifies the minimum amount of time between the
	// client expecting to receive heartbeat notifications from the server
	RecvTimeoutHeartbeat time.Duration `json:"recvTimeoutHeartbeat" default:"2s"`

	TLS TLSConfig `json:"tlsConfig"`
}

type TLSConfig struct {
	// UseTLS is a flag to enable or disable TLS.
	UseTLS bool `json:"useTLS" default:"false"`

	// ClientKeyPath is the path to the client key file.
	ClientKeyPath string `json:"clientKeyPath"`

	// ClientCertPath is the path to the client certificate file.
	ClientCertPath string `json:"clientCertPath"`

	// CaCertPath is the path to the CA certificate file.
	CaCertPath string `json:"caCertPath"`
}

type SourceConfig struct {
	Config
}

type DestinationConfig struct {
	Config
}
