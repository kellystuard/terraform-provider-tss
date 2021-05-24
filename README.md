# Terraform Provider for Thycotic Secret Server

:information_source: Since writing this provider, Thycotic has published their own provider to the registry. It can be found at https://registry.terraform.io/providers/thycotic/tss/latest.

:information_source: This provider is in no way authorized or affiliated with Thycotic. I have no association with the company, other than working for a company that uses their Secret Server product. To my knowledge, neither I nor Thycotic provide commercial support for this product. If you have a problem with how the product works, either file an issue or submit a pull request.

This provider has been published to the Terraform Registry at https://registry.terraform.io/providers/kellystuard/tss.

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.14.x
-	[Go](https://golang.org/doc/install) >= 1.16

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

```tf
terraform {
  required_providers {
    tss = {
      source = "kellystuard/tss"
    }
  }
}

provider "tss" {
  username = "john"      # Can be set through environment variable `TSS_USERNAME`.
  password = "<secret>"  # Can be set through environment variable `TSS_PASSWORD`.
  tenant   = "johnvault" # Can be set through environment variable `TSS_TENANT`.
                         # Set to the tenant portion of `https://tenant.secretservercloud.com/`.
                         # If using an on-premise installation, set to the full URI of the server (e.g. -- `https://my-server/SecretServer`).
}

data "tss_secret_field" "test_username" {
  number = 42
  slug   = "username"
}

data "tss_secret_field" "test_password" {
  number = 42
  slug   = "password"
}

output "test_password" {
  value     = "${data.tss_secret_field.test_username.value} : ${data.tss_secret_field.test_password.value}"
  sensitive = true
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
