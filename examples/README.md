# Examples

This directory contains Terraform examples for the SearchStax provider.

## Layout

| Path | Purpose |
|------|---------|
| `provider/provider.tf` | Provider configuration (used on the docs index page) |
| `data-sources/<name>/data-source.tf` | Minimal example for each data source (registry docs) |
| `resources/<name>/resource.tf` | Minimal example for each resource (registry docs) |
| `complete/` | Runnable stack against an **existing** deployment |

The documentation generator (`go generate`) reads `provider.tf`, `data-source.tf`, and `resource.tf` files only. Other files (e.g. `complete/variables.tf`) are for manual testing.

## Provider configuration

Credentials can be set in Terraform or via environment variables:

- `SEARCHSTAX_USERNAME`
- `SEARCHSTAX_PASSWORD`
- `SEARCHSTAX_HOST` (optional, defaults to the SearchStax cloud API)

## Run the complete example

```shell
cd examples/complete
terraform init
terraform plan \
  -var 'ssx_username=you@example.com' \
  -var 'ssx_pwd=secret' \
  -var 'account_name=my_account' \
  -var 'deployment_uid=ss123456'
```

For local development with a provider build override, skip `terraform init` and use `TF_REATTACH_PROVIDERS` as described in the root README.

## Format examples

```shell
terraform fmt -recursive ./examples/
```

This runs automatically via `go generate` in the repository root.
