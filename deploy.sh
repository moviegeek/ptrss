#!/bin/bash

gcloud functions deploy \
  ptrss \
  --entry-point UpdateRss \
  --memory 128MB \
  --retry \
  --runtime go111 \
  --service-account pt-rss-app@movie-221500.iam.gserviceaccount.com \
  --source . \
  --timeout 60s \
  --env-vars-file .env.yaml \
  --trigger-topic cronjob-topic