package main

import (
	"awspoc/pkg/aws_internet_gateway/awsigw"
	awsvpc "awspoc/pkg/aws_vpc/vpc_util"
	"context"
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
	igwID := awsigw.CreateIgw(cli)
	awsigw.AttachIgwToVPC(cli, igwID, vpcID)
	awsigw.AttachIgwRouteToSubnet(cli, igwID, rtbID)
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
