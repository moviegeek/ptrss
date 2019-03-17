module "scheduler" {
    source = "./scheduler"

    topic_name = "cronjob-topic"
}

module "storage" {
    source = "./storage"

    bucket_name = "ptrss-files"
}