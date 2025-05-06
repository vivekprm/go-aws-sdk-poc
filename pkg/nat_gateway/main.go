package main

import (
	"awspoc/pkg/aws_security_group/awssg"
	awsvpc "awspoc/pkg/aws_vpc/vpc_util"
	awsnatgw "awspoc/pkg/nat_gateway/aws_natgw"
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	resp := awsvpc.CreateVpc(cli, "10.0.0.0/16")
	vpcID := resp.Vpc.VpcId

	// vpcID := aws.String("vpc-09768d14cd79a0fcf")
	subnet1 := awsvpc.CreateSubnet(cli, vpcID, "ap-south-1a", "10.0.0.0/25", awsvpc.SubnetTypePrivate)
	subnet2 := awsvpc.CreateSubnet(cli, vpcID, "ap-south-1b", "10.0.0.128/25", awsvpc.SubnetTypePrivate)
	rtb1ID := awsvpc.CreateRouteTable(cli, subnet1, vpcID)
	// rtb2ID := awsvpc.CreateRouteTable(cli, subnet2, vpcID)

	sgID1 := aws.String("vivekm-custom-sg-1")
	sid1 := awssg.CreateSecurityGroup(cli, sgID1, vpcID)
	awssg.AllowPort(cli, "1", *sid1, "10.0.0.0/25", -1, -1)

	sgID2 := aws.String("vivekm-custom-sg-2")
	sid2 := awssg.CreateSecurityGroup(cli, sgID2, vpcID)
	awssg.AllowPort(cli, "1", *sid2, "10.0.0.128/25", -1, -1)

	awsvpc.CreateInstance(cli, subnet1, *sid1, "ami-0f1dcc636b69a6438")
	fmt.Printf("Security group id is: %s\n", *sid1)
	awsvpc.CreateEc2InstanceConnect(cli, subnet1, []string{*sid1})

	awsvpc.CreateInstance(cli, subnet2, *sid2, "ami-0f1dcc636b69a6438")
	// awsvpc.CreateEc2InstanceConnect(cli, subnet2, []string{*sid})

	// ngwid1 := awsnatgw.CreateNatGateway(cli, *subnet1, types.ConnectivityTypePrivate)
	ngwid2 := awsnatgw.CreateNatGateway(cli, *subnet2, types.ConnectivityTypePrivate)
	time.Sleep(5 * time.Second)
	out, err := cli.CreateRoute(context.Background(), &ec2.CreateRouteInput{
		RouteTableId:         rtb1ID,
		NatGatewayId:         ngwid2,
		DestinationCidrBlock: aws.String("10.0.0.128/25"),
	})
	if err != nil {
		log.Fatalf("error in creating route to natgateway: %v\n", err)
	}
	log.Printf("Route to nat gateway created successfully: %v\n", out.ResultMetadata)
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
}
