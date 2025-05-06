package awsigw

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func CreateIgw(cli *ec2.Client) *string {
	ctx := context.Background()
	resp, err := cli.CreateInternetGateway(ctx, &ec2.CreateInternetGatewayInput{
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInternetGateway,
				Tags: []types.Tag{
					{
						Key:   aws.String("created-by"),
						Value: aws.String("vivek"),
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Creation of internet gateway failed: %v\n", err)
	}
	return resp.InternetGateway.InternetGatewayId
}

func AttachIgwToVPC(cli *ec2.Client, igwID *string, vpcID string) {
	resp, err := cli.AttachInternetGateway(context.Background(), &ec2.AttachInternetGatewayInput{
		InternetGatewayId: igwID,
		VpcId:             aws.String(vpcID),
	})
	if err != nil {
		log.Fatalf("Error in attaching IGW to VPC: %v\n", err)
	}
	log.Printf("IGW attached to VPC successfully: %v\n", resp.ResultMetadata)
}

func AttachIgwRouteToSubnet(cli *ec2.Client, igwID *string, rtbID string) {
	CreateRouteToIGW(context.Background(), cli, igwID, aws.String(rtbID), "0.0.0.0/0")
}

func CreateRouteToIGW(ctx context.Context, cli *ec2.Client, igwID, rtbID *string, destCIDR string) {
	resp1, err := cli.CreateRoute(ctx, &ec2.CreateRouteInput{
		RouteTableId:         rtbID,
		DestinationCidrBlock: aws.String(destCIDR),
		GatewayId:            igwID,
	})
	if err != nil {
		log.Fatalf("Error in route creation: %v\n", err)
	}
	log.Printf("Route creation successful: %v\n", resp1.ResultMetadata)
}