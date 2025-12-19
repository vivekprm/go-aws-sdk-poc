package main

import (
	"log"

	"github.com/blushft/go-diagrams/diagram"
	"github.com/blushft/go-diagrams/nodes/aws"
	"github.com/blushft/go-diagrams/nodes/generic"
)

func main() {
	d, err := diagram.New(diagram.Filename("approach1"), diagram.Label("Approach 1"), diagram.Direction("LR"))
	if err != nil {
		log.Fatal(err)
	}

	cloudtrail := aws.Management.Cloudtrail(diagram.NodeLabel("CloudTrail"))
	s3 := aws.Storage.SimpleStorageServiceS3(diagram.NodeLabel("S3"))
	sqs := aws.Integration.SimpleQueueServiceSqs(diagram.NodeLabel("SQS"))
	eventbridge := aws.Integration.Eventbridge(diagram.NodeLabel("EventBridge"))
	lambda := aws.Compute.Lambda(diagram.NodeLabel("Lambda"))
	serv := generic.Compute.Rack(diagram.NodeLabel("GC Service"))

	dc := diagram.NewGroup("AWS")

	d.Connect(cloudtrail, s3, diagram.Forward()).Group(dc)
	d.Connect(s3, eventbridge, diagram.Forward()).Group(dc)
	d.Connect(eventbridge, lambda, diagram.Forward()).Group(dc)
	d.Connect(lambda, sqs, diagram.Forward()).Group(dc)
	d.Connect(sqs, serv, diagram.Forward()).Group(dc)

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
}
