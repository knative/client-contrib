## Knative Client Contributions

This repository is the place for curated contributions for the [Knative client](https://github.com/knative/client), especially Knative plugins


### Plugins

All plugins are stored below the `plugins/` directory. Currently you can find the following plugins there:

#### kn-hello

[kn-hello](plugin/hello) is a "Hello World" plugin which also serves as a blueprint for new plugins.  It contains the pieces that a mandatory for any new (golang based) plugin.

I.e it demonstrates:

* How the package structure should look like
* How the README and documentation should be structured
* How CI and testing in general can be setup
* Contains a sample build tool `build.sh` which can easily be customized

_list of plugins to be continued ..._

### How to contribute a plugin

First of all, thank you for considering to contribute a `kn` plugin. That's really so awesome, we love contributions !

Before you starting to craft a pull request, please consider to perform the following step:

* Enter the Knative [#cli](https://knative.slack.com/archives/CE4MVFVAQ) slack and discuss you plugin idea there first.
* When creating a PR please follow the following process:
  - Copy over the plugin from the directory `plugins/hello` to a new directory with your plugins short name (i.e. the command name), also below `plugins/`. E.g. `plugins/awesome` if you are about to create a `kn-awesome` plugin enriching kn with a `kn awesome` command.
  - Put your code in the `pkg/` and `cmd/` directories, similar to the existing code.
  - Check and adapt the top-level comments and variables in `hack/build.sh` to reflect your plugins settings.
  - Adapt and add to the `README.md` which is supposed to server as a full documentation for your plugin. See the example section given there and replace it with the content for your plugin. Especially a reference documentation is required as well as an example section.
  - Adapt the approved in the `OWNERS` file to fill in the maintainers of your plugin, but also leave the existing folks in that file (as a fallback)
  - Be sure that your plugin is completely **self-contained**. I.e. it must not depend on anything above its directory and should setup its own dependencies.
* Work with the reviewers to get your PR integrated.

Plese note that all plugins in this repository share the same release cycle and release cadence, which is currently 6 weeks together with the Knative client release.
