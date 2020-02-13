# kn-admin
A kn plugin for Knative cluster management
This kn-admin plugin is designed to help administrators and operators better manage a Knative platform installation with kn CLI. 
The plugin’s main objective is to make administration and operation workflows easier, for instance by making it easy to accomplish 
tasks such as feature flags enablement or disablement with one command, instead of many manual steps like modifying ConfigMaps or yaml files.

## Build
To build kn-admin
```bash
$ pwd
/Users/zhanggong/zhgworkspace/knative/src/github.com/knative/client-contrib/plugins/kn-admin
$ ./hack/build.sh
🧹 Format
🚧 Compile
$ ls build/_output
  kn-admin

```
To run kn-admin as plugin, please refer to https://github.com/knative/client/blob/master/docs/README.md#options

## kn-admin Usage
```bash
$ kn admin -h
A plugin of kn client to manage Knative for administrators. 

For example:
kn admin domain set - to set Knative route domain
kn admin private-registry enable - to enable deployment from the private registry

Usage:
  admin [command]

Available Commands:
  domain           Manage route domain
  help             Help about any command
  private-registry Manage private-registry

Flags:
      --config string   config file (default is $HOME/.config/kn/plugins/kn-admin.yaml)
  -h, --help            help for admin
  -t, --toggle          Help message for toggle

Use "admin [command] --help" for more information about a command.
```

## kn-admin Examples
### As a Knative administrator, I want to update Knative route domain with my custom domain.
```bash
# kn admin will update the default route domain if --selector no specified
$ Kn admin domain set --custom-domain mydomain.com
Updated Knative route domain mydomain.com

# Service with a label app=v1 will use test.com
$ Kn admin domain set --custom-domain test.com --selector app=v1
Updated Knative route domain test.com
```
### As a Knative administrator, I want to enable deploying from private registry.
```
$ kn admin private-registry enable \
  --secret-name=[SECRET_NAME]
  --docker-server=[PRIVATE_REGISTRY_SERVER_URL] \
  --docker-email=[PRIVATE_REGISTRY_EMAIL] \
  --docker-username=[PRIVATE_REGISTRY_USER] \
  --docker-password=[PRIVATE_REGISTRY_PASSWORD]
```

