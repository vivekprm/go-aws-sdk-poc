package main

import (
	"awspoc/pkg/aws_security_group/awssg"
	awsvpc "awspoc/pkg/aws_vpc/vpc_util"
	"context"
	"fmt"
	"log"
	"os"

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
	profileUS := os.Getenv("AWS_PROFILE")
	// Loading default
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// Loading a profile
	cli := getClient(profileUS)

	vpc1 := awsvpc.CreateVpc(cli, "192.168.0.0/16")
	subnet1 := awsvpc.CreateSubnet(cli, vpc1.Vpc.VpcId, "us-east-1a", "192.168.0.0/24", awsvpc.SubnetTypePrivate)
	rtb1 := awsvpc.CreateRouteTable(cli, subnet1, vpc1.Vpc.VpcId)
	sg1 := awssg.CreateSecurityGroup(cli, aws.String("vivekm-ap-s1c-sg"), vpc1.Vpc.VpcId)

	out, err := cli.AssociateRouteTable(context.Background(), &ec2.AssociateRouteTableInput{
		RouteTableId: rtb1,
		SubnetId: subnet1,
	})

	if err != nil {
		log.Fatalln("error in assoicating route table", err)
	}

	fmt.Println(out)

	awsvpc.CreateInstance(cli, subnet1, *sg1, "ami-068c0051b15cdb816")
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
}
