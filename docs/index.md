---
page_title: "B2 Provider"
subcategory: ""
description: |-
  
---
# B2 Provider

Terraform provider for [Backblaze B2](https://www.backblaze.com/b2/).

The provider is written in go, but it uses official [B2 python SDK](https://github.com/Backblaze/b2-sdk-python) embedded into the binary.

## Example Usage
```terraform
terraform {
  required_version = ">= 0.13"
  required_providers {
    b2 = {
      source  = "Backblaze/b2"
      version = "~> 0.2"
    }
  }
}

provider "b2" {
}

resource "b2_application_key" "example_key" {
  key_name     = "my-key"
  capabilities = ["readFiles"]
}

resource "b2_bucket" "example_bucket" {
  bucket_name = "my-b2-bucket"
  bucket_type = "allPublic"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **application_key** (String, Sensitive) B2 Application Key (B2_APPLICATION_KEY env)
- **application_key_id** (String, Sensitive) B2 Application Key ID (B2_APPLICATION_KEY_ID env)
- **endpoint** (String) B2 endpoint - production or custom URL (B2_ENDPOINT env)