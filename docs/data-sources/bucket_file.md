---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "b2_bucket_file Data Source - terraform-provider-b2"
subcategory: ""
description: |-
  B2 bucket file data source.
---

# Data Source `b2_bucket_file`

B2 bucket file data source.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **bucket_id** (String) The ID of the bucket.
- **file_name** (String) The file name.

### Optional

- **id** (String) The ID of this resource.
- **show_versions** (Boolean) Show all file versions.

### Read-only

- **file_versions** (List of Object) File versions. (see [below for nested schema](#nestedatt--file_versions))

<a id="nestedatt--file_versions"></a>
### Nested Schema for `file_versions`

Read-only:

- **action** (String)
- **content_md5** (String)
- **content_sha1** (String)
- **content_type** (String)
- **file_id** (String)
- **file_info** (Map of String)
- **file_name** (String)
- **size** (Number)
- **upload_timestamp** (Number)

