package activemq

import (
	"testing"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
	"github.com/matryer/is"
)

func TestMetadataFromMsg(t *testing.T) {
	is := is.New(t)

	header := &frame.Header{}
	header.Add("key", "value1")
	header.Add("key", "value2")

	header.Add("key2", "value3")
	header.Add("key3", "value4")

	metadata := metadataFromMsg(&stomp.Message{
		Header: header,
	})

	is.Equal(metadata["header-key"], "value1, value2")
	is.Equal(metadata["header-key2"], "value3")
	is.Equal(metadata["header-key3"], "value4")
}
