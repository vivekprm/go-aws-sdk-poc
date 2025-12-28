#!/usr/bin/env bash
aws s3api create-bucket \
  --profile AWS_HOTMAIL \
  --bucket enrichment-lambda \
  --region us-east-1

aws s3 cp central-ec2-event-enrichment.zip s3://enrichment-lambda/central-ec2-event-enrichment.zip \
  --profile AWS_HOTMAIL \
  --region us-east-1
