# =========================
# CONFIG (edit these)
# =========================
STACK_NAME_PREFIX			:= central-event-bridge
STACK_NAME     				:= event-enrichment-lambda
LAMBDA_NAME    				:= event-enrichment
ARTIFACT_KEY   				:= lambda/lambda.zip
LAMBDA_CREATION_TEMPLATE 	:= lambda-stack.yaml
CENTRAL_PROFILE        		:= AWS_HOTMAIL
CENTRAL_REGION 				:= us-east-1
CENTRAL_ACCOUNT 			:= 665096241598
SOURCE_STACK_PREFIX			:= source-event-bridge
SOURCE_PROFILE 				:= AWS_Personal
SOURCE_REGION 				:= ap-south-1
SOURCE_ACCOUNT 				:= 643716337869
QUEUE_NAME					:= enriched-management-events

# Bucket name will be read from CFN output
BUCKET_OUTPUT  := ArtifactBucketName

# =========================
# INTERNAL
# =========================
BUILD_DIR := build
ZIP_FILE  := function.zip

# =========================
# TARGETS
# =========================

.PHONY: all deploy build upload clean bucket-name

all: deploy

# -------------------------
# Build Go binary for Lambda
# -------------------------
build:
	mkdir -p $(BUILD_DIR)
	cd ../../../../../go-kafka-sqs-lambda && $(MAKE) clean build zip
	cp ../../../../../go-kafka-sqs-lambda/$(ZIP_FILE) $(BUILD_DIR)/

create-bucket:
	aws s3api create-bucket \
		--profile $(CENTRAL_PROFILE) \
		--bucket $(LAMBDA_NAME) \
		--region $(CENTRAL_REGION) \

# -------------------------
# Upload artifact to S3
# -------------------------
upload: create-bucket bucket-name
	aws s3 cp $(BUILD_DIR)/$(ZIP_FILE) \
		s3://$(LAMBDA_NAME)/$(ARTIFACT_KEY) \
		--region $(CENTRAL_REGION) \
		--profile $(CENTRAL_PROFILE)

# -------------------------
# Deploy CloudFormation
# -------------------------
deploy:
# 	step 1
	aws cloudformation deploy \
	--region $(CENTRAL_REGION) \
	--profile $(CENTRAL_PROFILE) \
	--template-file cross_account_event_bus_cf.yaml \
	--stack-name $(STACK_NAME_PREFIX)-v1

#   step 2
	aws events put-permission \
		--region $(CENTRAL_REGION) \
		--profile $(CENTRAL_PROFILE) \
		--event-bus-name central-event-management-bus \
		--statement-id AllowSourceAccount$(SOURCE_ACCOUNT) \
		--principal $(SOURCE_ACCOUNT) \
		--action events:PutEvents

#  	step 3
	aws cloudformation deploy \
		--region $(CENTRAL_REGION) \
		--profile $(CENTRAL_PROFILE) \
		--template-file cross_account_event_bus_permission_cf.yaml \
		--stack-name $(STACK_NAME_PREFIX)-target-permission-v1 \
		--parameter-overrides CentralEventBusArn=arn:aws:events:$(CENTRAL_REGION):$(CENTRAL_ACCOUNT):event-bus/central-event-management-bus \
		SourceAccountId=$(SOURCE_ACCOUNT) \
		--capabilities CAPABILITY_IAM
	
# 	step 4
	aws cloudformation deploy \
		--stack-name $(STACK_NAME_PREFIX)-$(STACK_NAME) \
		--template-file cross_account_event_processing_cf.yaml \
		--region $(CENTRAL_REGION) \
		--profile $(CENTRAL_PROFILE) \
		--parameter-overrides \
		CentralEventBusName=central-event-management-bus \
		S3Bucket=$(LAMBDA_NAME) \
		S3Key=$(ARTIFACT_KEY) \
		QueueName=$(QUEUE_NAME) \
		EventBridgeRuleIAMRoleArn=arn:aws:iam::$(CENTRAL_ACCOUNT):role/service-role/Amazon_EventBridge_Invoke_Sqs_1043361791 \
		--capabilities CAPABILITY_NAMED_IAM

# 	step 5
	aws cloudformation deploy \
		--region $(SOURCE_REGION) \
		--profile $(SOURCE_PROFILE) \
		--template-file forward_management_events_cf.yaml \
		--stack-name $(SOURCE_STACK_PREFIX)-forward-management-events \
		--parameter-overrides \
		TargetAccountId=$(CENTRAL_ACCOUNT) \
		TargetRegion=$(CENTRAL_REGION) \
		CentralEventBusArn=arn:aws:events:$(CENTRAL_REGION):$(CENTRAL_ACCOUNT):event-bus/central-event-management-bus \
		--capabilities CAPABILITY_NAMED_IAM

#  	step 6
	aws cloudformation deploy \
		--region $(SOURCE_REGION) \
		--profile $(SOURCE_PROFILE) \
		--stack-name central-ec2-tag-reader \
		--template-file central-account-assume-role-cf.yaml \
		--parameter-overrides CentralAccountId=$(CENTRAL_ACCOUNT) \
		--capabilities CAPABILITY_NAMED_IAM

consumer:
	go run main.go $(CENTRAL_ACCOUNT) $(QUEUE_NAME) $(CENTRAL_PROFILE) $(CENTRAL_REGION)
# Full pipeline
# -------------------------
release: build upload deploy 

# -------------------------
# Cleanup
# -------------------------
clean:
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME_PREFIX)-$(STACK_NAME) \
		--region $(CENTRAL_REGION) \
		--profile $(CENTRAL_PROFILE)
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME_PREFIX)-target-permission-v1 \
		--region $(CENTRAL_REGION) \
		--profile $(CENTRAL_PROFILE)
	aws cloudformation delete-stack \
		--stack-name $(SOURCE_STACK_PREFIX)-forward-management-events \
		--region $(SOURCE_REGION) \
		--profile $(SOURCE_PROFILE)
	aws cloudformation delete-stack \
		--stack-name central-ec2-tag-reader \
		--region $(SOURCE_REGION) \
		--profile $(SOURCE_PROFILE)
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME_PREFIX)-v1 \
		--region $(CENTRAL_REGION) \
		--profile $(CENTRAL_PROFILE)
	aws s3api delete-object \
		--bucket $(LAMBDA_NAME) \
		--key $(ARTIFACT_KEY) \
		--profile $(CENTRAL_PROFILE) \
		--region $(CENTRAL_REGION)
	aws s3api delete-bucket \
		--bucket $(LAMBDA_NAME) \
		--profile $(CENTRAL_PROFILE) \
		--region $(CENTRAL_REGION)
	rm -rf $(BUILD_DIR)
