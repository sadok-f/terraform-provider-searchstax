# Terraform Provider SearchStax
<div align="center">
    <img src="https://www.searchstax.com/docs/wp-content/themes/docs/images/logo.svg" width="400" alt="SearchStax" />
    <br/>
   <a href="https://www.searchstax.com/docs/searchstax-cloud-apis-overview/">
    <img src="https://img.shields.io/static/v1?label=Docs&message=API&color=000000&style=for-the-badge"  alt="SearchStax Cloud API Documentation"/>
    </a>

[![Tests](https://github.com/sadok-f/terraform-provider-searchstax/actions/workflows/test.yml/badge.svg)](https://github.com/sadok-f/terraform-provider-searchstax/actions/workflows/test.yml)
</div>

A Terraform provider for [SearchStax Cloud](https://www.searchstax.com/docs/searchstax-cloud-docs-home/) that lets you
manage deployments, users, backups, custom JARs, DNS records and more as code.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25

## Building SearchStax Provider locally

1. Clone the repository
```shell
git clone git@github.com:sadok-f/terraform-provider-searchstax.git
```
2. Enter the repository directory
```shell
cd terraform-provider-searchstax/
```
3. Build the provider using the Go `install` command:

```shell
go install
```

Terraform allows you to use local provider builds by setting a `dev_overrides` block in a configuration file called `.terraformrc`. This block overrides all other configured installation methods.

Terraform searches for the `.terraformrc` file in your home directory and applies any configuration settings you set.

```
provider_installation {

  dev_overrides {
      "registry.terraform.io/sadok-f/searchstax" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Your `<PATH>` may vary depending on how your Go environment variables are configured. Execute `go env GOBIN` to set it, then set the `<PATH>` to the value returned. If nothing is returned, set it to the default location, `$HOME/go/bin`.

> Skip terraform init when using provider development overrides. It is not necessary and may error unexpectedly.

## Using the provider

The provider is published on the [Terraform Registry](https://registry.terraform.io/providers/sadok-f/searchstax/latest).
Add it to your configuration and run `terraform init`:

```hcl
terraform {
  required_providers {
    searchstax = {
      source  = "sadok-f/searchstax"
      version = "~> 0.1"
    }
  }
}

provider "searchstax" {
  username = var.ssx_username
  password = var.ssx_pwd
}

```
## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

### Acceptance tests

*Note:* Acceptance tests runs on a SearchStax mock api server created by this [Repo](https://github.com/sadok-f/searchstax-mock-api), it returns a mocking response for each API call.

In order to run the full suite of Acceptance tests, run: 

```shell
./internal/provider/acceptance.sh
```

### Testing the provider locally with a debugger
Should you want to validate a change locally, the `--debug` flag allows you to execute the provider against a terraform instance locally.

This also allows for debuggers (e.g. delve) to be attached to the provider.

```sh
go run main.go --debug
# Copy the TF_REATTACH_PROVIDERS env var
# In a new terminal
cd examples/resources/deployment
TF_REATTACH_PROVIDERS=... terraform init
TF_REATTACH_PROVIDERS=... terraform apply
```

More details about running a terraform provider with a debugger:
[https://opencredo.com/blogs/running-a-terraform-provider-with-a-debugger/](https://opencredo.com/blogs/running-a-terraform-provider-with-a-debugger/)

### Implemented Domains

Data sources currently implemented:
- `searchstax_account_backups`
- `searchstax_alert_metrics`
- `searchstax_alerts`
- `searchstax_api_key_deployments`
- `searchstax_auth_token`
- `searchstax_backup_schedules`
- `searchstax_basic_auth`
- `searchstax_custom_jars`
- `searchstax_deployment`
- `searchstax_deployment_api_keys`
- `searchstax_deployment_backups`
- `searchstax_deployment_collections_health`
- `searchstax_deployment_health`
- `searchstax_deployment_servers`
- `searchstax_deployment_server_host_status`
- `searchstax_deployment_users`
- `searchstax_deployments`
- `searchstax_deployments_by_tag`
- `searchstax_dns_record`
- `searchstax_dns_records`
- `searchstax_heartbeats`
- `searchstax_incidents`
- `searchstax_ip_filters`
- `searchstax_plans`
- `searchstax_private_vpc`
- `searchstax_restore_status`
- `searchstax_tags`
- `searchstax_usage`
- `searchstax_usage_extended`
- `searchstax_users`
- `searchstax_webhooks`
- `searchstax_zookeeper_config`
- `searchstax_zookeeper_config_download`
- `searchstax_zookeeper_configs`

Resources currently implemented:
- `searchstax_account_backup`
- `searchstax_alert`
- `searchstax_api_key`
- `searchstax_api_key_association`
- `searchstax_auth_session`
- `searchstax_backup_schedule`
- `searchstax_basic_auth`
- `searchstax_custom_jar`
- `searchstax_deployment`
- `searchstax_deployment_backup`
- `searchstax_deployment_rolling_restart`
- `searchstax_deployment_solr`
- `searchstax_deployment_user`
- `searchstax_dns_record`
- `searchstax_heartbeat`
- `searchstax_ip_filter`
- `searchstax_restore`
- `searchstax_tags_set`
- `searchstax_user`
- `searchstax_webhook`
- `searchstax_zookeeper_config`
