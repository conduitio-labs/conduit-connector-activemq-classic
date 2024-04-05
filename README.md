# Conduit Connector ActiveMQ Classic

The ActiveMQ Classic connector is one of [Conduit](https://conduit.io) plugins. The connector provides both a source and a destination connector for [ActiveMQ Classic](https://activemq.apache.org/components/classic/).

It uses the [stomp protocol](https://stomp.github.io/) to connect to ActiveMQ.

## What data does the OpenCDC record consist of?

| Field                   | Description                                                                 |
|-------------------------|-----------------------------------------------------------------------------|
| `record.Position`       | json object with the queue name and the messageId frame header.             |
| `record.Operation`      | currently fixed as "create".                                                |
| `record.Metadata`       | a string to string map, with keys prefixed as `activemq.header.{STOMP_HEADER_NAME}`. |
| `record.Key`            | the messageId frame header.                                                 |
| `record.Payload.Before` | <empty>                                                                     |
| `record.Payload.After`  | the message body                                                            |

## How to build?
Run `make build` to build the connector.

## Testing
Run `make test` to run all the tests. The command will handle starting and stopping docker containers for you.


## Configuration

Both the source and destination connectors share these configuration parameters:

| name | description | required | default value |
| ---- | ----------- | -------- | ------------- |
| `url` | URL of the ActiveMQ classic broker. | true |  |
| `queue` | Name of the queue to read from or write to. | true |  |
| `user` | Username to use when connecting to the broker. | true |  |
| `password` | Password to use when connecting to the broker. | true |  |
| `sendTimeoutHeartbeat` | Specifies the maximum amount of time between the client sending heartbeat notifications from the server. | true | 2s (*) |
| `recvTimeoutHeartbeat` | Specifies the minimum amount of time between the client expecting to receive heartbeat notifications from the server. | true | 2s (*) |
| `tls.enabled` | Flag to enable or disable TLS. | false | `false` |
| `tls.clientKeyPath` | Path to the client key file. | false |  |
| `tls.clientCertPath` | Path to the client certificate file. | false |  |
| `tls.caCertPath` | Path to the CA certificate file. | false |  |
| `tls.insecureSkipVerify` | Flag to skip verification of the server's certificate chain and host name | false |  |

(*) Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".


### Source configuration

| name | description | required | default value |
| ---- | ----------- | -------- | ------------- |
| `clientID`                   | Specifies the JMS clientID which is used in combination with the activemq.subscriptionName to denote a durable subscriber. | No       |               |
| `dispatchAsync`              | Specifies whether messages should be dispatched synchronously or asynchronously from the producer thread for non-durable topics in the broker. | No       | |
| `exclusive`                  | Indicates the desire to be the sole consumer from the queue.                                    | No       | |
| `maximumPendingMessageLimit` | Specifies the upper limit of pending messages allowed for slow consumers on non-durable topics. | No       |               |
| `noLocal`                    | Indicates if messages sent from the local connection should be excluded from subscriptions.    | No       | |
| `prefetchSize`               | Determines the maximum number of messages to dispatch to the client before it acknowledges a message. | No       |               |
| `priority`                   | Specifies the consumer's priority level for weighted dispatching order.                        | No       |               |
| `retroactive`                | If set to true, makes the subscription retroactive for non-durable topics.                      | No       | |
| `subscriptionName`           | Specifies the name used for durable topic subscriptions.                                        | No       |               |
| `selector`                   | Defines a JMS Selector employing SQL 92 syntax as delineated in the JMS 1.1 specification, enabling a filter to be applied on each message associated with the subscription. | No |               |


### Destination configuration

There are no specific destination parameters as of writing.


## Example pipeline.yml file

Here's an example of a `pipeline.yml` file using `file to activemq classic` and `activemq classic to file` pipelines: 

```yaml
version: 2.0
pipelines:
  - id: file-to-activemq-classic
    status: running
    connectors:
      - id: file.in
        type: source
        plugin: builtin:file
        name: file-destination
        settings:
          path: ./file.in
      - id: activemq-classic.out
        type: destination
        plugin: standalone:activemq
        name: activemq-classic-source
        settings:
          url: localhost:61613
          user: admin
          password: admin
          queue: demo-queue
          sdk.record.format: template
          sdk.record.format.options: '{{ printf "%s" .Payload.After }}'

  - id: activemq-classic-to-file
    status: running
    connectors:
      - id: activemq-classic.in
        type: source
        plugin: standalone:activemq
        name: activemq-classic-source
        settings:
          url: localhost:61613
          user: admin
          password: admin
          queue: demo-queue

      - id: file.out
        type: destination
        plugin: builtin:file
        name: file-destination
        settings:
          path: ./file.out
          sdk.record.format: template
          sdk.record.format.options: '{{ printf "%s" .Payload.After }}'
```
