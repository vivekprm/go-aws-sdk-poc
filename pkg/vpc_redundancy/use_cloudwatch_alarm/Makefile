AWS_PROFILE := AWS_Personal
AWS_ACCOUNT_ID := 643716337869
AWS_REGION := ap-south-1
QUEUE_NAME := redundancy-events
NAT_GATEWAY_ID := nat-0134cca92fee8caca

eventbridge:
	aws events put-rule \
		--name ec2-state-change-rule \
		--profile ${AWS_PROFILE} \
		--event-pattern '{ \
			"source": ["aws.ec2"], \
			"detail-type": ["EC2 Instance State-change Notification"], \
			"detail": { \
			"state": ["stopped", "terminated"] \
			} \
		}'
	aws events put-rule \
		--name cloudwatch-alarm-to-eventbridge \
		--profile ${AWS_PROFILE} \
		--event-pattern '{ \
			"source": ["aws.cloudwatch"], \
			"detail-type": ["CloudWatch Alarm State Change"], \
			"detail": { \
			"alarmName": ["NAT-Gateway-Failure"], \
			"state": { \
				"value": ["ALARM"] \
			} \
			} \
		}'

cloudwatch:
	aws cloudwatch put-metric-alarm \
		--profile ${AWS_PROFILE} \
		--alarm-name NAT-Gateway-Failure \
		--alarm-description "NAT Gateway functional failure detected" \
		--namespace AWS/NATGateway \
		--metric-name BytesOutToDestination \
		--dimensions Name=NatGatewayId,Value=${NAT_GATEWAY_ID} \
		--statistic Sum \
		--period 60 \
		--evaluation-periods 3 \
		--threshold 0 \
		--comparison-operator LessThanOrEqualToThreshold \
		--treat-missing-data breaching

iam:
	@QUEUE_URL=$$(aws sqs get-queue-url \
		--queue-name ${QUEUE_NAME} \
		--profile ${AWS_PROFILE} \
		--query 'QueueUrl' \
		--output text); \
	echo "Setting policy for SQS Queue: $$QUEUE_URL"; \
	aws sqs set-queue-attributes \
		--queue-url $$QUEUE_URL \
		--profile ${AWS_PROFILE} \
		--attributes '{"Policy": "{\"Version\":\"2012-10-17\", \"Statement\":[{ \"Effect\":\"Allow\", \"Principal\":{\"Service\":\"events.amazonaws.com\"}, \"Action\":\"sqs:SendMessage\", \"Resource\":\"arn:aws:sqs:${AWS_REGION}:${AWS_ACCOUNT_ID}:${QUEUE_NAME}\"}]}" \
		}'

target:
	aws sqs create-queue \
		--queue-name ${QUEUE_NAME} \
		--profile ${AWS_PROFILE}
	aws events put-targets \
		--rule cloudwatch-alarm-to-eventbridge \
		--profile ${AWS_PROFILE} \
		--targets '[ \
			{ \
			"Id": "cw-alarm-sqs", \
			"Arn": "arn:aws:sqs:${AWS_REGION}:${AWS_ACCOUNT_ID}:${QUEUE_NAME}" \
			} \
		]'
	aws events put-targets \
		--rule ec2-state-change-rule \
		--profile ${AWS_PROFILE} \
		--targets '[ \
			{ \
			"Id": "ec2-state-change-sqs", \
			"Arn": "arn:aws:sqs:${AWS_REGION}:${AWS_ACCOUNT_ID}:${QUEUE_NAME}" \
			} \
		]'

consumer:
	@QUEUE_URL=$$(aws sqs get-queue-url \
		--queue-name ${QUEUE_NAME} \
		--profile ${AWS_PROFILE} \
		--query 'QueueUrl' \
		--output text); \
	echo "Consuming messages from SQS Queue: $$QUEUE_URL"; \
	go run ../cross_account_event_bridge/aws_setup/main.go $(AWS_ACCOUNT_ID) $(QUEUE_NAME) $(AWS_PROFILE) $(AWS_REGION)

release: cloudwatch eventbridge target iam consumer

clean:
	aws cloudwatch delete-alarms \
		--alarm-names NAT-Gateway-Failure \
		--profile ${AWS_PROFILE}
	aws events remove-targets \
		--rule cloudwatch-alarm-to-eventbridge \
		--ids cw-alarm-sqs \
		--profile ${AWS_PROFILE}
	aws events remove-targets \
		--rule ec2-state-change-rule \
		--ids ec2-state-change-sqs \
		--profile ${AWS_PROFILE}
	aws events delete-rule \
		--name cloudwatch-alarm-to-eventbridge \
		--profile ${AWS_PROFILE}
	aws events delete-rule \
		--name ec2-state-change-rule \
		--profile ${AWS_PROFILE}
	aws sqs delete-queue \
		--profile ${AWS_PROFILE} \
		--queue-url $$(aws sqs get-queue-url \
			--queue-name ${QUEUE_NAME} \
			--profile ${AWS_PROFILE} \
			--query 'QueueUrl' \
			--output text)