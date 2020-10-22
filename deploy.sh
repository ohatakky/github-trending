#!/bin/sh

function all() {
  gcloud app deploy cmd/all/app.yaml
}

function rust() {
  gcloud app deploy cmd/rust/app.yaml
}

function cron() {
  gcloud app deploy cron.yaml
}

$1
