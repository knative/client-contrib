# github

Knative Client plugin `github`

## Usage

### github

Knative eventing GitHub source plugin

#### Synopsis

Manage your Knative GitHub eventing sources

#### Options

```
  -h, --help   help for github
```

#### SEE ALSO

* [github create](#github-create)	 - create NAME
* [github delete](#github-delete)	 - delete NAME
* [github describe](#github-describe)	 - describe NAME
* [github update](#github-update)	 - update NAME

### github create

create NAME

#### Synopsis

create a GitHub source

```
github create NAME [flags]
```

#### Examples

```
# Creates a new GitHub source with NAME using credentials
kn source github create NAME  --access-token $MY_ACCESS_TOKEN --secret-token $MY_SECRET_TOKEN

# Creates a new GitHub source with NAME with specified organization and repository using credentials
kn source github create NAME --org knative --repo client-contrib --access-token $MY_ACCESS_TOKEN --secret-token $MY_SECRET_TOKEN
```

#### Options

```
      --access-token string   The GitHub access-token to use
  -A, --all-namespaces        If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
      --api-url string        The GitHub API URL to use (default "https://api.github.com")
  -h, --help                  help for create
  -n, --namespace string      Specify the namespace to operate in.
      --org string            The GitHub organization or username
      --repo string           Repository name to consume messages from
      --secret-token string   The GitHub secret-token to use
  -s, --sink string           Addressable sink for events
```

#### SEE ALSO

* [github](#github)	 - Knative eventing GitHub source plugin

### github delete

delete NAME

#### Synopsis

delete a GitHub source

```
github delete NAME [flags]
```

#### Examples

```
# Deletes a GitHub source with NAME
kn source github delete NAME
```

#### Options

```
  -A, --all-namespaces     If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
  -h, --help               help for delete
  -n, --namespace string   Specify the namespace to operate in.
```

#### SEE ALSO

* [github](#github)	 - Knative eventing GitHub source plugin

### github describe

describe NAME

#### Synopsis

update a GitHub source

```
github describe NAME [flags]
```

#### Examples

```
# Describes a GitHub source with NAME
kn source github describe NAME
```

#### Options

```
  -A, --all-namespaces     If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
  -h, --help               help for describe
  -n, --namespace string   Specify the namespace to operate in.
```

#### SEE ALSO

* [github](#github)	 - Knative eventing GitHub source plugin

### github update

update NAME

#### Synopsis

update a GitHub source

```
github update NAME [flags]
```

#### Examples

```
# Updates a GitHub source with NAME
kn source github update NAME
```

#### Options

```
      --access-token string   The GitHub access-token to use
  -A, --all-namespaces        If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
      --api-url string        The GitHub API URL to use (default "https://api.github.com")
  -h, --help                  help for update
  -n, --namespace string      Specify the namespace to operate in.
      --org string            The GitHub organization or username
      --repo string           Repository name to consume messages from
      --secret-token string   The GitHub secret-token to use
  -s, --sink string           Addressable sink for events
```

#### SEE ALSO

* [github](#github)	 - Knative eventing GitHub source plugin

## More information
	
* [Knative Client](https://github.com/knative/client)
* [How to contribute a plugin](https://github.com/knative/client-contrib#how-to-contribute-a-plugin)

