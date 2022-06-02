## akv-entrypoint

This program is intended to be used as a Docker entrypoint. It fetches all secrets from the
specified Azure KeyVault and presents them as environment variable for the program.

Authenticanion takes whatever `azidentity.NewDefaultAzureCredential` can, see the Go SDK docs, but generally it should just work.

It works the following way:

1. List all secrets from the keyvault.
1. Fetch the secrets in the Azure listing order.
1. If a secret has `application/json` content-type, then treat it as structure and merge it with the resulting structure.
1. Otherwise, take the secret name as the env var name and take the value.
1. Merging overrides the previous values.

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
