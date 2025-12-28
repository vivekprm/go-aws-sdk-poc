aws cloudformation deploy \
  --region ap-south-1 \
  --profile AWS_Personal \
  --stack-name central-ec2-tag-reader \
  --template-file central-account-assume-role-cf.yaml \
  --parameter-overrides CentralAccountId=665096241598 \
  --capabilities CAPABILITY_NAMED_IAM
