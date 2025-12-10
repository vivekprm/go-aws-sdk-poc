import json
import logging
import boto3

# Configure logging
logger = logging.getLogger()
logger.setLevel(logging.INFO)

ec2 = boto3.client('ec2')

def lambda_handler(event, context):
    """
    Receives and processes an event from Amazon EventBridge.

    Args:
        event (dict): The EventBridge event object.
        context (object): The Lambda context object.
    """
    logger.info(f"Received EventBridge event: {json.dumps(event)}")

    # Extract specific details from the event, if needed
    # For example, if the event contains a 'detail' field with custom data:
    if 'detail' in event:
        custom_data = event['detail']
        logger.info(f"Custom data from event detail: {json.dumps(custom_data)}")
        # Process custom_data as required by your application logic
        if custom_data.get("state") == "stopped":
            logger.info("The state is stopped. Performing necessary actions.")
            route_table_id = "rtb-001224d7ab6518238"
            destination_cidr_block = "0.0.0.0/0"
            target_id = "eni-0bc7f97d9263339ee"

            try:
                # First, check if the route already exists and replace it if necessary
                # or simply add it if it doesn't exist.
                # For simplicity, this example focuses on creating/replacing.
                # You might want to delete existing routes before creating new ones if needed.

                response = ec2.replace_route(
                    DestinationCidrBlock=destination_cidr_block,
                    RouteTableId=route_table_id,
                    NetworkInterfaceId=target_id # Or NatGatewayId, VpcPeeringConnectionId, etc.
                )
                print(f"Route replaced successfully: {response}")

                return {
                    'statusCode': 200,
                    'body': 'Route updated successfully!'
                }
            except ec2.exceptions.ClientError as e:
                if "InvalidRoute.NotFound" in str(e):
                    # If route not found, attempt to create it
                    try:
                        response = ec2.create_route(
                            DestinationCidrBlock=destination_cidr_block,
                            RouteTableId=route_table_id,
                            NetworkInterfaceId=target_id # Or NatGatewayId, VpcPeeringConnectionId, etc.
                        )
                        print(f"Route created successfully: {response}")
                        return {
                            'statusCode': 200,
                            'body': 'Route created successfully!'
                        }
                    except Exception as create_e:
                        print(f"Error creating route: {create_e}")
                        return {
                            'statusCode': 500,
                            'body': f'Error creating route: {str(create_e)}'
                        }
                else:
                    print(f"Error updating route: {e}")
                    return {
                        'statusCode': 500,
                        'body': f'Error updating route: {str(e)}'
                    }
            except Exception as e:
                print(f"An unexpected error occurred: {e}")
                return {
                    'statusCode': 500,
                    'body': f'An unexpected error occurred: {str(e)}'
                }
        elif custom_data.get("state") == "terminating":
            logger.info("The state is terminating. Performing necessary actions.")
        elif custom_data.get("state") == "terminated":
            logger.info("The state is terminated. Performing necessary actions.")
        else:
            logger.info("The state is %s. No actions taken.", custom_data.get("state"))

    # You can also access other fields like source, detail-type, time, etc.
    source = event.get('source')
    detail_type = event.get('detail-type')
    time = event.get('time')

    logger.info(f"Event Source: {source}")
    logger.info(f"Detail Type: {detail_type}")
    logger.info(f"Event Time: {time}")

    # Your application logic to process the EventBridge message goes here
    # For example, you might interact with other AWS services, perform calculations, etc.

    return {
        'statusCode': 200,
        'body': json.dumps('Event processed successfully!')
    }