#!/bin/sh

function trending_function() {
  gcloud functions deploy TrendingHTTP --runtime go113 \
  --trigger-http \
  --entry-point=Handler \
  --region=asia-northeast1 \
  --env-vars-file .env.yaml
  # --ingress-settings=internal-only \
}

function trending_scheduler() {
  gcloud scheduler jobs create TrendingScheduler http \
  --schedule="0 10 * * *" \
  --time-zone="Asia/Tokyo" \
  --uri=${TRENDING_FUNCTION_URI} \
  --oidc-service-account-email=${TRENDING_FUNCTION_SERVICE_ACCOUNT}
}

$1
