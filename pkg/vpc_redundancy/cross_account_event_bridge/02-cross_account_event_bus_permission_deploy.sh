#!/usr/bin/env bash
aws events put-permission \
  --region us-east-1 \
  --profile AWS_HOTMAIL \
  --event-bus-name central-event-management-bus \
  --statement-id AllowSourceAccount643716337869 \
  --principal 643716337869 \
  --action events:PutEvents

aws cloudformation deploy \
  --region us-east-1 \
  --profile AWS_HOTMAIL \
  --template-file cross_account_event_bus_permission_cf.yaml \
  --stack-name central-eventbridge-target-permission-v1 \
  --parameter-overrides CentralEventBusArn=arn:aws:events:us-east-1:665096241598:event-bus/central-event-management-bus \
  SourceAccountId=643716337869 \
  --capabilities CAPABILITY_IAM
