LOCATION := "eastus2"
RG := vm-events-rg
STORAGE := vmfuncstorage
FUNC_APP := send-event-func
EVENT_TOPIC := vm-events-topic
TARGETPLATFORM := linux/arm64
AWS_PROFILE := AWS_HOTMAIL
SUBSCRIPTION_ID := 22e88069-5977-4fc2-9d6e-ef7e14cfbb70
export AWS_ACCOUNT_ID := 665096241598
export EVENT_BUS_NAME := central-event-management-bus

setup:
	az group create \
		--name ${RG} \
		--location ${LOCATION}

	az storage account create \
		--name ${STORAGE} \
		--location ${LOCATION} \
		--resource-group ${RG} \
		--sku Standard_LRS

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

release: funcapp publish eventgrid createvm

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
	@openssl x509 -in ms-root.pem -fingerprint -sha1 -noout | cut -d= -f2 | tr -d ':' | tr 'A-F' 'a-f'

azure_tenant_name:
	@az account show --query tenantId -o tsv

oidcsetup:
	@AZURE_TENANT_ID=$$(make -s azure_tenant_name); \
	@AZURE_THUMBPRINT=$$(make -s thumbprint); \
	aws iam create-open-id-connect-provider \
		--url https://sts.windows.net/$$AZURE_TENANT_ID/ \
		--client-id-list api://$$AZURE_TENANT_ID/AzureEventBridgePublisher \
		--thumbprint-list $$AZURE_THUMBPRINT \
		--profile ${AWS_PROFILE}

aws_role:
	@export AZURE_TENANT_ID=$$(make -s azure_tenant_name); \
	envsubst < trust-policy.json.tpl > temp-trust-policy.json
	aws iam create-role \
		--role-name AzureEventBridgePublisher \
		--profile ${AWS_PROFILE} \
		--assume-role-policy-document file://temp-trust-policy.json
	@rm temp-trust-policy.json

aws_policy:
	envsubst < eventbridge-policy.json.tpl > temp-eventbridge-policy.json
	aws iam create-policy \
		--policy-name AzurePutEventsPolicy \
		--profile ${AWS_PROFILE} \
		--policy-document file://temp-eventbridge-policy.json
	@rm temp-eventbridge-policy.json

	aws iam attach-role-policy \
		--role-name AzureEventBridgePublisher \
		--policy-arn arn:aws:iam::${AWS_ACCOUNT_ID}:policy/AzurePutEventsPolicy \
		--profile ${AWS_PROFILE}

azure_oidc:
	STORAGE_CONNECTION_STRING=$$(make -s storage_account_connstring); \
	AZURE_TENANT_ID=$$(make -s azure_tenant_name); \
	az ad app create \
		--display-name AzureEventBridgePublisher \
		--identifier-uris api://${AZURE_TENANT_ID}/AzureEventBridgePublisher \
		--sign-in-audience AzureADMyOrg
	az functionapp identity assign \
		--resource-group ${RG} \
		--name ${FUNC_APP}
# 	Get the identity object id
	MANAGED_IDENTITY_OBJECT_ID=$(az functionapp identity show \
	--resource-group $RG \
	--name $FUNC_APP \
	--query principalId -o tsv)

# 	Grant App Role to Managed Identity
	APP_ID=$(az ad app list \
	--display-name AzureEventBridgePublisher \
	--query '[0].appId' -o tsv)

# 	App Object ID (IMPORTANT – not appId)
	APP_OBJECT_ID=$(az ad app show \
	--id $APP_ID \
	--query id -o tsv)

	az ad app federated-credential create \
		--id ${APP_OBJECT_ID} \
		--parameters '{
			"name": "AzureFunctionCredential",
			"issuer": "https://sts.windows.net/${AZURE_TENANT_ID}",
			"subject": "system:principalid:${MANAGED_IDENTITY_OBJECT_ID}",
			"audiences": ["api://${AZURE_TENANT_ID}/AzureEventBridgePublisher"]
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

clean:
	az functionapp delete --name ${FUNC_APP} --resource-group ${RG}
	az storage account delete --name ${STORAGE} --resource-group ${RG} --yes
	az group delete --name ${RG} --yes --no-wait