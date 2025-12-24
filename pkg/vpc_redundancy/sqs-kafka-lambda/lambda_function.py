import os
import json
import logging
from confluent_kafka import Producer

# Configure logging
logger = logging.getLogger()
logger.setLevel(logging.INFO)

# Environment variables
BOOTSTRAP_SERVERS = "2.tcp.us-cal-1.ngrok.io:10428"
KAFKA_TOPIC = "redundancy-events"
SECURITY_PROTOCOL = "PLAINTEXT"

# Kafka config
conf = {
    'bootstrap.servers': BOOTSTRAP_SERVERS,
    'security.protocol': SECURITY_PROTOCOL
}

if SECURITY_PROTOCOL.startswith('SASL'):
    conf.update({
        'sasl.mechanisms': 'PLAIN',
        'sasl.username': SASL_USERNAME,
        'sasl.password': SASL_PASSWORD
    })

producer = Producer(**conf)

def delivery_report(err, msg):
    """Callback after message delivery"""
    if err:
        print(f"Delivery failed: {err}")
    else:
        print(f"Delivered message to {msg.topic()} [{msg.partition()}]")

def lambda_handler(event, context):
    """
    Lambda handler triggered by SQS event
    event['Records'] contains all messages
    """
    logger.info(f"Received EventBridge event: {json.dumps(event)}")
    if 'detail' in event:
        custom_data = event['detail']
        try:
            producer.produce(
                topic=KAFKA_TOPIC,
                value=custom_data.encode('utf-8'),
                callback=delivery_report
            )
        except Exception as e:
            print(f"Error producing message: {e}")

    # Flush producer to ensure delivery
    producer.flush()
    
    return {
        'statusCode': 200,
        'body': json.dumps(f"Processed {len(event['Records'])} messages")
    }
