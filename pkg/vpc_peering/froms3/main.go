package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	cli := getClient(profile)
	dataFile, err := cli.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("redundancy-data-poc"),
		Key:    aws.String("data.json"),
	})
	if err != nil {
		log.Printf("Error reading data file from S3: %v\n", err)
	}
	log.Printf("Data file downloaded successfully: %v\n", dataFile)
	defer dataFile.Body.Close()
}

func getClient(profile string) *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return s3.NewFromConfig(cfg)
}
