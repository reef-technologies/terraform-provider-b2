---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "b2_bucket Resource - terraform-provider-b2"
subcategory: ""
description: |-
  B2 bucket resource.
---

# Resource `b2_bucket`

B2 bucket resource.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **bucket_name** (String) The name of the bucket.
- **bucket_type** (String) The bucket type.

### Optional

- **bucket_info** (Map of String) The bucket info.
- **cors_rules** (Block List) CORS rules. (see [below for nested schema](#nestedblock--cors_rules))
- **id** (String) The ID of this resource.
- **lifecycle_rules** (Block List) Lifecycle rules. (see [below for nested schema](#nestedblock--lifecycle_rules))

### Read-only

- **account_id** (String) Account ID that the bucket belongs to.
- **bucket_id** (String) The ID of the bucket.
- **options** (Set of String) List of bucket options.
- **revision** (Number) Bucket revision.

<a id="nestedblock--cors_rules"></a>
### Nested Schema for `cors_rules`

Required:

- **allowed_operations** (List of String)
- **allowed_origins** (List of String)
- **cors_rule_name** (String)
- **max_age_seconds** (Number)

Optional:

- **allowed_headers** (List of String)
- **expose_headers** (List of String)


<a id="nestedblock--lifecycle_rules"></a>
### Nested Schema for `lifecycle_rules`

Required:

- **file_name_prefix** (String)

Optional:

- **days_from_hiding_to_deleting** (Number)
- **days_from_uploading_to_hiding** (Number)

