To install camel kafka sqs plugin:
```sh
curl -L -o camel-aws2-sqs-kafka-connector-0.11.5-package.tar.gz https://repo1.maven.org/maven2/org/apache/camel/kafkaconnector/camel-aws2-sqs-kafka-connector/0.11.5/camel-aws2-sqs-kafka-connector-0.11.5-package.tar.gz
```

Extract this plugin and mount it to ```/usr/share/java/plugins``` directory.

To see the list of plugins installed:
```sh
curl http://localhost:8083/connector-plugins | jq
```

After running 02-camelqueue.sh run below command to look at the config used:
```sh
curl http://localhost:8083/connectors/camel-sqs-source/config
```

Check logs:
```sh
docker logs kafka-connect | grep -i plugin
```

# To use Kafka with ngrok
https://rmoff.net/2023/11/01/using-apache-kafka-with-ngrok/

# Running camel-sqs-kafka plugin in linux
Make sure plugin directory looks like below:

camel-aws2-sqs-kafka-connector
    ├── camel-aws2-sqs-kafka-connector-0.11.5.jar
    ├── camel-kafka-connector-0.11.5.jar
    └── lib
        ├── annotations-13.0.jar
        ├── annotations-2.16.88.jar
        ├── apache-client-2.16.88.jar
        ├── apicurio-registry-common-1.3.2.Final.jar
        ├── apicurio-registry-rest-client-1.3.2.Final.jar
        ├── apicurio-registry-utils-converter-1.3.2.Final.jar
        ├── apicurio-registry-utils-serde-1.3.2.Final.jar
        ├── auth-2.16.88.jar
        ├── avro-1.10.2.jar
        ├── aws-core-2.16.88.jar
        ├── aws-query-protocol-2.16.88.jar
        ├── camel-api-3.11.5.jar
        ├── camel-aws2-sqs-3.11.5.jar
        ├── camel-aws2-sqs-kafka-connector-0.11.5.jar
        ├── camel-base-3.11.5.jar
        ├── camel-base-engine-3.11.5.jar
        ├── camel-core-engine-3.11.5.jar
        ├── camel-core-languages-3.11.5.jar
        ├── camel-core-model-3.11.5.jar
        ├── camel-core-processor-3.11.5.jar
        ├── camel-core-reifier-3.11.5.jar
        ├── camel-direct-3.11.5.jar
        ├── camel-jackson-3.11.5.jar
        ├── camel-kafka-3.11.5.jar
        ├── camel-kafka-connector-0.11.5.jar
        ├── camel-main-3.11.5.jar
        ├── camel-management-api-3.11.5.jar
        ├── camel-seda-3.11.5.jar
        ├── camel-support-3.11.5.jar
        ├── camel-util-3.11.5.jar
        ├── commons-codec-1.15.jar
        ├── commons-compress-1.20.jar
        ├── commons-io-2.10.0.jar
        ├── commons-logging-1.2.jar
        ├── eventstream-1.0.1.jar
        ├── http-client-spi-2.16.88.jar
        ├── httpclient-4.5.13.jar
        ├── httpcore-4.4.14.jar
        ├── jackson-dataformat-avro-2.15.2.jar
        ├── javax.annotation-api-1.3.2.jar
        ├── jboss-jaxrs-api_2.1_spec-2.0.1.Final.jar
        ├── jctools-core-3.3.0.jar
        ├── kafka-clients-2.8.0.jar
        ├── kotlin-reflect-1.3.20.jar
        ├── kotlin-stdlib-1.3.20.jar
        ├── kotlin-stdlib-common-1.3.20.jar
        ├── lz4-java-1.7.1.jar
        ├── medeia-validator-core-1.1.1.jar
        ├── medeia-validator-jackson-1.1.1.jar
        ├── metrics-spi-2.16.88.jar
        ├── netty-buffer-4.1.66.Final.jar
        ├── netty-codec-4.1.66.Final.jar
        ├── netty-codec-http-4.1.66.Final.jar
        ├── netty-codec-http2-4.1.66.Final.jar
        ├── netty-common-4.1.66.Final.jar
        ├── netty-handler-4.1.66.Final.jar
        ├── netty-nio-client-2.16.88.jar
        ├── netty-reactive-streams-2.0.5.jar
        ├── netty-reactive-streams-http-2.0.5.jar
        ├── netty-resolver-4.1.66.Final.jar
        ├── netty-transport-4.1.66.Final.jar
        ├── netty-transport-native-epoll-4.1.66.Final-linux-x86_64.jar
        ├── netty-transport-native-unix-common-4.1.66.Final.jar
        ├── okhttp-4.8.1.jar
        ├── okio-2.7.0.jar
        ├── profiles-2.16.88.jar
        ├── protobuf-java-3.13.0.jar
        ├── protocol-core-2.16.88.jar
        ├── reactive-streams-1.0.3.jar
        ├── regions-2.16.88.jar
        ├── retrofit-2.9.0.jar
        ├── sdk-core-2.16.88.jar
        ├── slf4j-api-1.7.30.jar
        ├── snappy-java-1.1.8.1.jar
        ├── sqs-2.16.88.jar
        ├── utils-2.16.88.jar
        └── zstd-jni-1.4.9-1.jar

After extracting the we need to add below extra jar inside lib:
```sh
curl https://repo1.maven.org/maven2/com/fasterxml/jackson/dataformat/jackson-dataformat-avro/2.15.2/jackson-dataformat-avro-2.15.2.jar -o jackson-dataformat-avro-2.15.2.jar
curl https://repo1.maven.org/maven2/org/jctools/jctools-core/3.3.0/jctools-core-3.3.0.jar -o jctools-core-3.3.0.jar
```
docker-compose as below:
```yaml
version: "3.9"

services:
  kafka:
    image: confluentinc/cp-kafka:8.1.1
    container_name: kafka
    hostname: kafka
    ports:
      - "9092:9092"
    environment:
      KAFKA_NODE_ID: 1
      KAFKA_PROCESS_ROLES: "broker,controller"
      KAFKA_CONTROLLER_QUORUM_VOTERS: "1@kafka:9093"
      KAFKA_CONTROLLER_LISTENER_NAMES: "CONTROLLER"
      KAFKA_KRAFT_CLUSTER_ID: "MkU3OEVBNTcwNTJENDM2Qk"
      KAFKA_LISTENERS: "PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093"
      KAFKA_ADVERTISED_LISTENERS: "PLAINTEXT://kafka:9092"
      KAFKA_INTER_BROKER_LISTENER_NAME: "PLAINTEXT"
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: "PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT"
      KAFKA_LOG_DIRS: "/var/lib/kafka/data"
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
      KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0
      KAFKA_CREATE_TOPICS: "redundancy-events:1:1"
    volumes:
      - kafka-data:/var/lib/kafka/data

  kafka-connect:
    image: confluentinc/cp-kafka-connect-base:8.1.1
    volumes:
    - ./plugins:/usr/share/connect/plugins
    container_name: kafka-connect
    depends_on:
      - kafka
    ports:
      - "8083:8083"
    environment:
      CONNECT_BOOTSTRAP_SERVERS: "kafka:9092"
      CONNECT_REST_PORT: 8083
      CONNECT_REST_ADVERTISED_HOST_NAME: kafka-connect
      CONNECT_GROUP_ID: "connect-cluster"
      CONNECT_CONFIG_STORAGE_TOPIC: "_connect-configs"
      CONNECT_OFFSET_STORAGE_TOPIC: "_connect-offsets"
      CONNECT_STATUS_STORAGE_TOPIC: "_connect-status"
      CONNECT_CONFIG_STORAGE_REPLICATION_FACTOR: 1
      CONNECT_OFFSET_STORAGE_REPLICATION_FACTOR: 1
      CONNECT_STATUS_STORAGE_REPLICATION_FACTOR: 1
      CONNECT_KEY_CONVERTER: "org.apache.kafka.connect.storage.StringConverter"
      CONNECT_VALUE_CONVERTER: "org.apache.kafka.connect.storage.StringConverter"
      CONNECT_PLUGIN_PATH: "/usr/share/connect/plugins"
      # CONNECT_LOG4J_ROOT_LOGLEVEL: "DEBUG"
      AWS_REGION: ap-south-1
      AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID}
      AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY}
volumes:
  kafka-data:

```

Add sqs-connector-config.json as below:
```json
{
  "name": "camel-aws2sqs-source",
  "config": {
    "connector.class": "org.apache.camel.kafkaconnector.aws2sqs.CamelAws2sqsSourceConnector",
    "tasks.max": "1",

    "topics": "aws.sqs.events",

    "camel.source.endpoint.queueName": "redundancy-events",
    "camel.source.endpoint.region": "ap-south-1",

    "camel.source.endpoint.deleteAfterRead": "true",
    "camel.source.endpoint.waitTimeSeconds": "20",
    "camel.source.endpoint.maxMessagesPerPoll": "10",

    "key.converter": "org.apache.kafka.connect.storage.StringConverter",
    "value.converter": "org.apache.kafka.connect.storage.StringConverter",

    "errors.tolerance": "all",
    "errors.log.enable": "true",
    "errors.log.include.messages": "true"
  }
}
```

Run below curl command to apply the config:

```sh
curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d @sqs-connector-config.json
```