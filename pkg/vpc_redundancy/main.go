package main

import (
	"awspoc/pkg/aws_security_group/awssg"
	awsvpc "awspoc/pkg/aws_vpc/vpc_util"
	"awspoc/pkg/vpc_redundancy/receiver"
	"awspoc/pkg/vpc_redundancy/utils"
	"bytes"
	"context"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
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
	region := os.Getenv("AWS_REGION")
	az := os.Getenv("AWS_REGION_AZ")
	log.Printf("Using profile %s and region %s\n", profile, region)
	// Loading default
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// Loading a profile
	cli := utils.GetEC2Client(profile, region)
	ctx := context.Background()

	cfCli := utils.GetCloudFormationClient(profile, region)

	content, err := os.ReadFile("pkg/vpc_redundancy/vpc_redundancy_cf_tpl.yaml")
	if err != nil {
		log.Fatalln("error reading CloudFormation template file", err)
	}

	tmpl, err := template.New("cloudfront-template").Parse(string(content))

	if err != nil {
		log.Fatalln("error in parsing template", err)
	}

	var tmplOut bytes.Buffer
	err = tmpl.Execute(&tmplOut, struct {
		SiteName string
	}{SiteName: "redundancy-demo-site"})

	if err != nil {
		log.Fatalln("error in executing template", err)
	}

	log.Println("Creating and executing cloudformation stack")
	stackout, err := cfCli.CreateStack(ctx, &cloudformation.CreateStackInput{
		StackName:    aws.String("vivekm-redundancy-test-stack"),
		TemplateBody: aws.String(tmplOut.String()),
		// DisableRollback: aws.Bool(true),
		Tags: []types.Tag{
			{
				Key:   aws.String("SiteName"),
				Value: aws.String("redundancy-demo-site"),
			},
			{
				Key:   aws.String("PocType"),
				Value: aws.String("vpc-redundancy"),
			},
		},
	})

	if err != nil {
		log.Fatalln("error in creating stack", err)
	}
	log.Println("Stack creation completed, waiting for stack execution to complete", *stackout.StackId)

	cfWaiter := cloudformation.NewStackCreateCompleteWaiter(cfCli)

	if err = cfWaiter.Wait(ctx, &cloudformation.DescribeStacksInput{
		StackName: stackout.StackId,
	}, 10*time.Minute); err != nil {
		log.Fatalln("error waiting for stack creation to complete", err)
	}

	log.Println("Stack execution completed.")

	vpc1 := awsvpc.CreateVpc(cli, "192.168.0.0/16")
	subnet1 := awsvpc.CreateSubnet(cli, vpc1.Vpc.VpcId, az, "192.168.0.0/24", awsvpc.SubnetTypePrivate)
	rtb1 := awsvpc.CreateRouteTable(cli, subnet1, vpc1.Vpc.VpcId)
	sg1 := awssg.CreateSecurityGroup(cli, aws.String("vivekm-test-sg"), vpc1.Vpc.VpcId)
	awsvpc.CreateNatGateway(ctx, cli, subnet1)

	out, err := cli.AssociateRouteTable(ctx, &ec2.AssociateRouteTableInput{
		RouteTableId: rtb1,
		SubnetId:     subnet1,
	})

	if err != nil {
		log.Fatalln("error in assoicating route table", err)
	}

	log.Printf("Associated route table %s with subnet %s, association id %s\n", *rtb1, *subnet1, *out.AssociationId)

	awsvpc.CreateInstance(cli, subnet1, *sg1, "ami-00ca570c1b6d79f36")

	str, err := cfCli.DescribeStackResources(ctx, &cloudformation.DescribeStackResourcesInput{
		StackName: stackout.StackId,
	})

	if err != nil {
		log.Fatalln("error in describing stack", err)
	}

	var queueID string
	for _, resource := range str.StackResources {
		// log.Printf("Resource Type: %s, Resource ID: %s\n", *resource.ResourceType, *resource.PhysicalResourceId)

		if *resource.ResourceType == "AWS::SQS::Queue" {
			queueID = *resource.PhysicalResourceId
		}
	}

	log.Println("Starting message consumer...")
	startConsumer(ctx, profile, region, queueID)
}

func startConsumer(ctx context.Context, profile, region, queueID string) {
	cli := utils.GetSQSClient(profile, region)

	msgChan := make(chan receiver.Message)
	go receiver.Consume(cli, msgChan, queueID)

	log.Println("Started message consumer, waiting for messages...")

	for msg := range msgChan {
		log.Println("Received message:", msg.Body)
	}
}
