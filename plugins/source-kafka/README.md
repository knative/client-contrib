## kn-source_kafka

`kn-source_kafka` is a plugin of Knative Client, which allows you to management of Kafka event source interactively from the command line.

### Description

`kn-source_kafka` is a plugin of Knative Client. You could create, update, describe and delete Kafka event sources in Knative Eventing. Go to [Knative Eventing document](https://knative.dev/docs/eventing/samples/kafka/source/) to understand more about Kafka event sources.

### Build and Install

You must [set up your development environment](https://github.com/knative/client/blob/master/docs/DEVELOPMENT.md#prerequisites) before you build `kn-source_kafka`.

**Building:**

Once you've set up your development environment, let's build `kn-source_kafka`. Run below command under the directory of `client-contrib/plugins/source-kafka`.

```sh
$ hack/build.sh
```

**Installing:**

You will get an excuatable file `kn-source_kafka` under the directory of `client-contrib/plugins/source-kafka` after you run the build command. Then let's install it to become a Knative Client `kn` plugin.

Install a plugin by simply copying the excuatable file `kn-source_kafka` to the folder of the `kn` plugins directory. You will be able to invoke it by `kn source_kafka`.

### Usage

```
$ kn source_kafka
Manage your Knative Kafka eventing sources

Usage:
  kafka [command]

Available Commands:
  create      create NAME
  delete      delete NAME
  describe    describe NAME
  help        Help about any command
  update      update NAME

Flags:
  -h, --help   help for kafka

Use "kafka [command] --help" for more information about a command.
```

#### `kn source_kafka create`

```
$ kn source_kafka create --help
create NAME

Usage:
  kafka create NAME [flags]

Examples:
#Creates a new Kafka source with mykafkasrc which subscribes a Kafka server 'my-cluster-kafka-bootstrap.kafka.svc:9092' at topic 'test-topic' using the consumer group ID 'test-consumer-group' and sends the event messages to service 'event-display'
kn source_kafka create mykafkasrc --servers my-cluster-kafka-bootstrap.kafka.svc:9092 --topics test-topic --consumergroup test-consumer-group --sink svc:event-display

Flags:
  -A, --all-namespaces         If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
      --consumergroup string   the consumer group ID
  -h, --help                   help for create
  -n, --namespace string       Specify the namespace to operate in.
      --servers string         Kafka bootstrap servers that the consumer will connect to, consist of a hostname plus a port pair, e.g. my-kafka-bootstrap.kafka:9092
  -s, --sink string            Addressable sink for events
      --topics string          Topics to consume messages from
```
