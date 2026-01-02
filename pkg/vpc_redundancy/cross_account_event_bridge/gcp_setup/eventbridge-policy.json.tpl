{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "events:PutEvents",
      "Resource": "arn:aws:events:us-east-1:${AWS_ACCOUNT_ID}:event-bus/${EVENT_BUS_NAME}"
    }
  ]
}
