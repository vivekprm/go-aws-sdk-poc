#!/usr/bin/env bash
docker exec -it kafka \
  kafka-topics \
  --bootstrap-server kafka:9092 \
  --create \
  --topic redundancy-events \
  --partitions 1 \
  --replication-factor 1
