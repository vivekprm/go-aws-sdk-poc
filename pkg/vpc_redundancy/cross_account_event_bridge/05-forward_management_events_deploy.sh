aws cloudformation deploy \
  --region ap-south-1 \
  --profile AWS_Personal \
  --template-file forward_management_events_cf.yaml \
  --stack-name forward-management-events \
  --parameter-overrides \
      TargetAccountId=665096241598 \
      TargetRegion=us-east-1 \
      CentralEventBusArn=arn:aws:events:us-east-1:665096241598:event-bus/central-event-management-bus \
  --capabilities CAPABILITY_NAMED_IAM
