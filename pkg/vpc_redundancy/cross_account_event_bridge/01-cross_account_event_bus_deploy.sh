#!/usr/bin/env bash
aws cloudformation deploy \
  --region us-east-1 \
  --profile AWS_HOTMAIL \
  --template-file cross_account_event_bus_cf.yaml \
  --stack-name central-target-eventbridge-v1