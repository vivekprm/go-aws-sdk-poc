#!/usr/bin/env bash

docker exec -it kafka kafka-console-consumer \
  --bootstrap-server 2.tcp.us-cal-1.ngrok.io:11630 \
  --topic redundancy-events \
  --from-beginning \
  --property print.key=true \
  --property print.value=true
