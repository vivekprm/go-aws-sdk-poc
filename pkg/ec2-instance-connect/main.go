package main

import (
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
	profile := os.Getenv("AWS_PROFILE_PERSONAL")

	// Loading default
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// Loading a profile
	cli := getClient(profile)
	// profileUS := os.Getenv("AWS_PROFILE_PERSONAL_US")
	// cliUs := getClient(profileUS)

	vpcInfo := awsvpc.GetVPCData()

	instance := awsvpc.CreateInstance(cli, aws.String(vpcInfo.SubnetID), vpcInfo.SecurityGroupID, "ami-0f1dcc636b69a6438")
	fmt.Printf("instance: %s\n", *instance)

	awsvpc.CreateEc2InstanceConnect(cli, aws.String(vpcInfo.SubnetID), []string{vpcInfo.SecurityGroupID})
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
}
