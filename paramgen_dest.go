// Code generated by paramgen. DO NOT EDIT.
// Source: github.com/ConduitIO/conduit-commons/tree/main/paramgen

package activemq

import (
	"github.com/conduitio/conduit-commons/config"
)

const (
	DestinationConfigPassword              = "password"
	DestinationConfigQueue                 = "queue"
	DestinationConfigRecvTimeoutHeartbeat  = "recvTimeoutHeartbeat"
	DestinationConfigSendTimeoutHeartbeat  = "sendTimeoutHeartbeat"
	DestinationConfigTlsCaCertPath         = "tls.caCertPath"
	DestinationConfigTlsClientCertPath     = "tls.clientCertPath"
	DestinationConfigTlsClientKeyPath      = "tls.clientKeyPath"
	DestinationConfigTlsEnabled            = "tls.enabled"
	DestinationConfigTlsInsecureSkipVerify = "tls.insecureSkipVerify"
	DestinationConfigUrl                   = "url"
	DestinationConfigUser                  = "user"
)

func (DestinationConfig) Parameters() map[string]config.Parameter {
	return map[string]config.Parameter{
		DestinationConfigPassword: {
			Default:     "",
			Description: "Password is the password to use when connecting to the broker.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{
				config.ValidationRequired{},
			},
		},
		DestinationConfigQueue: {
			Default:     "",
			Description: "Queue is the name of the queue to write to.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{
				config.ValidationRequired{},
			},
		},
		DestinationConfigRecvTimeoutHeartbeat: {
			Default:     "2s",
			Description: "RecvTimeoutHeartbeat specifies the minimum amount of time between the\nclient expecting to receive heartbeat notifications from the server",
			Type:        config.ParameterTypeDuration,
			Validations: []config.Validation{},
		},
		DestinationConfigSendTimeoutHeartbeat: {
			Default:     "2s",
			Description: "SendTimeoutHeartbeat specifies the maximum amount of time between the\nclient sending heartbeat notifications to the server",
			Type:        config.ParameterTypeDuration,
			Validations: []config.Validation{},
		},
		DestinationConfigTlsCaCertPath: {
			Default:     "",
			Description: "CaCertPath is the path to the CA certificate file.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{},
		},
		DestinationConfigTlsClientCertPath: {
			Default:     "",
			Description: "ClientCertPath is the path to the client certificate file.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{},
		},
		DestinationConfigTlsClientKeyPath: {
			Default:     "",
			Description: "ClientKeyPath is the path to the client key file.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{},
		},
		DestinationConfigTlsEnabled: {
			Default:     "false",
			Description: "Enabled is a flag to enable or disable TLS.",
			Type:        config.ParameterTypeBool,
			Validations: []config.Validation{},
		},
		DestinationConfigTlsInsecureSkipVerify: {
			Default:     "false",
			Description: "InsecureSkipVerify is a flag to skip verification of the server's\ncertificate chain and host name.",
			Type:        config.ParameterTypeBool,
			Validations: []config.Validation{},
		},
		DestinationConfigUrl: {
			Default:     "",
			Description: "URL is the URL of the ActiveMQ classic broker.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{
				config.ValidationRequired{},
			},
		},
		DestinationConfigUser: {
			Default:     "",
			Description: "User is the username to use when connecting to the broker.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{
				config.ValidationRequired{},
			},
		},
	}
}
