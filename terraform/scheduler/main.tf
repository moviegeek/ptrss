variable "topic_name" {
    type = "string"
    description = "cloud pub/sub topic name to publish to, when shedule comes"  
}


variable "region" {
    type = "string"
    description = "crontab schedule region"
    # default to per hour
    default = "asia-northeast1"
}
variable "schedule" {
    type = "string"
    description = "crontab schedule definition"
    # default to per hour
    default = "0 * * * *"
}

variable "timezone" {
    type = "string"
    description = "time zone used for the schedule"
    default = "Asia/Tokyo"
}

resource "google_pubsub_topic" "topic" {
  name = "${var.topic_name}"
}

resource "google_cloud_scheduler_job" "job" {
  provider = "google-beta"
  name     = "pt-rss-job"
  region = "${var.region}"
  description = "trigger the pt rss cloud function every 1 hour"
  schedule = "${var.schedule}"
  time_zone = "${var.timezone}"

  pubsub_target {
    topic_name = "${google_pubsub_topic.topic.id}"
    data = "${base64encode("run")}"
  }
}
