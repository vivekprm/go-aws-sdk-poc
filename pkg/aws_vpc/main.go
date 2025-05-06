package main

import (
	"awspoc/pkg/aws_security_group/awssg"
	awsvpc "awspoc/pkg/aws_vpc/vpc_util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
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

	sgID := aws.String("vivekm-custom-sg")
	sid := awssg.CreateSecurityGroup(cli, sgID, vpcID)

	// vpcID := aws.String("vpc-09768d14cd79a0fcf")
	subnet := awsvpc.CreateSubnet(cli, vpcID, "ap-south-1a", "10.0.0.0/24", awsvpc.SubnetTypePrivate)
	rtbID := awsvpc.CreateRouteTable(cli, subnet, vpcID)

	fmt.Printf("subnet 1: %s\n", *subnet)

	vpcInfo := &awsvpc.VpcInfo{
		VpcID:           *vpcID,
		SecurityGroupID: *sid,
		SubnetID:        *subnet,
		RouteTableID:    *rtbID,
	}

	bytes, err := json.Marshal(vpcInfo)
	if err != nil {
		log.Fatalf("Unable to marshal vpcinfo: %v\n", err)
	}
	f, err := os.Create("vpcdata.txt")
	if err != nil {
		log.Fatalf("Unable to create file: %v\n", err)
	}
	n, err := fmt.Fprint(f, string(bytes))
	if err != nil {
		log.Fatalf("Unable to write to file: %v\n", err)
	}
	log.Println("Wrote to the file successfully", n)
	// instance1 := aws.String("i-0c0401346d681bebf")
	// instance2 := aws.String("i-00542ecb7aa2b9d62")
	// subnet1 := aws.String("subnet-0429f21907ddb2d7f")
	// subnet2 := aws.String("subnet-0a8a03099d73d6586")
	// // vpcID := resp.Vpc.VpcId
	// vpcID := aws.String("vpc-06223026e4c5e98b0")
	// cleanUpInstances(pcli, instance1, instance2)
	// cleanUpSubnets(pcli, subnet1, subnet2)
	// cleanUpVpc(pcli, vpcID)
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
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
