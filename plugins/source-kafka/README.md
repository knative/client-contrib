# kn-source-kafka

`kn-source-kafka` is a plugin of Knative Client, for management of kafka event
source interactively from the command line.

## Description

`kn-source-kafka` is a plugin of Knative Client. You can create, describe and
delete kafka event sources. Go to
[Knative Eventing document](https://knative.dev/docs/eventing/samples/kafka/source/)
to understand more about kafka event sources.

## Build and Install

You must
[set up your development environment](https://github.com/knative/client/blob/master/docs/DEVELOPMENT.md#prerequisites)
before you build `kn-source-kafka`.

**Building:**

Once you've set up your development environment, let's build `kn-source-kafka`.
Run below command under the directory of `client-contrib/plugins/source-kafka`.

```sh
$ hack/build.sh
```

**Installing:**

You will get an executable file `kn-source-kafka` under the directory of
`client-contrib/plugins/source-kafka` after you run the build command. Then
let's install it to become a Knative Client `kn` plugin.

Install the plugin by simply copying the executable file `kn-source-kafka` to
the folder of the `kn` plugins directory. You will be able to invoke it by
`kn source kafka`.

## Usage

### kafka

Knative eventing kafka source plugin

#### Synopsis

Manage Knative kafka eventing sources

#### Options

```
  -h, --help   help for kafka
```

#### SEE ALSO

* [kafka create](#kafka-create)	 - Create a kafka source
* [kafka delete](#kafka-delete)	 - Delete a kafka source
* [kafka describe](#kafka-describe)	 - Describe a kafka source

### kafka create

Create a kafka source

#### Synopsis

Create a kafka source

```
kafka create NAME --servers SERVERS --topics TOPICS --consumergroup GROUP --sink SINK [flags]
```

#### Examples

```
# Create a new kafka source 'mykafkasrc' which subscribes a kafka server 'my-cluster-kafka-bootstrap.kafka.svc:9092' at topic 'test-topic' using the consumer group ID 'test-consumer-group' and sends the events to service 'event-display'
kn source kafka create mykafkasrc --servers my-cluster-kafka-bootstrap.kafka.svc:9092 --topics test-topic --consumergroup test-consumer-group --sink svc:event-display
```

#### Options

```
  -A, --all-namespaces         If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
      --consumergroup string   the consumer group ID
  -h, --help                   help for create
  -n, --namespace string       Specify the namespace to operate in.
      --servers string         Kafka bootstrap servers that the consumer will connect to, consist of a hostname plus a port pair, e.g. my-kafka-bootstrap.kafka:9092
  -s, --sink string            Addressable sink for events
      --topics string          Topics to consume messages from
```

#### SEE ALSO

* [kafka](#kafka)	 - Knative eventing kafka source plugin

### kafka delete

Delete a kafka source

#### Synopsis

Delete a kafka source

```
kafka delete NAME [flags]
```

#### Examples

```
# Delete a kafka source with name 'mykafkasrc'
kn source kafka delete mykafkasrc
```

#### Options

```
  -A, --all-namespaces     If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
  -h, --help               help for delete
  -n, --namespace string   Specify the namespace to operate in.
```

#### SEE ALSO

* [kafka](#kafka)	 - Knative eventing kafka source plugin

### kafka describe

Describe a kafka source

#### Synopsis

Describe a kafka source

```
kafka describe NAME [flags]
```

#### Examples

```
# Describe a kafka source with NAME
kn source kafka describe kafka-name
```

#### Options

```
  -A, --all-namespaces     If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
  -h, --help               help for describe
  -n, --namespace string   Specify the namespace to operate in.
```

#### SEE ALSO

* [kafka](#kafka)	 - Knative eventing kafka source plugin

## More information
	
* [Knative Client](https://github.com/knative/client)
* [How to contribute a plugin](https://github.com/knative/client-contrib#how-to-contribute-a-plugin)

