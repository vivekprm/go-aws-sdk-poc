#!/usr/bin/env bash

export LOCATION=eastus2
export RG=vm-events-$RANDOM
export STORAGE=vmfuncstorage$RANDOM
export FUNC_APP=send-event-func-$RANDOM
export ACR_NAME=eventpocacr$RANDOM
export TARGETPLATFORM=linux/arm64
export IDENTITY_NAME=uami-azure-events-dev-$RANDOM

az group create \
  --name $RG \
  --location $LOCATION

# Register your subscription with Microsoft.Storage Resource Provider
# Portal-> < Your-Subscription> -> Resource providers -> search Microsoft. Storage -> Register
az storage account create --name $STORAGE --location $LOCATION --resource-group $RG --sku Standard_LRS

az identity create \
  --name $IDENTITY_NAME \
  --resource-group $RG \
  --location $LOCATION

CLIENT_ID=$(az identity show \
  --name $IDENTITY_NAME \
  --resource-group $RG \
  --query "{clientId:clientId,tenantId:tenantId,principalId:principalId}" | jq ".clientId")
TENANT_ID=$(az identity show \
  --name $IDENTITY_NAME \
  --resource-group $RG \
  --query "{clientId:clientId,tenantId:tenantId,principalId:principalId}" | jq ".tenantId")
  
# Register your subscription with Microsoft.ContainerRegistry Resource Provider
cd ../../../../../py-azure-sendevent-func

# Register your subscription with Microsoft.Web Resource Provider
az functionapp create \
  --name $FUNC_APP \
  --resource-group $RG \
  --storage-account $STORAGE \
  --consumption-plan-location $LOCATION \
  --runtime python \
  --runtime-version 3.10 \
  --functions-version 4 \
  --os-type Linux

STORAGE_CONNECTION_STRING=$(az storage account show-connection-string \
  --name $STORAGE \
  --resource-group $RG \
  --query connectionString \
  -o tsv
)

az functionapp config appsettings set \
  --resource-group $RG \
  --name $FUNC_APP \
  --settings \
    AWS_ROLE_ARN=arn:aws:iam::665096241598:role/AzureEventBridgeRole \
    AWS_REGION=$LOCATION \
    AWS_EVENT_BUS=central-event-management-bus \
    AZURE_TENANT_ID=$TENANT_ID \
    AZURE_EXTERNAL_ID=azure-eventgrid

az functionapp restart \
  -g $RG \
  -n $FUNC_APP

az functionapp show \
  -g $RG \
  -n $FUNC_APP \
  --query state

# Needs dotnet sdk installed
# In Mac: brew install --cask dotnet-sdk
func azure functionapp publish $FUNC_APP

az functionapp function show \
  --resource-group $RG \
  --name $FUNC_APP \
  --query "invokeUrlTemplate" \
  -o tsv

EVENTGRID_KEY=$(az functionapp keys list \
  -g $RG \
  -n $FUNC_APP \
  --query systemKeys.eventgrid_extension \
  -o tsv)

az monitor diagnostic-settings create \
  --name activity-to-eventgrid \
  --resource-group $RG \
  --storage-account $STORAGE \
  --resource /subscriptions/22e88069-5977-4fc2-9d6e-ef7e14cfbb70 \
  --logs '[{"category":"Administrative","enabled":true}]'

CONNECTION_STRING=$(az monitor app-insights component show \
  -g $RG \
  -a $FUNC_APP \
  --query connectionString)

az functionapp config appsettings set \
  -g $RG \
  -n $FUNC_APP \
  --settings APPLICATIONINSIGHTS_CONNECTION_STRING=$CONNECTION_STRING  

# Register your subscription with Microsoft.EventGrid Resource Provider
az eventgrid event-subscription create \
  --name send-vm-events-to-aws \
  --source-resource-id /subscriptions/22e88069-5977-4fc2-9d6e-ef7e14cfbb70 \
  --endpoint "https://$FUNC_APP.azurewebsites.net/runtime/webhooks/eventgrid?functionName=EventGridTrigger&code=$EVENTGRID_KEY" \
  --included-event-types \
    Microsoft.Compute.VirtualMachines.Write \
    Microsoft.Compute.VirtualMachines.Deallocate

az vm create \
  --resource-group $RG \
  --name vm-monitor-01 \
  --image Ubuntu2204 \
  --size Standard_DC1s_v3 \
  --admin-username azureuser \
  --generate-ssh-keys \
  --tags f5xc-site-name=azure-vm-events

# OIDC Provider
AZURE_TENANT_ID=$(az account show --query tenantId -o tsv)
AZURE_THUMBPRINT=$(openssl x509 -in ms-root.pem -fingerprint -sha1 -noout | cut -d= -f2 | tr -d ':' | tr 'A-F' 'a-f')

aws iam create-open-id-connect-provider \
  --url https://sts.windows.net/$AZURE_TENANT_ID/ \
  --client-id-list api://$AZURE_TENANT_ID/AzureEventBridgePublisher \
  --thumbprint-list $AZURE_THUMBPRINT \
  --profile AWS_HOTMAIL


aws iam create-role \
  --role-name AzureEventBridgePublisher \
  --profile AWS_HOTMAIL \
  --assume-role-policy-document '{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::665096241598:oidc-provider/sts.windows.net/0f4a9737-a4ba-48bd-8174-f54ed1e247e0/"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "sts.windows.net/0f4a9737-a4ba-48bd-8174-f54ed1e247e0/:aud":
          "api://0f4a9737-a4ba-48bd-8174-f54ed1e247e0/AzureEventBridgePublisher"
        }
      }
    }
  ]
}
'


# Create policy
aws iam create-policy \
  --policy-name AzurePutEventsPolicy \
  --profile AWS_HOTMAIL \
  --policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": "events:PutEvents",
        "Resource": "arn:aws:events:us-east-1:665096241598:event-bus/central-event-management-bus"
      }
    ]
  }'

# Attach policy to role
aws iam attach-role-policy \
  --role-name AzureEventBridgePublisher \
  --policy-arn arn:aws:iam::665096241598:policy/AzurePutEventsPolicy \
  --profile AWS_HOTMAIL

# Thumbprint is Microsoft’s root CA (stable & widely used)
# ENV Vars needed
# AWS_ROLE_ARN=arn:aws:iam::665096241598:role/AzureEventBridgePublisher
# AWS_WEB_IDENTITY_TOKEN_FILE=/var/run/secrets/azure/tokens/azure-identity-token
# AWS_REGION=us-east-1
# EVENT_BUS_NAME=central-event-management-bus

# Create Azure AD App (OIDC Audience)
az ad app create \
  --display-name AzureEventBridgePublisher \
  --identifier-uris api://$AZURE_TENANT_ID/AzureEventBridgePublisher \
  --sign-in-audience AzureADMyOrg

# Enable Managed Identity on Azure Function
az functionapp identity assign \
  --resource-group $RG \
  --name $FUNC_APP

# Get the identity object id
MANAGED_IDENTITY_OBJECT_ID=$(az functionapp identity show \
  --resource-group $RG \
  --name $FUNC_APP \
  --query principalId -o tsv)

# Grant App Role to Managed Identity
APP_ID=$(az ad app list \
  --display-name AzureEventBridgePublisher \
  --query '[0].appId' -o tsv)

# App Object ID (IMPORTANT – not appId)
APP_OBJECT_ID=$(az ad app show \
  --id $APP_ID \
  --query id -o tsv)

az ad app federated-credential create \
  --id $APP_OBJECT_ID \
  --parameters '{
    "name": "AzureFunctionCredential",
    "issuer": "https://sts.windows.net/'$AZURE_TENANT_ID'",
    "subject": "system:principalid:'$MANAGED_IDENTITY_OBJECT_ID'",
    "audiences": ["api://'$AZURE_TENANT_ID'/AzureEventBridgePublisher"]
  }'

# Set Azure Function Environment Variables
az functionapp config appsettings set \
  --resource-group $RG \
  --name $FUNC_APP \
  --settings \
    AWS_ROLE_ARN=arn:aws:iam::665096241598:role/AzureEventBridgePublisher \
    AWS_WEB_IDENTITY_TOKEN_FILE=/var/run/secrets/azure/tokens/azure-identity-token \
    AWS_REGION=us-east-1 \
    EVENT_BUS_NAME=central-event-management-bus

# Cleanup
az functionapp delete --name $FUNC_APP --resource-group $RG
az group delete --name $RG --location $LOCATION