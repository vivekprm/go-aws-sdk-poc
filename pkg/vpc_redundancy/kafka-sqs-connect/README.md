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