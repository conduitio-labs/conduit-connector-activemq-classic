package activemq

import (
	"encoding/json"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/go-stomp/stomp"
	"github.com/go-stomp/stomp/frame"
)

// version is set during the build process with ldflags (see Makefile).
// Default version matches default from runtime/debug.
var version = "(devel)"

// Specification returns the connector's specification.
func Specification() sdk.Specification {
	return sdk.Specification{
		Name:        "activemq",
		Summary:     "An ActiveMQ classic source and destination plugin for Conduit, written in Go.",
		Description: "An ActiveMQ classic source and destination plugin for Conduit, written in Go.",
		Version:     version,
		Author:      "Meroxa, Inc.",
	}
}

// Connector combines all constructors for each plugin in one struct.
var Connector = sdk.Connector{
	NewSpecification: Specification,
	NewSource:        NewSource,
	NewDestination:   NewDestination,
}

//go:generate paramgen -output=paramgen_src.go SourceConfig
//go:generate paramgen -output=paramgen_dest.go DestinationConfig

type Config struct {
	// URL is the URL of the ActiveMQ classic broker.
	URL string `json:"url" validate:"required"`

	// User is the username to use when connecting to the broker.
	User string `json:"user" validate:"required"`

	// Password is the password to use when connecting to the broker.
	Password string `json:"password" validate:"required"`

	TLSConfig `json:"tlsConfig"`
}

type TLSConfig struct {
	// UseTLS is a flag to enable or disable TLS.
	UseTLS bool `json:"useTLS" default:"false"`

	// ClientKeyPath is the path to the client key file.
	ClientKeyPath string `json:"clientKeyPath" validate:"required"`

	// ClientCertPath is the path to the client certificate file.
	ClientCertPath string `json:"clientCertPath" validate:"required"`

	// CaCertPath is the path to the CA certificate file.
	CaCertPath string `json:"caCertPath" validate:"required"`
}

type SourceConfig struct {
	Config
	// Queue is the name of the queue to read from.
	Queue string `json:"queue" validate:"required"`

	// ContentType is the content type of the message.
	ContentType string `json:"contentType" default:"text/plain"`
}

type DestinationConfig struct {
	Config

	// Queue is the name of the queue to write to.
	Queue string `json:"queue" validate:"required"`
	// ContentType is the content type of the message.
	ContentType string `json:"contentType" default:"text/plain"`
}

type Position struct {
	MessageID string `json:"message_id"`
	Queue     string `json:"queue"`
}

func NewPosition(msg *stomp.Message) Position {
	messageID := msg.Header.Get(frame.MessageId)

	return Position{
		MessageID: messageID,
		Queue:     msg.Destination,
	}
}

func parseSDKPosition(sdkPos sdk.Position) (Position, error) {
	var p Position
	err := json.Unmarshal([]byte(sdkPos), &p)
	return p, err
}

func (p Position) ToSdkPosition() sdk.Position {
	bs, err := json.Marshal(p)
	if err != nil {
		// this should never happen
		panic(err)
	}

	return sdk.Position(bs)
}

func (p Position) toMsg(s *Source) *stomp.Message {
	var header frame.Header
	header.Add(frame.MessageId, p.MessageID)

	return &stomp.Message{Header: &header, Destination: p.Queue}
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
