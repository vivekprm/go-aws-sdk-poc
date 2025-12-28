#!/usr/bin/env bash

aws cloudformation deploy \
  --region us-east-1 \
  --profile AWS_HOTMAIL \
  --stack-name event-processing-stack-v1 \
  --template-file cross_account_event_processing_cf.yaml \
  --capabilities CAPABILITY_NAMED_IAM
