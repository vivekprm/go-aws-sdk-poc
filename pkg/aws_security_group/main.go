package main

import (
	"awspoc/pkg/aws_security_group/awssg"
	"context"
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
	vpcID := aws.String("vpc-0f00b3f6f34ac63d9")
	name := aws.String("vm-custom-sg")
	awssg.CreateSecurityGroup(cli, name, vpcID)
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
}
