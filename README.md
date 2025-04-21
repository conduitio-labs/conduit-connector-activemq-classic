# Conduit Connector ActiveMQ Classic

The ActiveMQ Classic connector is one of [Conduit](https://conduit.io) plugins. The connector provides both a source and a destination connector for [ActiveMQ Classic](https://activemq.apache.org/components/classic/).

It uses the [stomp protocol](https://stomp.github.io/) to connect to ActiveMQ.

## What data does the OpenCDC record consist of?

| Field                   | Description                                                                          |
| ----------------------- | ------------------------------------------------------------------------------------ |
| `record.Position`       | json object with the queue name and the messageId frame header.                      |
| `record.Operation`      | currently fixed as "create".                                                         |
| `record.Metadata`       | a string to string map, with keys prefixed as `activemq.header.{STOMP_HEADER_NAME}`. |
| `record.Key`            | the messageId frame header.                                                          |
| `record.Payload.Before` | <empty>                                                                              |
| `record.Payload.After`  | the message body                                                                     |

## How to build?

Run `make build` to build the connector.

## Testing

Run `make test` to run all the tests. The command will handle starting and stopping docker containers for you.

## Source configuration

<!-- readmegen:source.parameters.yaml -->
```yaml
version: 2.2
pipelines:
  - id: example
    status: running
    connectors:
      - id: example
        plugin: "activemq"
        settings:
          # The password to use when connecting to the broker.
          # Type: string
          # Required: yes
          password: ""
          # The name of the queue to write to.
          # Type: string
          # Required: yes
          queue: ""
          # The URL of the ActiveMQ classic broker.
          # Type: string
          # Required: yes
          url: ""
          # The username to use when connecting to the broker.
          # Type: string
          # Required: yes
          user: ""
          # Whether messages should be dispatched synchronously or
          # asynchronously from the producer thread for non-durable topics in
          # the broker. Maps to the activemq.dispatchAsync header.
          # Type: bool
          # Required: no
          activemq.dispatchAsync: "false"
          # Whether the desire to be the sole consumer from the queue. Maps to
          # the activemq.exclusive header.
          # Type: bool
          # Required: no
          activemq.exclusive: "false"
          # The upper limit of pending messages allowed for slow consumers on
          # non-durable topics. When this limit is reached, older messages will
          # be discarded to handle slow consumer backlog. Maps to the
          # activemq.maximumPendingMessageLimit header.
          # Type: int
          # Required: no
          activemq.maximumPendingMessageLimit: "0"
          # Whether messages sent from the local connection should be excluded
          # from subscriptions. When set to true, locally sent messages will be
          # ignored. Maps to the activemq.noLocal header.
          # Type: bool
          # Required: no
          activemq.noLocal: "false"
          # The maximum number of messages to dispatch to the client before it
          # acknowledges a message. No further messages are dispatched once this
          # limit is hit. For fair message distribution across consumers,
          # consider setting this to a value greater than 1. Maps to the
          # activemq.prefetchSize header.
          # Type: int
          # Required: no
          activemq.prefetchSize: "0"
          # The consumer's priority level for weighted dispatching order. Maps
          # to the activemq.priority header.
          # Type: string
          # Required: no
          activemq.priority: ""
          # Whether the subscription is retroactive for non-durable topics. Maps
          # to the activemq.retroactive header.
          # Type: bool
          # Required: no
          activemq.retroactive: "false"
          # The name used for durable topic subscriptions. Prior to ActiveMQ
          # version 5.7.0, both clientID on the connection and subscriptionName
          # on the subscribe operation must match. Maps to the
          # activemq.subscriptionName header.
          # Type: string
          # Required: no
          activemq.subscriptionName: ""
          # The JMS clientID which is used in combination with the
          # activemq.subcriptionName to denote a durable subscriber. Maps to the
          # client-id header.
          # Type: string
          # Required: no
          clientID: ""
          # The minimum amount of time between the client expecting to receive
          # heartbeat notifications from the server
          # Type: duration
          # Required: no
          recvTimeoutHeartbeat: "2s"
          # A JMS Selector employing SQL 92 syntax as delineated in the JMS 1.1
          # specification, enabling a filter to be applied on each message
          # associated with the subscription. Maps to the selector header.
          # Type: string
          # Required: no
          selector: ""
          # The maximum amount of time between the client sending heartbeat
          # notifications to the server
          # Type: duration
          # Required: no
          sendTimeoutHeartbeat: "2s"
          # The path to the CA certificate file.
          # Type: string
          # Required: no
          tls.caCertPath: ""
          # The path to the client certificate file.
          # Type: string
          # Required: no
          tls.clientCertPath: ""
          # The path to the client key file.
          # Type: string
          # Required: no
          tls.clientKeyPath: ""
          # Flag to enable or disable TLS.
          # Type: bool
          # Required: no
          tls.enabled: "false"
          # Flag to skip verification of the server's certificate chain and host
          # name.
          # Type: bool
          # Required: no
          tls.insecureSkipVerify: "false"
          # Maximum delay before an incomplete batch is read from the source.
          # Type: duration
          # Required: no
          sdk.batch.delay: "0"
          # Maximum size of batch before it gets read from the source.
          # Type: int
          # Required: no
          sdk.batch.size: "0"
          # Specifies whether to use a schema context name. If set to false, no
          # schema context name will be used, and schemas will be saved with the
          # subject name specified in the connector (not safe because of name
          # conflicts).
          # Type: bool
          # Required: no
          sdk.schema.context.enabled: "true"
          # Schema context name to be used. Used as a prefix for all schema
          # subject names. If empty, defaults to the connector ID.
          # Type: string
          # Required: no
          sdk.schema.context.name: ""
          # Whether to extract and encode the record key with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.key.enabled: "true"
          # The subject of the key schema. If the record metadata contains the
          # field "opencdc.collection" it is prepended to the subject name and
          # separated with a dot.
          # Type: string
          # Required: no
          sdk.schema.extract.key.subject: "key"
          # Whether to extract and encode the record payload with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.payload.enabled: "true"
          # The subject of the payload schema. If the record metadata contains
          # the field "opencdc.collection" it is prepended to the subject name
          # and separated with a dot.
          # Type: string
          # Required: no
          sdk.schema.extract.payload.subject: "payload"
          # The type of the payload schema.
          # Type: string
          # Required: no
          sdk.schema.extract.type: "avro"
```
<!-- /readmegen:source.parameters.yaml -->

## Destination configuration

<!-- readmegen:destination.parameters.yaml -->
```yaml
version: 2.2
pipelines:
  - id: example
    status: running
    connectors:
      - id: example
        plugin: "activemq"
        settings:
          # The password to use when connecting to the broker.
          # Type: string
          # Required: yes
          password: ""
          # The name of the queue to write to.
          # Type: string
          # Required: yes
          queue: ""
          # The URL of the ActiveMQ classic broker.
          # Type: string
          # Required: yes
          url: ""
          # The username to use when connecting to the broker.
          # Type: string
          # Required: yes
          user: ""
          # The minimum amount of time between the client expecting to receive
          # heartbeat notifications from the server
          # Type: duration
          # Required: no
          recvTimeoutHeartbeat: "2s"
          # The maximum amount of time between the client sending heartbeat
          # notifications to the server
          # Type: duration
          # Required: no
          sendTimeoutHeartbeat: "2s"
          # The path to the CA certificate file.
          # Type: string
          # Required: no
          tls.caCertPath: ""
          # The path to the client certificate file.
          # Type: string
          # Required: no
          tls.clientCertPath: ""
          # The path to the client key file.
          # Type: string
          # Required: no
          tls.clientKeyPath: ""
          # Flag to enable or disable TLS.
          # Type: bool
          # Required: no
          tls.enabled: "false"
          # Flag to skip verification of the server's certificate chain and host
          # name.
          # Type: bool
          # Required: no
          tls.insecureSkipVerify: "false"
          # Maximum delay before an incomplete batch is written to the
          # destination.
          # Type: duration
          # Required: no
          sdk.batch.delay: "0"
          # Maximum size of batch before it gets written to the destination.
          # Type: int
          # Required: no
          sdk.batch.size: "0"
          # Allow bursts of at most X records (0 or less means that bursts are
          # not limited). Only takes effect if a rate limit per second is set.
          # Note that if `sdk.batch.size` is bigger than `sdk.rate.burst`, the
          # effective batch size will be equal to `sdk.rate.burst`.
          # Type: int
          # Required: no
          sdk.rate.burst: "0"
          # Maximum number of records written per second (0 means no rate
          # limit).
          # Type: float
          # Required: no
          sdk.rate.perSecond: "0"
          # The format of the output record. See the Conduit documentation for a
          # full list of supported formats
          # (https://conduit.io/docs/using/connectors/configuration-parameters/output-format).
          # Type: string
          # Required: no
          sdk.record.format: "opencdc/json"
          # Options to configure the chosen output record format. Options are
          # normally key=value pairs separated with comma (e.g.
          # opt1=val2,opt2=val2), except for the `template` record format, where
          # options are a Go template.
          # Type: string
          # Required: no
          sdk.record.format.options: ""
          # Whether to extract and decode the record key with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.key.enabled: "true"
          # Whether to extract and decode the record payload with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.payload.enabled: "true"
```
<!-- /readmegen:destination.parameters.yaml -->

## NOTES

- The source `activemq.subscriptionName` parameter is only supported for
  activemq classic v5.0. When using this connector with previous versions of
  activemq, this parameter will be ignored, as the previous header name for
  this parameter was `activemq.subcriptionName`.
