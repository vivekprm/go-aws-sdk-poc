#!/usr/bin/env bash

aws events put-rule \
  --name azure-gcp-to-sqs \
  --event-bus-name central-event-management-bus \
  --event-pattern file://event-pattern.json \
  --region us-east-1 \
  --profile AWS_HOTMAIL

aws events put-targets \
  --event-bus-name central-event-management-bus \
  --rule azure-gcp-to-sqs \
  --targets "Id"="SqsTarget","Arn"="arn:aws:sqs:us-east-1:665096241598:enriched-management-events" \
  --region us-east-1 \
  --profile AWS_HOTMAIL