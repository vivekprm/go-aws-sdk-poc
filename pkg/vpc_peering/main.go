package main

import (
	"awspoc/pkg/aws_security_group/awssg"
	awsvpc "awspoc/pkg/aws_vpc/vpc_util"
	"awspoc/pkg/vpc_peering/peerutil"
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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
	profileUS := os.Getenv("AWS_PROFILE_PERSONAL_US")
	// Loading default
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// Loading a profile
	cli := getClient(profile)
	cliUs := getClient(profileUS)

	vpcInfo := awsvpc.GetVPCData()

	vpc1 := awsvpc.CreateVpc(cli, "192.168.0.0/16")
	subnet1 := awsvpc.CreateSubnet(cli, vpc1.Vpc.VpcId, "ap-south-1a", "192.168.0.0/24")
	rtb1 := awsvpc.CreateRouteTable(cli, subnet1, vpc1.Vpc.VpcId)
	sg1 := awssg.CreateSecurityGroup(cli, aws.String("vivekm-ap-s1c-sg"), vpc1.Vpc.VpcId)

	vpc2 := awsvpc.CreateVpc(cliUs, "172.16.0.0/16")
	subnet2 := awsvpc.CreateSubnet(cliUs, vpc2.Vpc.VpcId, "us-east-1a", "172.16.0.0/24")
	rtb2 := awsvpc.CreateRouteTable(cliUs, subnet2, vpc2.Vpc.VpcId)
	sg2 := awssg.CreateSecurityGroup(cliUs, aws.String("vivekm-us-s1a-sg"), vpc2.Vpc.VpcId)

	connection1 := peerutil.CreatePeering(cli, aws.String(vpcInfo.VpcID), vpc1.Vpc.VpcId, "ap-south-1")
	connection2 := peerutil.CreatePeering(cli, vpc1.Vpc.VpcId, vpc2.Vpc.VpcId, "us-east-1")

	peerutil.CreatePeeringRoute(cli, vpcInfo.RouteTableID, "192.168.0.0/16", *connection1)
	awssg.AllowPort(cli, "1", vpcInfo.SecurityGroupID, "192.168.0.0/16", -1, -1)

	peerutil.CreatePeeringRoute(cli, *rtb1, "10.0.0.0/16", *connection1)
	awssg.AllowPort(cli, "1", *sg1, "10.0.0.0/16", -1, -1)

	peerutil.CreatePeeringRoute(cli, *rtb1, "172.16.0.0/16", *connection2)
	awssg.AllowPort(cli, "1", *sg1, "172.16.0.0/16", -1, -1)

	peerutil.CreatePeeringRoute(cliUs, *rtb2, "192.168.0.0/16", *connection2)
	awssg.AllowPort(cliUs, "1", *sg2, "192.168.0.0/16", -1, -1)

	awsvpc.CreateInstance(cli, subnet1, *sg1, "ami-0f1dcc636b69a6438")
	awsvpc.CreateInstance(cliUs, subnet2, *sg2, "ami-0e449927258d45bc4")
	time.Sleep(10 * time.Second)
	awsvpc.CreateEc2InstanceConnect(cli, subnet1, []string{*sg1})
	awsvpc.CreateEc2InstanceConnect(cliUs, subnet2, []string{*sg2})
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
}
