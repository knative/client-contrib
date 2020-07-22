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
