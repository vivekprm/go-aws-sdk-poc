package receiver

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Message struct {
	Body string
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
