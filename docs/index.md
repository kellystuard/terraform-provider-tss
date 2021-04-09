---
page_title: "Thycotic Secret Server Terraform Provider"
subcategory: ""
description: |-
  
---

# Thycotic Secret Server Terraform Provider



## Example Usage

```terraform
provider tss {
  username   = "my-username" // or set env TSS_USERNAME
  password   = "my-password" // or set env TSS_PASSWORD
  tenant     = "my-tenant"   // or set env TSS_TENANT
  timeout    = "10s"         // or set env TSS_TIMEOUT
  grant_type = "password"    // or set env TSS_GRANT_TYPE
}
```

## Schema
