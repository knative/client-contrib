## kn-curl

This plugins makes it easy to curl a Knative service.

### Assumptions

1.  Your Knative installation uses Istio for its mesh networking.
2.  You follow the Dependencies listed below.
3.  You don't have a custom domain for your services.

### Dependencies

1. Get the latest [`kn`](https://github.com/knative/client) binary for your
   platform.
2. Make sure the `kn-curl` BASH file is located in your `~/.config/kn/plugins`
   directory. Of you are using a custom plugins repository, then add it there.
3. Make sure the `kn-curl` BASH file has executable priviledges. You can achieve
   this with `chmod +x kn-curl`

### Examples

```bash
# curl the service named hello
kn curl hello
```

### Next steps

1. Remove Istio assumption or at least deal with cases of non-istio Knative
   deployments.
2. Determine the domain and use that or default `example.com` if not set.
3. Add e2e test for this BASH plugin.
