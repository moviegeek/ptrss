variable "project" {
  type = "string"
  description = "google project to use"
  default = "movie-221500"
}

variable "credentials" {
    type = "string"
    description = "google service account key"
    default = "/Users/xiaohan/Documents/credentials/pt-rss-deploy_sa-key.json"  
}

variable "region" {
    type = "string"
    description = "google region"
    default = "asia-northeast1"  
}

