LOCATION := "eastus2"
RG := vm-events-rg
STORAGE := vmfuncstorage
FUNC_APP := send-event-func
EVENT_TOPIC := vm-events-topic
TARGETPLATFORM := linux/arm64
AWS_PROFILE := AWS_HOTMAIL
AWS_REGION := us-east-1
SUBSCRIPTION_ID := 22e88069-5977-4fc2-9d6e-ef7e14cfbb70
EVENT_SUBSCRIPTION := vm-events-subscription
export AWS_ACCOUNT_ID := 665096241598
export EVENT_BUS_NAME := central-event-management-bus
TENANT_ID_FILE := tenant_id.txt
${shell touch ${TENANT_ID_FILE}}

setup:
	az group create \
		--name ${RG} \
		--location ${LOCATION}

	az storage account create \
		--name ${STORAGE} \
		--location ${LOCATION} \
		--resource-group ${RG} \
		--sku Standard_LRS
	@AZURE_TENANT_ID=$$(az account show --query tenantId -o tsv); \
	echo "$$AZURE_TENANT_ID" > ${TENANT_ID_FILE}

funcapp: setup
	az functionapp create \
		--resource-group ${RG} \
		--flexconsumption-location ${LOCATION} \
		--runtime python \
		--name ${FUNC_APP} \
		--storage-account ${STORAGE} \
		--os-type Linux \
		--runtime python \
		--runtime-version 3.10 \
		--functions-version 4

publish:
	cd ../../../../../py-azure-sendevent-func && func azure functionapp publish ${FUNC_APP}

eventgrid:
	az eventgrid system-topic create \
		--name ${EVENT_TOPIC} \
		--location global \
		--topic-type "Microsoft.Resources.Subscriptions" \
		--source "/subscriptions/${SUBSCRIPTION_ID}" \
		--resource-group ${RG}

	az eventgrid event-subscription create \
		--name ${EVENT_SUBSCRIPTION} \
		--event-delivery-schema "cloudeventschemav1_0" \
		--source-resource-id "/subscriptions/${SUBSCRIPTION_ID}" \
		--endpoint-type azurefunction \
		--endpoint "/subscriptions/${SUBSCRIPTION_ID}/resourceGroups/${RG}/providers/Microsoft.Web/sites/${FUNC_APP}/functions/EventGridTrigger" \
		--included-event-types \
    		Microsoft.Compute.VirtualMachines.Write \
    		Microsoft.Compute.VirtualMachines.Deallocate
createvm:
	az vm create \
		--resource-group ${RG} \
		--name vm-monitor-01 \
		--image Ubuntu2204 \
		--size Standard_DC1s_v3 \
		--admin-username azureuser \
		--generate-ssh-keys \
		--tags f5xc-site-name=azure-vm-events

thumbprint:
	@AZURE_THUMBPRINT=$$(openssl x509 -in ms-root.pem -fingerprint -sha1 -noout | cut -d= -f2 | tr -d ':' | tr 'A-F' 'a-f') && \
	echo "THUMBPRINT=$$AZURE_THUMBPRINT" > env.mk

aws_oidc:
	@AZURE_TENANT_ID=$$(cat ${TENANT_ID_FILE}); \
	aws iam create-open-id-connect-provider \
		--url https://sts.windows.net/$$AZURE_TENANT_ID/ \
		--client-id-list api://$$AZURE_TENANT_ID/AzureEventBridgePublisher \
		--thumbprint-list "626d44e704d1ceabe3bf0d53397464ac8080142c" \
		--profile ${AWS_PROFILE}

aws_role:
	@AZURE_TENANT_ID=$$(cat ${TENANT_ID_FILE}); \
	export AZURE_TENANT_ID; \
	envsubst < trust-policy.json.tpl > temp-trust-policy.json
	AWS_PAGER="" aws iam create-role \
		--role-name AzureEventBridgePublisher \
		--profile ${AWS_PROFILE} \
		--assume-role-policy-document file://temp-trust-policy.json
	@rm temp-trust-policy.json

aws_policy:
	envsubst < eventbridge-policy.json.tpl > temp-eventbridge-policy.json
	AWS_PAGER="" aws iam create-policy \
		--policy-name AzurePutEventsPolicy \
		--profile ${AWS_PROFILE} \
		--policy-document file://temp-eventbridge-policy.json
	@rm temp-eventbridge-policy.json

	aws iam attach-role-policy \
		--role-name AzureEventBridgePublisher \
		--policy-arn arn:aws:iam::${AWS_ACCOUNT_ID}:policy/AzurePutEventsPolicy \
		--profile ${AWS_PROFILE}

aws_iam: aws_oidc aws_role aws_policy

azure_oidc:
	@AZURE_TENANT_ID=$$(cat ${TENANT_ID_FILE}); \
	az ad app create \
		--display-name AzureEventBridgePublisher \
		--identifier-uris api://$$AZURE_TENANT_ID/AzureEventBridgePublisher \
		--sign-in-audience AzureADMyOrg
	az functionapp identity assign \
		--resource-group ${RG} \
		--name ${FUNC_APP}; \

	APP_ID=$$(az ad app list \
	--display-name AzureEventBridgePublisher \
	--query '[0].appId' -o tsv); \
	APP_OBJECT_ID=$$(az ad app show \
	--id $$APP_ID \
	--query id -o tsv); \
	MANAGED_IDENTITY_OBJECT_ID=$$(az functionapp identity show \
	--resource-group ${RG} \
	--name ${FUNC_APP} \
	--query principalId -o tsv); \
	echo "APP_ID=$$APP_ID"; \
	echo "APP_OBJECT_ID=$$APP_OBJECT_ID"; \
	echo "MANAGED_IDENTITY_OBJECT_ID=$$MANAGED_IDENTITY_OBJECT_ID"; \
	az ad app federated-credential create \
		--id $$APP_OBJECT_ID \
		--parameters '{ \
			"name": "AzureFunctionCredential", \
			"issuer": "https://sts.windows.net/$$AZURE_TENANT_ID", \
			"subject": "system:principalid:$$MANAGED_IDENTITY_OBJECT_ID", \
			"audiences": ["api://$$AZURE_TENANT_ID/AzureEventBridgePublisher"] \
		}'
		
	az functionapp config appsettings set \
		--resource-group ${RG} \
		--name ${FUNC_APP} \
		--settings \
			AWS_ROLE_ARN=arn:aws:iam::${AWS_ACCOUNT_ID}:role/AzureEventBridgePublisher \
			AWS_WEB_IDENTITY_TOKEN_FILE=/var/run/secrets/azure/tokens/azure-identity-token \
			AWS_REGION=${AWS_REGION} \
			EVENT_BUS_NAME=${EVENT_BUS_NAME} \
			AWS_WEB_IDENTITY_TOKEN_FILE="/var/run/secrets/azure/tokens/azure-identity-token"

azure_setup: setup funcapp publish azure_oidc createvm eventgrid

release: azure_setup aws_iam

azure_cleanup:
	az functionapp delete --name ${FUNC_APP} --resource-group ${RG}
	az storage account delete --name ${STORAGE} --resource-group ${RG} --yes
	APP_ID=$$(az ad app list --display-name AzureEventBridgePublisher --query '[0].appId' -o tsv); \
	az ad app delete --id $$APP_ID
	az eventgrid system-topic delete --name ${EVENT_TOPIC} --resource-group ${RG} --yes
	az vm delete --resource-group ${RG} --name vm-monitor-01 --yes
	az eventgrid event-subscription delete --name ${EVENT_SUBSCRIPTION} --source-resource-id "/subscriptions/${SUBSCRIPTION_ID}" --yes
	az group delete --name ${RG} --yes --no-wait

aws_cleanup:
	@AZURE_TENANT_ID=$$(cat ${TENANT_ID_FILE}); \
	aws iam delete-open-id-connect-provider \
		--open-id-connect-provider-arn arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/sts.windows.net/$$AZURE_TENANT_ID/ \
		--profile ${AWS_PROFILE}
	aws iam detach-role-policy \
		--role-name AzureEventBridgePublisher \
		--policy-arn "arn:aws:iam::${AWS_ACCOUNT_ID}:policy/AzurePutEventsPolicy" \
		--profile ${AWS_PROFILE} 
	aws iam delete-role \
		--role-name AzureEventBridgePublisher \
		--profile ${AWS_PROFILE}
	aws iam delete-policy \
		--policy-arn arn:aws:iam::${AWS_ACCOUNT_ID}:policy/AzurePutEventsPolicy \
		--profile ${AWS_PROFILE}

clean: azure_cleanup aws_cleanup
	rm -f ${TENANT_ID_FILE}