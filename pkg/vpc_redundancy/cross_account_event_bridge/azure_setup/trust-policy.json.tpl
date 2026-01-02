{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/sts.windows.net/${AZURE_TENANT_ID}/"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
        "sts.windows.net/${AZURE_TENANT_ID}/:aud":
        "api://${AZURE_TENANT_ID}/AzureEventBridgePublisher"
        }
      }
    }
  ]
}