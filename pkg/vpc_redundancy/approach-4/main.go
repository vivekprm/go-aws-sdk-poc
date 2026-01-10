package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	ebtypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v\n", err)
	}
	personalProfile := os.Getenv("AWS_PROFILE_PERSONAL")
	hotmailProfile := os.Getenv("AWS_PROFILE_HOTMAIL")

	ctx := context.Background()
	s3Cli := getS3Client(ctx, hotmailProfile)

	data, err := os.ReadFile("pkg/vpc_redundancy/approach-4/data.json")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	out, err := s3Cli.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(os.Getenv("S3_BUCKET_DATA_KEY")),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		log.Fatalf("Failed to upload file to S3: %v", err)
	}
	log.Printf("File uploaded successfully, ETag: %s", *out.ETag)

	fout, err := s3Cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(os.Getenv("S3_BUCKET_FUNCTION_KEY")),
	})
	if err != nil {
		log.Fatalf("Failed to get lambda code zip file from S3: %v", err)
	}
	defer fout.Body.Close()

	zipData, err := io.ReadAll(fout.Body)
	if err != nil {
		log.Fatalf("Failed to read lambda code from S3: %v", err)
	}

	lambdaCli := getLambdaClient(ctx, personalProfile)

	lout, err := lambdaCli.CreateFunction(ctx, &lambda.CreateFunctionInput{
		FunctionName: aws.String("vpc-redundancy-poc-function"),
		Runtime:      types.RuntimeProvidedal2023,
		Role:         aws.String(os.Getenv("LAMBDA_ROLE_ARN")),
		Handler:      aws.String("bootstrap"),
		Code: &types.FunctionCode{
			ZipFile: zipData,
		},
		Architectures: []types.Architecture{
			types.ArchitectureArm64,
		},
		Environment: &types.Environment{
			Variables: map[string]string{
				"BUCKET_NAME":            os.Getenv("S3_BUCKET_NAME"),
				"BUCKET_KEY":             os.Getenv("S3_BUCKET_DATA_KEY"),
				"ASSUME_ROLE_ARN":        os.Getenv("ASSUME_ROLE_ARN"),
				"SQS_QUEUE_URL":          os.Getenv("REDUNDANCY_EVENTS_QUEUE_URL"),
				"CENTRAL_ACCOUNT_REGION": "us-east-1",
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create Lambda function: %v", err)
	}
	log.Printf("Lambda function created successfully, ARN: %s", *lout.FunctionArn)

	ebCli := getEventBridgeClient(ctx, personalProfile)
	pout, err := ebCli.PutRule(ctx, &eventbridge.PutRuleInput{
		Name:         aws.String("to-vpc-redundancy-lambda"),
		EventBusName: aws.String("default"),
		EventPattern: aws.String("{\"source\":[\"aws.ec2\"], \"detail-type\": [\"EC2 Instance State-change Notification\"]}"),
		State:        ebtypes.RuleStateEnabled,
	})
	if err != nil {
		log.Fatalf("Failed to create EventBridge rule: %v", err)
	}
	log.Printf("EventBridge rule created successfully, ARN: %s", *pout.RuleArn)

	_, err = ebCli.PutTargets(ctx, &eventbridge.PutTargetsInput{
		Rule:         aws.String("to-vpc-redundancy-lambda"),
		EventBusName: aws.String("default"),
		Targets: []ebtypes.Target{
			{
				Id:      aws.String("VPCRedundancyLambdaTarget"),
				Arn:     lout.FunctionArn,
				RoleArn: aws.String("arn:aws:iam::643716337869:role/EventBridgeWriteToLambda"),
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to add target to EventBridge rule: %v", err)
	}
	log.Println("Target added to EventBridge rule successfully")
}

func getEventBridgeClient(ctx context.Context, profile string) *eventbridge.Client {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return eventbridge.NewFromConfig(cfg)
}

func getS3Client(ctx context.Context, profile string) *s3.Client {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return s3.NewFromConfig(cfg)
}

func getLambdaClient(ctx context.Context, profile string) *lambda.Client {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return lambda.NewFromConfig(cfg)
}
