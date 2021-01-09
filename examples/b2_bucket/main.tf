terraform {
  required_version = ">= 0.12"
  required_providers {
    b2 = {
      source  = "localhost/backblaze/b2"
      version = "~> 0.1"
    }
  }
}

provider "b2" {
}

resource "b2_bucket" "example" {
  bucket_name = "Example-TestBucket"
  bucket_type = "allPublic"
}

data "b2_bucket" "example" {
  bucket_name = b2_bucket.example.bucket_name
}

output "bucket" {
  value = data.b2_bucket.example
}
