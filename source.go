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
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

type Source struct {
	sdk.UnimplementedSource
	config Config

	conn         *stomp.Conn
	subscription *stomp.Subscription

	storedMessages messageMap
}

type messageMap struct {
	mutex *sync.Mutex
	msgs  map[string]*stomp.Message
}

func newMessageMap() messageMap {
	return messageMap{
		mutex: &sync.Mutex{},
		msgs:  make(map[string]*stomp.Message),
	}
}

func (m *messageMap) add(key string, msg *stomp.Message) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.msgs[key] = msg
}

func (m *messageMap) get(key string) (*stomp.Message, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	msg, ok := m.msgs[key]
	return msg, ok
}

func (m *messageMap) remove(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.msgs, key)
}

func NewSource() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{}, sdk.DefaultSourceMiddleware()...)
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	err := sdk.Util.ParseConfig(cfg, &s.config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	s.config.logConfig(ctx, "configured destination")

	s.storedMessages = newMessageMap()

	return nil
}

func (s *Source) Open(ctx context.Context, sdkPos sdk.Position) (err error) {
	s.conn, err = connect(ctx, s.config)
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

	s.subscription, err = s.conn.Subscribe(s.config.Queue, stomp.AckClientIndividual)
	if err != nil {
		return fmt.Errorf("failed to subscribe to queue: %w", err)
	}

	sdk.Logger(ctx).Debug().Msg("opened source")

	return nil
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	var rec sdk.Record

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
			key      = sdk.RawData(messageID)
			payload  = sdk.RawData(msg.Body)
		)

		rec = sdk.Util.Source.NewRecordCreate(sdkPos, metadata, key, payload)

		sdk.Logger(ctx).Trace().Str("queue", s.config.Queue).Msgf("read message")
		s.storedMessages.add(messageID, msg)

		return rec, nil
	}
}

// metadataFromMsg extracts all the present headers from a stomp.Message into
// sdk.Metadata.
func metadataFromMsg(msg *stomp.Message) sdk.Metadata {
	metadata := make(sdk.Metadata)

	for i := range msg.Header.Len() {
		k, v := msg.Header.GetAt(i)
		metadata["header-"+k] = v
	}

	return metadata
}

func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
	pos, err := parseSDKPosition(position)
	if err != nil {
		return fmt.Errorf("failed to parse position: %w", err)
	}

	msg, ok := s.storedMessages.get(pos.MessageID)
	if !ok {
		return fmt.Errorf("message with ID %q not found", pos.MessageID)
	}

	if err := s.conn.Ack(msg); err != nil {
		return fmt.Errorf("failed to ack message: %w", err)
	}

	s.storedMessages.remove(pos.MessageID)

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
