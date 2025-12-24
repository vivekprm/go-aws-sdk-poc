#!/usr/bin/env bash

docker exec -it kafka kafka-console-consumer \
  --bootstrap-server kafka:9092 \
  --topic redundancy-events \
  --from-beginning \
  --property print.key=true \
  --property print.value=true
