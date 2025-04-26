package main

import (
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
	profile := os.Getenv("AWS_PROFILE")
	pprofile := os.Getenv("AWS_PROFILE_PERSONAL")

	// Loading default
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// Loading a profile
	_ = getClient(profile)
	pcli := getClient(pprofile)

	resp := createVpc(pcli)
	subnet1, subnet2 := createSubnet(pcli, resp.Vpc.VpcId)
	instance1, instance2 := createInstances(pcli, subnet1, subnet2)

	cleanUpInstances(pcli, instance1, instance2)
	cleanUpSubnets(pcli, subnet1, subnet2)
	cleanUpVpc(pcli, resp.Vpc.VpcId)
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
}

func createVpc(cli *ec2.Client) *ec2.CreateVpcOutput {
	ctx := context.Background()
	resp, err := cli.CreateVpc(ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/16"),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeVpc,
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
		log.Fatal(err)
	}
	out, err := cli.AssociateVpcCidrBlock(ctx, &ec2.AssociateVpcCidrBlockInput{
		VpcId:     resp.Vpc.VpcId,
		CidrBlock: aws.String("100.64.1.0/24"),
	})
	if err != nil {
		log.Fatalf("Error in associating CIDR block: %v\n", err)
	}
	log.Printf("CIDR associated successfully: %v\n", *out.VpcId)
	return resp
}

func createSubnet(cli *ec2.Client, vpcID *string) (*string, *string) {
	resp1, err := cli.CreateSubnet(context.TODO(), &ec2.CreateSubnetInput{
		VpcId:            vpcID,
		AvailabilityZone: aws.String("us-east-1a"),
		CidrBlock:        aws.String("10.0.0.0/24"),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSubnet,
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
		log.Fatalf("Creation of subnet1 failed: %v\n", err)
	}

	log.Printf("Creation of subnet1 successful: %v\n", *resp1.Subnet.SubnetId)

	resp2, err := cli.CreateSubnet(context.TODO(), &ec2.CreateSubnetInput{
		VpcId:            vpcID,
		AvailabilityZone: aws.String("us-east-1b"),
		CidrBlock:        aws.String("100.64.1.0/28"),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSubnet,
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
		log.Fatalf("Creation of subnet2 failed: %v\n", err)
	}

	log.Printf("Creation of subnet2 successful: %v\n", *resp2.Subnet.SubnetId)
	return resp1.Subnet.SubnetId, resp2.Subnet.SubnetId
}

func createInstances(cli *ec2.Client, subnet1Id, subnet2Id *string) (*string, *string) {
	ctx := context.Background()
	resp1, err := cli.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-0e449927258d45bc4"),
		InstanceType: types.InstanceTypeT2Micro,
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		SubnetId:     aws.String(*subnet1Id),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
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
		log.Fatalf("Unable to create instance 1 %v\n", err)
	}
	log.Printf("Instance 1 created successfully: %s", *resp1.Instances[0].InstanceId)

	resp2, err := cli.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-0e449927258d45bc4"),
		InstanceType: types.InstanceTypeT2Micro,
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		SubnetId:     aws.String(*subnet2Id),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
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
		log.Fatalf("Unable to create instance 2 %v\n", err)
	}
	log.Printf("Instance 2 created successfully: %s", *resp2.Instances[0].InstanceId)
	return resp1.Instances[0].InstanceId, resp2.Instances[0].InstanceId
}

func cleanUpInstances(cli *ec2.Client, instance1ID, instance2ID *string) {
	ctx := context.Background()
	out, err := cli.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{*instance1ID, *instance2ID},
	})
	if err != nil {
		log.Fatalf("Unable to delete instances: %v\n", err)
	}
	log.Printf("Deletion of instances successful: %v\n", out.ResultMetadata)
}
func cleanUpSubnets(cli *ec2.Client, subnet1, subnet2 *string) {
	ctx := context.Background()
	resp1, err := cli.DeleteSubnet(ctx, &ec2.DeleteSubnetInput{
		SubnetId: subnet1,
	})
	if err != nil {
		log.Fatalf("Deletion of subnet 1 failed: %v\n", err)
	}
	log.Printf("Subnet 1 deleted successfully: %v\n", resp1.ResultMetadata)
	resp2, err := cli.DeleteSubnet(ctx, &ec2.DeleteSubnetInput{
		SubnetId: subnet2,
	})
	if err != nil {
		log.Fatalf("Deletion of subnet 2 failed: %v\n", err)
	}
	log.Printf("Subnet 2 deleted successfully: %v\n", resp2.ResultMetadata)
}
func cleanUpVpc(cli *ec2.Client, vpcId *string) {
	resp, err := cli.DeleteVpc(context.Background(), &ec2.DeleteVpcInput{
		VpcId: vpcId,
	})
	if err != nil {
		log.Fatalf("Deletion of vpc failed: %v\n", err)
	}
	log.Printf("VPC deleted successfully: %v\n", resp.ResultMetadata)
}
