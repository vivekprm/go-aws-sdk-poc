package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	// --- 1. Configure AWS credentials ---
	profile := os.Getenv("AWS_PROFILE_PERSONAL")
	cli := getClient(profile)

	ctx := context.Background()

	resp, err := cli.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-0df30fbc669b7b39a"),
		InstanceType: types.InstanceTypeM52xlarge,
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		SubnetId:     aws.String("subnet-06e2c4f6e8247f4e2"),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags: []types.Tag{
					{
						Key:   aws.String("Environment"),
						Value: aws.String("demo1"),
					},
					{
						Key:   aws.String("CostCenter"),
						Value: aws.String("7929"),
					},
					{
						Key:   aws.String("ManagerEmail"),
						Value: aws.String("v.enugurthi"),
					},
					{
						Key:   aws.String("Team"),
						Value: aws.String("f5xc-orchestration"),
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Unable to create instance %v\n", err)
	}
	log.Printf("Instance created successfully: %s", *resp.Instances[0].InstanceId)

	// --- 3. Wait until the instance is running ---
	waiter := ec2.NewInstanceRunningWaiter(cli)
	fmt.Println("Waiting for instance to be running...")

	err = waiter.Wait(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{*resp.Instances[0].InstanceId},
	}, 5*time.Minute)
	if err != nil {
		log.Fatalf("failed waiting for instance to run: %v", err)
	}

	// --- 4. Describe the instance to get its public IP ---
	describeOut, err := cli.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{*resp.Instances[0].InstanceId},
	})
	if err != nil {
		log.Fatalf("failed to describe instance: %v", err)
	}

	instance := describeOut.Reservations[0].Instances[0]
	fmt.Printf("Instance %s is running!\n", *instance.InstanceId)
	fmt.Printf("Public IP: %s\n", aws.ToString(instance.PublicIpAddress))
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
}
