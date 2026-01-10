AWS_ACCOUNT_ID := 643716337869
AWS_PROFILE := AWS_Personal
AWS_REGION := ap-south-1
PREFIX_NAME := vpc-redundancy
CENTRAL_ACCOUNT_ID := 665096241598
QUEUE_NAME := central-vpc-redundancy-queue

eventbridge:
	aws events put-rule \
		--name vpc-redundancy-rule \
		--profile ${AWS_PROFILE} \
		--event-pattern '{ \
			"source": ["custom.vpc-redundancy"], \
			"detail-type": ["VPC Redundancy Event"] \
		}'