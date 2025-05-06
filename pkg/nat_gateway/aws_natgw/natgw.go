package awsnatgw

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func CreateNatGateway(cli *ec2.Client, subnetID string, connectivity types.ConnectivityType) *string {
	out, err := cli.CreateNatGateway(context.Background(), &ec2.CreateNatGatewayInput{
		SubnetId:         aws.String(subnetID),
		ConnectivityType: connectivity,
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeNatgateway,
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
		log.Fatalf("Error in creating nat gateway: %v\n", err)
	}
	log.Printf("Creation of nat gateway successful: %s\n", *out.NatGateway.NatGatewayId)
	return out.NatGateway.NatGatewayId
}
