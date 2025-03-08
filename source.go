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
	"errors"
	"fmt"
	"strings"

	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
	"github.com/goccy/go-json"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type SourceConfig struct {
	sdk.DefaultSourceMiddleware

	Config

	// ClientID specifies the JMS clientID which is used in combination with
	// the activemq.subcriptionName to denote a durable subscriber.
	// Maps to the client-id header.
	ClientID string `json:"clientID"`

	// DispatchAsync specifies whether messages should be dispatched
	// synchronously or asynchronously from the producer thread for non-durable
	// topics in the broker.
	// Maps to the activemq.dispatchAsync header.
	DispatchAsync bool `json:"activemq.dispatchAsync"`

	// Exclusive indicates the desire to be the sole consumer from the queue.
	// Maps to the activemq.exclusive header.
	Exclusive bool `json:"activemq.exclusive"`

	// MaxPendingMessageLimit specifies the upper limit of pending messages
	// allowed for slow consumers on non-durable topics. When this limit is
	// reached, older messages will be discarded to handle slow consumer
	// backlog.
	// Maps to the activemq.maximumPendingMessageLimit header.
	MaxPendingMessageLimit int `json:"activemq.maximumPendingMessageLimit"`

	// NoLocal indicates if messages sent from the local connection should be
	// excluded from subscriptions. When set to true, locally sent messages
	// will be ignored.
	// Maps to the activemq.noLocal header.
	NoLocal bool `json:"activemq.noLocal"`

	// PrefetchSize determines the maximum number of messages to dispatch to the client
	// before it acknowledges a message. No further messages are dispatched once this
	// limit is hit. For fair message distribution across consumers, consider setting
	// this to a value greater than 1.
	// Maps to the activemq.prefetchSize header.
	PrefetchSize int `json:"activemq.prefetchSize"`

	// Priority specifies the consumer's priority level for weighted dispatching order.
	// Maps to the activemq.priority header.
	Priority byte `json:"activemq.priority"`

	// Retroactive, if set to true, makes the subscription retroactive for non-durable topics.
	// Maps to the activemq.retroactive header.
	Retroactive bool `json:"activemq.retroactive"`

	// SubscriptionName specifies the name used for durable topic subscriptions.
	// Prior to ActiveMQ version 5.7.0, both clientID on the connection and
	// subscriptionName  on the subscribe operation must match.
	// Maps to the activemq.subscriptionName header.
	SubscriptionName string `json:"activemq.subscriptionName"`

	// Selector defines a JMS Selector employing SQL 92 syntax as delineated in
	// the JMS 1.1 specification, enabling a filter to be applied on each
	// message associated with the subscription.
	// Maps to the selector header.
	Selector string `json:"selector"`
}

type Source struct {
	sdk.UnimplementedSource
	config SourceConfig

	conn         *stomp.Conn
	subscription *stomp.Subscription

	storedMessages cmap.ConcurrentMap[string, *stomp.Message]
}

func (s *Source) Config() sdk.SourceConfig {
	return &s.config
}

func NewSource() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{
		storedMessages: cmap.New[*stomp.Message](),
	})
}

// getSubscribeOpts gets all configurable STOMP SUBSCRIPTION frame options from SourceConfig.
// Header names come from here:
// https://activemq.apache.org/components/classic/documentation/stomp#activemq-classic-extensions-to-stomp
func getSubscribeOpts(config SourceConfig) []func(*frame.Frame) error {
	opts := []func(*frame.Frame) error{}
	addHeader := func(key, value string) {
		opts = append(opts, stomp.SubscribeOpt.Header(key, value))
	}

	if config.DispatchAsync {
		addHeader("activemq.dispatchAsync", "true")
	}

	if config.Exclusive {
		addHeader("activemq.exclusive", "true")
	}

	if config.MaxPendingMessageLimit > 0 {
		value := fmt.Sprint(config.MaxPendingMessageLimit)
		addHeader("activemq.maximumPendingMessageLimit", value)
	}

	if config.NoLocal {
		addHeader("activemq.noLocal", "true")
	}

	if config.PrefetchSize > 0 {
		value := fmt.Sprint(config.PrefetchSize)
		addHeader("activemq.prefetchSize", value)
	}

	if config.Priority > 0 {
		value := fmt.Sprint(config.Priority)
		addHeader("activemq.priority", value)
	}

	if config.Retroactive {
		addHeader("activemq.retroactive", "true")
	}

	if config.SubscriptionName != "" {
		addHeader("activemq.subscriptionName", config.SubscriptionName)
	}

	if config.Selector != "" {
		addHeader("selector", config.Selector)
	}

	return opts
}

func (s *Source) Open(ctx context.Context, sdkPos opencdc.Position) (err error) {
	s.conn, err = connectSource(ctx, s.config)
	if err != nil {
		return fmt.Errorf("failed to dial to ActiveMQ: %w", err)
	}

	if sdkPos != nil {
		pos, err := parseSDKPosition(sdkPos)
		if err != nil {
			return fmt.Errorf("failed to parse position: %w", err)
		}

		if s.config.Queue != "" && s.config.Queue != pos.Queue {
			return fmt.Errorf(
				"the old position contains a different queue name than the connector configuration (%q vs %q), please check if the configured queue name changed since the last run",
				pos.Queue, s.config.Queue,
			)
		}

		s.config.Queue = pos.Queue
		sdk.Logger(ctx).Debug().Str("queue", pos.Queue).Msg("got queue name from given position")
	}

	subscribeOpts := getSubscribeOpts(s.config)
	s.subscription, err = s.conn.Subscribe(
		s.config.Queue, stomp.AckClientIndividual,
		subscribeOpts...)
	if err != nil {
		return fmt.Errorf("failed to subscribe to queue: %w", err)
	}

	sdk.Logger(ctx).Debug().Msg("opened source")

	return nil
}

func (s *Source) Read(ctx context.Context) (opencdc.Record, error) {
	var rec opencdc.Record

	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			return rec, fmt.Errorf("context error: %w", err)
		}

		return rec, nil
	case msg, ok := <-s.subscription.C:
		if !ok {
			return rec, errors.New("source message channel closed")
		}

		if err := msg.Err; err != nil {
			return rec, fmt.Errorf("source message error: %w", err)
		}

		var (
			messageID = msg.Header.Get(frame.MessageId)
			pos       = Position{
				MessageID: messageID,
				Queue:     s.config.Queue,
			}
			sdkPos   = pos.ToSdkPosition()
			metadata = metadataFromMsg(msg)
			key      = opencdc.RawData(messageID)
			payload  = opencdc.RawData(msg.Body)
		)

		rec = sdk.Util.Source.NewRecordCreate(sdkPos, metadata, key, payload)

		sdk.Logger(ctx).Trace().Str("queue", s.config.Queue).Msgf("read message")
		s.storedMessages.Set(messageID, msg)

		return rec, nil
	}
}

// metadataFromMsg extracts all the present headers from a stomp.Message into
// opencdc.Metadata.
func metadataFromMsg(msg *stomp.Message) opencdc.Metadata {
	metadata := make(opencdc.Metadata)

	for i := range msg.Header.Len() {
		k, v := msg.Header.GetAt(i)

		// Prefix to avoid collisions with other metadata keys
		headerKey := "activemq.header." + k

		// According to the STOMP protocol, headers can have multiple values for
		// the same key. We concatenate them with a comma and a space.
		if headerVal, ok := metadata[headerKey]; ok {
			var sb strings.Builder
			sb.Grow(len(headerVal) + len(v) + 2)
			sb.WriteString(headerVal)
			sb.WriteString(", ")
			sb.WriteString(v)

			metadata[headerKey] = sb.String()
		} else {
			metadata[headerKey] = v
		}
	}

	return metadata
}

func (s *Source) Ack(ctx context.Context, position opencdc.Position) error {
	pos, err := parseSDKPosition(position)
	if err != nil {
		return fmt.Errorf("failed to parse position: %w", err)
	}

	msg, ok := s.storedMessages.Get(pos.MessageID)
	if !ok {
		return fmt.Errorf("message with ID %q not found", pos.MessageID)
	}

	if err := s.conn.Ack(msg); err != nil {
		return fmt.Errorf("failed to ack message: %w", err)
	}

	s.storedMessages.Remove(pos.MessageID)

	sdk.Logger(ctx).Trace().Str("queue", s.config.Queue).Msgf("acked message")

	return nil
}

func (s *Source) Teardown(ctx context.Context) error {
	return teardown(ctx, s.subscription, s.conn)
}

type Position struct {
	MessageID string `json:"message_id"`
	Queue     string `json:"queue"`
}

func parseSDKPosition(sdkPos opencdc.Position) (Position, error) {
	decoder := json.NewDecoder(bytes.NewBuffer(sdkPos))
	decoder.DisallowUnknownFields()

	var p Position
	if err := decoder.Decode(&p); err != nil {
		return p, fmt.Errorf("failed to parse position: %w", err)
	}

	return p, nil
}

func (p Position) ToSdkPosition() opencdc.Position {
	bs, err := json.Marshal(p)
	if err != nil {
		// this should never happen
		panic(err)
	}

	return opencdc.Position(bs)
}
