#!/bin/sh

function all() {
  gcloud app deploy --appyaml=app.yaml
}

function rust() {
  gcloud app deploy --appyaml=app-rust.yaml
}

function cron() {
  gcloud app deploy cron.yaml
}
