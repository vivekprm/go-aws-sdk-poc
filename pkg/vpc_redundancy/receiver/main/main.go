package main

import (
	"awspoc/pkg/vpc_redundancy/receiver"
	"awspoc/pkg/vpc_redundancy/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v\n", err)
	}
	profile := os.Getenv("AWS_PROFILE_PERSONAL")
	region := os.Getenv("AWS_REGION")
	cli := utils.GetSQSClient(profile, region)

	queueURL := os.Getenv("REDUNDANCY_EVENTS_QUEUE_URL")	

	msgChan := make(chan receiver.Message)
	go receiver.Consume(cli, msgChan, queueURL)

	for msg := range msgChan {
		log.Println("Received message:", msg.Body)
	}
}