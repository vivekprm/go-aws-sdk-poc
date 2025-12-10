package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v\n", err)
	}
	profileUS := os.Getenv("AWS_PROFILE")
	cli := getClient(profileUS)
	out, err := cli.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/643716337869/vpc-redundancy-queue"),
	}) 

	if err != nil {
		log.Fatalln("error in receiving msg", err)
	}
	log.Println(out.Messages)
	for _, msg := range out.Messages {
		log.Println(*msg.Body)
	}
}

func getClient(profile string) *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return sqs.NewFromConfig(cfg)
}
