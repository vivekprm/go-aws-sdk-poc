# Helper commands
# Create connector
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

# Get connectors list, config and status
curl http://localhost:8083/connectors
curl http://localhost:8083/connectors/camel-sqs-source/config | jq
curl http://localhost:8083/connectors/camel-sqs-source/status | jq

# To connect to kafka cluster
kcat -L -b kafka:9092
kcat -b kafka:9092 -t redundancy-events
echo "hello world" | kcat -b kafka:9092 -P -t redundancy-events 