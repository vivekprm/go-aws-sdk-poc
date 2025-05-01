package main

import (
	awsvpc "awspoc/pkg/aws_vpc/vpc_util"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v\n", err)
	}
	// Using the SDK's default configuration, load additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	profile := os.Getenv("AWS_PROFILE_PERSONAL")

	// Loading default
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// Loading a profile
	cli := getClient(profile)

	vpcInfo := awsvpc.GetVPCData()

	vpcID := vpcInfo.VpcID
	subnetID := vpcInfo.SubnetID
	rtbID := vpcInfo.RouteTableID
	sgID := vpcInfo.SecurityGroupID

	modifySubnetAttribute(cli, subnetID)
	igwID := createIgw(cli)
	attachIgwToVPC(cli, igwID, vpcID)
	attachIgwRouteToSubnet(cli, igwID, rtbID)
	awsvpc.CreateInstance(cli, aws.String(subnetID), sgID, "ami-0f1dcc636b69a6438")
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
}

func modifySubnetAttribute(cli *ec2.Client, subnetID string) {
	out, err := cli.ModifySubnetAttribute(context.Background(), &ec2.ModifySubnetAttributeInput{
		SubnetId: aws.String(subnetID),
		MapPublicIpOnLaunch: &types.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
	})
	if err != nil {
		log.Fatalf("Error in modifying subnet attribute: %v\n", err)
	}
	log.Printf("Modification of subnet attribute successful: %v\n", out.ResultMetadata)
}

func createIgw(cli *ec2.Client) *string {
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

func attachIgwToVPC(cli *ec2.Client, igwID *string, vpcID string) {
	resp, err := cli.AttachInternetGateway(context.Background(), &ec2.AttachInternetGatewayInput{
		InternetGatewayId: igwID,
		VpcId:             aws.String(vpcID),
	})
	if err != nil {
		log.Fatalf("Error in attaching IGW to VPC: %v\n", err)
	}
	log.Printf("IGW attached to VPC successfully: %v\n", resp.ResultMetadata)
}

func attachIgwRouteToSubnet(cli *ec2.Client, igwID *string, rtbID string) {
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
