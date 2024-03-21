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
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/go-stomp/stomp/v3"
)

type Position struct {
	MessageID string `json:"message_id"`
	Queue     string `json:"queue"`
}

func parseSDKPosition(sdkPos sdk.Position) (Position, error) {
	decoder := json.NewDecoder(bytes.NewBuffer(sdkPos))
	decoder.DisallowUnknownFields()

	var p Position
	if err := decoder.Decode(&p); err != nil {
		return p, fmt.Errorf("failed to parse position: %w", err)
	}

	return p, nil
}

func (p Position) ToSdkPosition() sdk.Position {
	bs, err := json.Marshal(p)
	if err != nil {
		// this should never happen
		panic(err)
	}

	return sdk.Position(bs)
}

// metadataFromMsg extracts all the present headers from a stomp.Message into
// sdk.Metadata.
func metadataFromMsg(msg *stomp.Message) sdk.Metadata {
	metadata := make(sdk.Metadata)

	for i := range msg.Header.Len() {
		k, v := msg.Header.GetAt(i)
		metadata[k] = v
	}

	return metadata
}

func connect(ctx context.Context, config Config) (*stomp.Conn, error) {
	loginOpt := stomp.ConnOpt.Login(config.User, config.Password)
	heartbeat := stomp.ConnOpt.HeartBeat(config.SendTimeoutHeartbeat, config.RecvTimeoutHeartbeat)
	if !config.TLS.Enabled {
		conn, err := stomp.Dial("tcp", config.URL, loginOpt, heartbeat)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to ActiveMQ: %w", err)
		}
		sdk.Logger(ctx).Debug().Msg("opened connection to ActiveMQ")

		return conn, nil
	}

	sdk.Logger(ctx).Debug().Msg("using TLS to connect to ActiveMQ")

	cert, err := tls.LoadX509KeyPair(config.TLS.ClientCertPath, config.TLS.ClientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client key pair: %w", err)
	}
	sdk.Logger(ctx).Debug().Msg("loaded client key pair")

	caCert, err := os.ReadFile(config.TLS.CaCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load CA cert: %w", err)
	}
	sdk.Logger(ctx).Debug().Msg("loaded CA cert")

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,

		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// version will be overwritten at compile time when building a release,
	// so this should only be true when running in development mode.
	if version == "(devel)" {
		tlsConfig.InsecureSkipVerify = true
	}

	netConn, err := tls.Dial("tcp", config.URL, tlsConfig)
	if err != nil {
		panic(err)
	}
	sdk.Logger(ctx).Debug().Msg("TLS connection established")

	conn, err := stomp.Connect(netConn, loginOpt, heartbeat)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ActiveMQ: %w", err)
	}
	sdk.Logger(ctx).Debug().Msg("STOMP connection using tls established")

	return conn, nil
}

func teardown(ctx context.Context, subs *stomp.Subscription, conn *stomp.Conn) error {
	if subs != nil {
		err := subs.Unsubscribe()
		if errors.Is(err, stomp.ErrCompletedSubscription) {
			sdk.Logger(ctx).Debug().Msg("subscription already unsubscribed")
		} else if err != nil {
			return fmt.Errorf("failed to unsubscribe: %w", err)
		}
	}

	if conn != nil {
		if err := conn.Disconnect(); err != nil {
			return fmt.Errorf("failed to disconnect from ActiveMQ: %w", err)
		}
	}

	return nil
}
