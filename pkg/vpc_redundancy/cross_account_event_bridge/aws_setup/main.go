package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
)

type Message struct {
	Body string
}

func main() {
	accountID := os.Args[1]
	queue := os.Args[2]
	profile := os.Args[3]
	region := os.Args[4]

	log.Println("Receiving events from:", queue)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	cli := sqs.NewFromConfig(cfg)

	queueURL := "https://sqs." + region + ".amazonaws.com/" + accountID + "/" + queue
	msgChan := make(chan Message)
	go Consume(cli, msgChan, queueURL)

	for msg := range msgChan {
		log.Println("Received message:", msg.Body)
	}
}

func Consume(cli *sqs.Client, msgChan chan Message, queueUrl string) {
	for {
		out, err := cli.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(queueUrl),
		})

		if err != nil {
			log.Fatalln("error in receiving msg", err)
			close(msgChan)
		}

		if out == nil || len(out.Messages) == 0 {
			continue
		}
		for _, msg := range out.Messages {
			msgChan <- Message{Body: *msg.Body}
		}
	}
}
