AWS_PROFILE := AWS_HOTMAIL
AWS_REGION := us-east-1
export AWS_ACCOUNT_ID := 665096241598
export EVENT_BUS_NAME := central-event-management-bus
PROJECT_ID := project-fa5880b0-4880-4755-92e
PROJECT_NUMBER := 657741258816
PUBSUB_TOPIC := gcp-vm-events
GCP_REGION := us-central1
GCP_ZONE := us-central1-a
CLOUDRUN_APP := gcp-vm-event-bridge
CLOUDRUN_URL_FILE := cloudrun_url.txt
${shell touch ${CLOUDRUN_URL_FILE}}
SERVICE_ACCOUNT_ID_FILE := service_account_id.txt
${shell touch ${SERVICE_ACCOUNT_ID_FILE}}

gcp_setup:
	gcloud config set project ${PROJECT_ID}
	gcloud services enable \
  		compute.googleapis.com \
  		logging.googleapis.com \
  		pubsub.googleapis.com \
		run.googleapis.com \
		artifactregistry.googleapis.com

gcp_pubsub: cloudrun_hash
	gcloud pubsub topics create ${PUBSUB_TOPIC}
	gcloud logging sinks create gcp-vm-sink \
		pubsub.googleapis.com/projects/${PROJECT_ID}/topics/${PUBSUB_TOPIC} \
		--log-filter='resource.type="gce_instance" AND (protoPayload.methodName="v1.compute.instances.insert" OR protoPayload.methodName="v1.compute.instances.delete" OR protoPayload.methodName="v1.compute.instances.stop" OR protoPayload.methodName="v1.compute.instances.start")'
	@CLOUDRUN_URL=$$(cat $(CLOUDRUN_URL_FILE)); \
	gcloud pubsub subscriptions create gcp-vm-events-sub \
		--topic ${PUBSUB_TOPIC} \
		--push-endpoint=$$CLOUDRUN_URL \
		--push-auth-service-account=eventbridge-publisher@${PROJECT_ID}.iam.gserviceaccount.com
	
	gcloud pubsub topics add-iam-policy-binding ${PUBSUB_TOPIC} \
  		--member="serviceAccount:service-${PROJECT_NUMBER}@gcp-sa-logging.iam.gserviceaccount.com" \
  		--role="roles/pubsub.publisher"

createvm:
	gcloud compute instances create test-vm-1 \
	--zone=us-central1-a \
	--machine-type=e2-micro \
	--image-family=debian-12 \
	--image-project=debian-cloud \
	--boot-disk-size=10GB \
	--tags=vm-events-test

gcp_iam:
	gcloud iam service-accounts create eventbridge-publisher
	@sleep 5
	@gcloud iam service-accounts describe eventbridge-publisher@${PROJECT_ID}.iam.gserviceaccount.com --format="value(uniqueId)"> ${SERVICE_ACCOUNT_ID_FILE}
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
		--member="serviceAccount:eventbridge-publisher@${PROJECT_ID}.iam.gserviceaccount.com" \
		--role="roles/pubsub.subscriber"
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  		--member="serviceAccount:eventbridge-publisher@${PROJECT_ID}.iam.gserviceaccount.com" \
  		--role="roles/iam.serviceAccountTokenCreator"

create-registry:
	gcloud artifacts repositories create gcr.io \
		--repository-format=docker \
		--location=us \
		--description="Container Registry compatibility repo"

cloudrun_deploy: create-registry
# 	Needs cloud build permission
# 	gcloud projects add-iam-policy-binding <PROJECT_ID> \
  		--member="user:user@gmail.com" \
  		--role="roles/cloudbuild.builds.editor"
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  		--member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
  		--role="roles/artifactregistry.writer"
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  		--member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
  		--role="roles/artifactregistry.writer"
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
		--member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
		--role="roles/logging.logWriter"
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  		--member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
  		--role="roles/storage.objectViewer"
	@CLOUDRUN_URL=$$(cat $(CLOUDRUN_URL_FILE)); \
	cd ../../../../../go-gcp-sendevent-cloudrun && \
		gcloud builds submit --tag gcr.io/${PROJECT_ID}/${CLOUDRUN_APP}:latest . && \
	gcloud run deploy ${CLOUDRUN_APP} \
		--image gcr.io/${PROJECT_ID}/${CLOUDRUN_APP}:latest \
		--platform managed \
		--region ${GCP_REGION} \
		--allow-unauthenticated \
		--service-account=eventbridge-publisher@${PROJECT_ID}.iam.gserviceaccount.com \
		--set-env-vars AWS_ROLE_ARN="arn:aws:iam::${AWS_ACCOUNT_ID}:role/GcpEventBridgePublisher",AWS_REGION="${AWS_REGION}",EVENT_BUS_NAME="${EVENT_BUS_NAME}",AWS_EC2_METADATA_DISABLED=true,GOOGLE_AUDIENCE=$$CLOUDRUN_URL

cloudrun_hash:
	@CLOUDRUN_URL=$$(gcloud run services describe ${CLOUDRUN_APP} \
  		--region ${GCP_REGION} \
  		--format="value(status.url)"); \
	echo "$$CLOUDRUN_URL" > ${CLOUDRUN_URL_FILE}

aws_iam:
	SERVICE_ACCOUNT_ID=$$(cat ${SERVICE_ACCOUNT_ID_FILE});\
	export SERVICE_ACCOUNT_ID; \
	aws iam create-open-id-connect-provider \
  		--url https://accounts.google.com \
  		--client-id-list $$SERVICE_ACCOUNT_ID \
  		--thumbprint-list cf23df2207d99a74fbe169e3eba035e633b65d94 \
		--profile ${AWS_PROFILE}; \
	envsubst < trust-policy.json.tpl > temp-trust-policy.json
	AWS_PAGER="" aws iam create-role \
		--role-name GcpEventBridgePublisher \
		--profile ${AWS_PROFILE} \
		--assume-role-policy-document file://temp-trust-policy.json
	@rm temp-trust-policy.json

	envsubst < eventbridge-policy.json.tpl > temp-eventbridge-policy.json
	AWS_PAGER="" aws iam create-policy \
		--policy-name GcpPutEventsPolicy \
		--profile ${AWS_PROFILE} \
		--policy-document file://temp-eventbridge-policy.json
	@rm temp-eventbridge-policy.json
	aws iam attach-role-policy \
		--role-name GcpEventBridgePublisher \
		--policy-arn arn:aws:iam::${AWS_ACCOUNT_ID}:policy/GcpPutEventsPolicy \
		--profile ${AWS_PROFILE}

release: gcp_setup gcp_iam cloudrun_deploy gcp_pubsub createvm aws_iam

gcp_cleanup:
	gcloud projects remove-iam-policy-binding ${PROJECT_ID} \
		--member="serviceAccount:eventbridge-publisher@${PROJECT_ID}.iam.gserviceaccount.com" \
		--role="roles/pubsub.subscriber"
	gcloud iam service-accounts delete eventbridge-publisher@${PROJECT_ID}.iam.gserviceaccount.com -q
	gcloud projects remove-iam-policy-binding ${PROJECT_ID} \
  		--member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
  		--role="roles/storage.objectViewer"
	gcloud projects remove-iam-policy-binding ${PROJECT_ID} \
  		--member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
  		--role="roles/logging.logWriter"
	gcloud storage rm gs://${PROJECT_ID}_cloudbuild/source/ --recursive
	gcloud run services delete ${CLOUDRUN_APP} --region=${GCP_REGION} -q
	gcloud artifacts repositories delete gcr.io --location=us -q
	gcloud pubsub topics delete ${PUBSUB_TOPIC}
	gcloud logging sinks delete gcp-vm-sink -q
	gcloud pubsub subscriptions delete gcp-vm-events-sub
	gcloud compute instances delete test-vm-1 --zone=${GCP_ZONE} -q

aws_cleanup:
	aws iam delete-open-id-connect-provider \
  		--open-id-connect-provider-arn arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/accounts.google.com \
		--profile ${AWS_PROFILE}
	aws iam detach-role-policy \
		--role-name GcpEventBridgePublisher \
		--policy-arn arn:aws:iam::${AWS_ACCOUNT_ID}:policy/GcpPutEventsPolicy \
		--profile ${AWS_PROFILE}
	aws iam delete-policy \
		--policy-arn arn:aws:iam::${AWS_ACCOUNT_ID}:policy/GcpPutEventsPolicy \
		--profile ${AWS_PROFILE}
	aws iam delete-role \
		--role-name GcpEventBridgePublisher \
		--profile ${AWS_PROFILE}

clean: gcp_cleanup aws_cleanup
	rm -f ${CLOUDRUN_URL_FILE} ${SERVICE_ACCOUNT_ID_FILE}
	