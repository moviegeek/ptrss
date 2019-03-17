variable "bucket_name" {
    type = "string"
    description = "storage bucket name"  
}

variable "storage_class" {
    type = "string"
    description = "storage class for the bucket"    
    default = "REGIONAL"
}

variable "location" {
    type = "string"
    description = "storage location"
    default = "asia-northeast1"  
}

resource "google_storage_bucket" "main" {
  name     = "${var.bucket_name}"
  storage_class = "${var.storage_class}"
  location = "${var.location}"

  versioning {
      enabled = true
  }
}

resource "google_storage_bucket_object" "feeds-xml" {
  name   = "pt-rss.xml"
  content = "dummy"
  bucket = "${google_storage_bucket.main.name}"

  content_type = "application/rss+xml;charset=utf-8"
}

resource "google_storage_object_acl" "rss-acl" {
  bucket = "${google_storage_bucket.main.name}"
  object = "${google_storage_bucket_object.feeds-xml.output_name}"

  predefined_acl = "publicRead"
}
resource "google_storage_bucket_object" "feeds-json" {
  name   = "pt-rss.json"
  content = "{}"
  bucket = "${google_storage_bucket.main.name}"

  content_type = "application/json"
}