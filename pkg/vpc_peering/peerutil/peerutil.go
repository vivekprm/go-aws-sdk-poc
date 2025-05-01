package peerutil

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
)

func CreatePeering(cli *ec2.Client, vpc1, vpc2 *string, peerRegion string) *string {
	out, err := cli.CreateVpcPeeringConnection(context.Background(), &ec2.CreateVpcPeeringConnectionInput{
		VpcId:      vpc1,
		PeerRegion: aws.String(peerRegion),
		PeerVpcId:  vpc2,
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeVpcPeeringConnection,
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
		log.Fatalf("error in creating vpc peering: %v\n", err)
	}
	log.Printf("vpc peering created successfully: %s\n", *out.VpcPeeringConnection.VpcPeeringConnectionId)

	time.Sleep(5 * time.Second)
	acceptPeeringConnection(cli, out.VpcPeeringConnection.VpcPeeringConnectionId)

	return out.VpcPeeringConnection.VpcPeeringConnectionId
}

func acceptPeeringConnection(cli *ec2.Client, connectionID *string) {
	out, err := cli.AcceptVpcPeeringConnection(context.Background(), &ec2.AcceptVpcPeeringConnectionInput{
		VpcPeeringConnectionId: connectionID,
	})
	if err != nil {
		log.Fatalf("Error in accepting peering connection: %v\n", err)
	}
	log.Printf("Accepting peering connection successful: %v\n", out.VpcPeeringConnection.AccepterVpcInfo)
}

func CreatePeeringRoute(cli *ec2.Client, rtbID, vpcCIDR, peeringConnection string) {
	out, err := cli.CreateRoute(context.Background(), &ec2.CreateRouteInput{
		RouteTableId:           aws.String(rtbID),
		DestinationCidrBlock:   aws.String(vpcCIDR),
		VpcPeeringConnectionId: aws.String(peeringConnection),
	})
	if err != nil {
		log.Fatalf("Creation of route to vpc %s failed, rtbid %s: %v\n", vpcCIDR, rtbID, err)
	}
	fmt.Printf("Creation of route to vpc %s successful, rtbid %s: %v\n", vpcCIDR, rtbID, out)
}
