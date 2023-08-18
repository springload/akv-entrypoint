## Why would I use it?

There are multiple products on Azure Cloud that allow running containers. The problem is they're all configured differently.

1. Web Apps for Containers have full support for referencing secrets from Azure Key Vault.
2. Azure Container Apps will have full support for referencing secrets at some point (this functionality is in preview as of Aug 2023)
   but at the moment it's only ACA-native secrets.
3. Azure Container Instances don't support KeyVault at all and have their own secrets.

For your workload you might have to use several products, i.e. Web Apps only support applications with web endpoints.
For workloads without a web endpoint you might want to use ACA. One-off jobs in ACA are in preview (as of Aug 2023) so
you might have to use ACI.

This middleware solves that problem. It runs as an entrypoint to your container application. It fetches all (or specified) secrets from the
specified Azure KeyVault and presents them as environment variables to the program.

Authentication takes whatever `azidentity.NewDefaultAzureCredential` can, see the Go SDK docs, but generally it should just work.

Note: for ACA when using UAI, you have to specify `AZURE_CLIENT_ID` as per https://github.com/microsoft/azure-container-apps/issues/325/

It works the following way:

1. If `AKV_ENTRYPOINT_SECRET_NAMES` env var is not set list all secrets from the KeyVault
1. Fetche the secrets in the Azure listing order.
1. If a secret has `application/json` content-type, then treat it as structure and merge it with the resulting structure.
1. Otherwise, take the secret name as the env var name and take the value.
1. Merging overrides the previous values.

When `AKV_ENTRYPOINT_SECRET_NAMES` is set `akv-entrypoint` only fetches the specified secret names, which are comma-separated.
This way it can't know in advance secrets content-type so it must be hinted if it should treat it as a json, for example:

```bash

AKV_ENTRYPOINT_SECRET_NAMES=json:secret1,secret2

```
It'll treat secret1 as JSON and secret2 as plain text.

Note: KeyVault secret names can't have underscores so it's normal to have dashes in them as separators.
At the same time environment variables in Linux can't have dashes. So the middleware will replace all dashes
with underscores for all secret names, i.e. if you have a KeyVault secret `mega-secret` the environment name you'll
get will become `mega_secret` to conform to the spec.


## Examples

### CLI usage for testing

```bash
$akv-entrypoint -d print -n your-keyvault # to see env vars printed in JSON
$akv-entrypoint -d run -n your-keyvault -- your-app # to run the app
```

### Docker entrypoint

```
# somewhere above
FROM ghcr.io/springload/akv-entrypoint:latest as akv-entrypoint
...
# later on in your app stage
COPY --from=akv-entrypoint /usr/bin/akv-entrypoint /usr/bin/akv-entrypoint
...
# at the end of your production stage
ENTRYPOINT ["/usr/bin/akv-entrypoint", "run", "-n", "your-keyvault",  "--"]
CMD ["/usr/local/bin/gunicorn", "--config", "/app/docker/gunicorn.py", "app.wsgi" ]
```
