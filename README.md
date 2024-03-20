# Conduit Connector ActiveMQ Classic

The ActiveMQ Classic connector is one of [Conduit](https://conduit.io) plugins. The connector provides both a source and a destination connector for [ActiveMQ Classic](https://activemq.apache.org/components/classic/).

It uses the [stomp protocol](https://stomp.github.io/) to connect to ActiveMQ.

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
| `tlsConfig.useTLS` | Flag to enable or disable TLS. | false | `false` |
| `tlsConfig.clientKeyPath` | Path to the client key file. | false |  |
| `tlsConfig.clientCertPath` | Path to the client certificate file. | false |  |
| `tlsConfig.caCertPath` | Path to the CA certificate file. | false |  |

(*) Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

### Destination configuration

The destination connector accepts an optional `contentType` configuration parameter:

| name | description | required | default value |
| ---- | ----------- | -------- | ------------- |
| `contentType` | Content type of the message. | false | `text/plain` |



Example of a `pipeline.yml` file using `file to activemq classic` and `activemq classic to file` pipelines: 

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
