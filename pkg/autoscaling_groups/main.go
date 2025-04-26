package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
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
	profile := os.Getenv("AWS_PROFILE")
	pprofile := os.Getenv("AWS_PROFILE_PERSONAL")

	// Loading default
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// Loading a profile
	cli := getClient(profile)
	pcli := getClient(pprofile)

	createAutoscalingGroup(pcli)
	deleteAutoscalingGroup(pcli)

	// Build the request with its input parameters
	listAutoScalingGroups(cli)
}

func getClient(profile string) *autoscaling.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return autoscaling.NewFromConfig(cfg)
}

func createAutoscalingGroup(cli *autoscaling.Client) {
	resp, err := cli.CreateAutoScalingGroup(context.TODO(), &autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String("my-asg"),
		LaunchTemplate: &types.LaunchTemplateSpecification{
			LaunchTemplateName: aws.String("my-template-for-auto-scaling"),
		},
		MinSize:           aws.Int32(1),
		MaxSize:           aws.Int32(4),
		DesiredCapacity:   aws.Int32(1),
		VPCZoneIdentifier: aws.String("subnet-0163bd2b8c13cb139,subnet-0747a2a8e6aa1b256,subnet-014251eee77cce8d1"),
	})
	if err != nil {
		log.Fatalf("Unable to create autoscaling group: %v", err)
	}
	fmt.Printf("Creation of autoscaling group successful: %v\n", resp.ResultMetadata)
}

func deleteAutoscalingGroup(cli *autoscaling.Client) {
	resp, err := cli.DeleteAutoScalingGroup(context.TODO(), &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: aws.String("my-asg"),
		ForceDelete:          aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Autoscaling group deleted successfully: %v\n", resp)
}

func listAutoScalingGroups(cli *autoscaling.Client) {
	resp, err := cli.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("failed to list templates, %v", err)
	}

	for i, asg := range resp.AutoScalingGroups {
		fmt.Printf("Autoscaling Groups %d\n", i)
		fmt.Println("--------------------------------")
		fmt.Printf("Name: %s\n", *asg.AutoScalingGroupName)
		if asg.LaunchConfigurationName != nil {
			fmt.Printf("Launch Configuration Name: %s\n\n", *asg.LaunchConfigurationName)
		}
	}
}
