# Terraform Provider SearchStax (WIP)
<div align="center">
    <img src="https://www.searchstax.com/docs/wp-content/themes/docs/images/logo.svg" width="400" alt="SearchStax" />
    <br/>
   <a href="https://www.searchstax.com/docs/searchstax-cloud-apis-overview/">
    <img src="https://img.shields.io/static/v1?label=Docs&message=API Ref&color=000000&style=for-the-badge"  alt="SearchStax Cloud API Documentation"/>
    </a>

[![Tests](https://github.com/sadok-f/terraform-provider-searchstax/actions/workflows/test.yml/badge.svg)](https://github.com/sadok-f/terraform-provider-searchstax/actions/workflows/test.yml)
</div>

This repository represents an base code for a terraform provider for [SearchStax Cloud](https://www.searchstax.com/docs/searchstax-cloud-docs-home/).
The code is still **WIP** and not published yet to the Terraform registry.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

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


## Using the provider

```hcl
terraform {
  required_providers {
    searchstax = {
      source = "hashicorp.com/sadok-f/searchstax"
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

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests runs on a SearchStax mock api server created by this [Repo](https://github.com/sadok-f/searchstax-mock-api), it returns a mocking response for each API call.

```shell
make testacc
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

[https://opencredo.com/blogs/running-a-terraform-provider-with-a-debugger/](https://opencredo.com/blogs/running-a-terraform-provider-with-a-debugger/)
