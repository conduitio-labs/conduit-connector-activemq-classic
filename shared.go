package activemq

import sdk "github.com/conduitio/conduit-connector-sdk"

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
	URL string `json:"url" validate:"required"`
}

type SourceConfig struct {
	Config
	Queue       string `json:"queue" validate:"required"`
	ContentType string `json:"content" default:"text/plain"`
}

type DestinationConfig struct {
	Config
	Queue string `json:"queue" validate:"required"`
}
