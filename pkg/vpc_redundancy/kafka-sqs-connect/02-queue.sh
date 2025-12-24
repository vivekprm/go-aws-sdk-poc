curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "sqs-source-connector",
    "config": {
      "connector.class": "io.confluent.connect.sqs.source.SqsSourceConnector",
      "tasks.max": "1",
      "sqs.url": "https://sqs.ap-south-1.amazonaws.com/643716337869/redundancy-events-queue",
      "confluent.topic.bootstrap.servers": "kafka:9092",
      "aws.region": "ap-south-1",
      "aws.access.key.id": "'"$AWS_ACCESS_KEY_ID"'",
      "aws.secret.key.id": "'"$AWS_SECRET_ACCESS_KEY"'",
      "kafka.topic": "redundancy-events",
      "polling.wait.time.ms": "20000",
      "visibility.timeout.ms": "30000",
      "max.number.of.messages": "10",
      "behavior.on.error": "log",
      "key.converter": "org.apache.kafka.connect.storage.StringConverter",
      "value.converter": "org.apache.kafka.connect.json.JsonConverter",
      "value.converter.schemas.enable": "false"
    }
  }'
