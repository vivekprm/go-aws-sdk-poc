package main

import (
	"awspoc/pkg/aws_internet_gateway/awsigw"
	"awspoc/pkg/aws_security_group/awssg"
	awsvpc "awspoc/pkg/aws_vpc/vpc_util"
	awsnatgw "awspoc/pkg/nat_gateway/aws_natgw"
	"context"
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

	pubSubnet1 := awsvpc.CreateSubnet(cli, vpcID, "ap-south-1a", "10.0.0.0/18", awsvpc.SubnetTypePublic)
	pubrtb1ID := awsvpc.CreateRouteTable(cli, pubSubnet1, vpcID)

	privSubnet1 := awsvpc.CreateSubnet(cli, vpcID, "ap-south-1a", "10.0.64.0/18", awsvpc.SubnetTypePrivate)
	// privrtb1ID := awsvpc.CreateRouteTable(cli, privSubnet1, vpcID)

	pubSubnet2 := awsvpc.CreateSubnet(cli, vpcID, "ap-south-1b", "10.0.128.0/18", awsvpc.SubnetTypePublic)
	pubrtb2ID := awsvpc.CreateRouteTable(cli, pubSubnet2, vpcID)

	privSubnet2 := awsvpc.CreateSubnet(cli, vpcID, "ap-south-1b", "10.0.192.0/18", awsvpc.SubnetTypePrivate)
	privrtb2ID := awsvpc.CreateRouteTable(cli, privSubnet2, vpcID)

	igwID := awsigw.CreateIgw(cli)
	awsigw.AttachIgwToVPC(cli, igwID, *vpcID)
	awsigw.AttachIgwRouteToSubnet(cli, igwID, *pubrtb1ID)
	awsigw.AttachIgwRouteToSubnet(cli, igwID, *pubrtb2ID)

	pubsgID1 := aws.String("vivekm-custom-pub-sg-1")
	pubsid1 := awssg.CreateSecurityGroup(cli, pubsgID1, vpcID)
	awssg.AllowPort(cli, "1", *pubsid1, "10.0.0.0/18", -1, -1)

	privsgID1 := aws.String("vivekm-custom-priv-sg-1")
	privsid1 := awssg.CreateSecurityGroup(cli, privsgID1, vpcID)
	awssg.AllowPort(cli, "1", *privsid1, "10.0.64.0/18", -1, -1)

	pubsgID2 := aws.String("vivekm-custom-pub-sg-2")
	pubsid2 := awssg.CreateSecurityGroup(cli, pubsgID2, vpcID)
	awssg.AllowPort(cli, "1", *pubsid2, "10.0.128.0/18", -1, -1)

	privsgID2 := aws.String("vivekm-custom-priv-sg-2")
	privsid2 := awssg.CreateSecurityGroup(cli, privsgID2, vpcID)
	awssg.AllowPort(cli, "1", *privsid2, "10.0.0.0/16", -1, -1)

	awsvpc.CreateInstance(cli, pubSubnet1, *pubsid1, "ami-0f1dcc636b69a6438")
	awsvpc.CreateInstance(cli, privSubnet1, *privsid1, "ami-0f1dcc636b69a6438")
	awsvpc.CreateEc2InstanceConnect(cli, privSubnet1, []string{*privsid1})

	awsvpc.CreateInstance(cli, pubSubnet2, *pubsid2, "ami-0f1dcc636b69a6438")
	awsvpc.CreateInstance(cli, privSubnet2, *privsid2, "ami-0f1dcc636b69a6438")

	ngwid2 := awsnatgw.CreateNatGateway(cli, *privSubnet2, types.ConnectivityTypePrivate)
	time.Sleep(5 * time.Second)
	out, err := cli.CreateRoute(context.Background(), &ec2.CreateRouteInput{
		RouteTableId:         privrtb2ID,
		NatGatewayId:         ngwid2,
		DestinationCidrBlock: aws.String("10.0.192.0/18"),
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
