curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "camel-sqs-source",
    "config": {
      "connector.class": "org.apache.camel.kafkaconnector.aws2sqs.CamelAws2sqsSourceConnector",
      "tasks.max": "1",
      "camel.source.path.queueNameOrArn": "redundancy-events-queue",
      "camel.source.endpoint.region": "ap-south-1",
      "camel.component.aws2-sqs.accessKey": "'"$AWS_ACCESS_KEY_ID"'",
      "camel.component.aws2-sqs.secretKey": "'"$AWS_SECRET_ACCESS_KEY"'",
      "topics": "redundancy-events"
    }
  }'
