package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2_types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v\n", err)
	}
	profile := os.Getenv("AWS_PROFILE_PERSONAL_US")
	cli := getClient(profile)

	subnets := []string{
		"subnet-010fb94965225ba14",
	}
	ctx := context.Background()
	out, err := cli.CreateTransitGatewayVpcAttachment(ctx, &ec2.CreateTransitGatewayVpcAttachmentInput{
		SubnetIds:        subnets,
		TransitGatewayId: aws.String("tgw-0941d03f2fff144db"),
		VpcId:            aws.String("vpc-075c59840495e1469"),
		TagSpecifications: []ec2_types.TagSpecification{
			{
				ResourceType: ec2_types.ResourceTypeTransitGatewayAttachment,
				Tags: []ec2_types.Tag{
					{
						Key:   aws.String("CreatedBy"),
						Value: aws.String("vivek-poc"),
					},
				},
			},
		},
		Options: &ec2_types.CreateTransitGatewayVpcAttachmentRequestOptions{
			ApplianceModeSupport:            ec2_types.ApplianceModeSupportValueDisable,
			DnsSupport:                      ec2_types.DnsSupportValueEnable,
			Ipv6Support:                     ec2_types.Ipv6SupportValueDisable,
			SecurityGroupReferencingSupport: ec2_types.SecurityGroupReferencingSupportValueDisable,
		},
	})

	if err != nil {
		log.Fatalf("Error creating transit gateway attachment %v\n", err)
	}
	log.Printf("Created transit gateway attachment: %v\n", out)

	ok, err := pollTransitGatewayAttachmentStatus(ctx, "tgw-0941d03f2fff144db", *out.TransitGatewayVpcAttachment.TransitGatewayAttachmentId, cli)

	if err != nil {
		log.Fatalf("Error while waiting for the attachment to become available: %v", err)
	}

	if !ok {
		err := fmt.Errorf("TGW attachment status for %s is not found", *out.TransitGatewayVpcAttachment.TransitGatewayAttachmentId)
		log.Fatalf("Error while waiting for the attachment to become available: %v", err)
	}
	log.Printf("Transit gateway attachment %s is now available", *out.TransitGatewayVpcAttachment.TransitGatewayAttachmentId)
}

func pollTransitGatewayAttachmentStatus(ctx context.Context, tgwID, tgwAttachID string, cli *ec2.Client) (bool, error) {
	statusCheckInterval := time.NewTicker(10 * time.Second)
	defer statusCheckInterval.Stop()

	statusChecktimeOut := time.After(120 * time.Second)

	for {
		select {
		case <-statusChecktimeOut:
			err := fmt.Errorf("timeout while waiting for the attachment %s to become available", tgwAttachID)

			return false, err
		case <-statusCheckInterval.C:
			o := []ec2_types.TransitGatewayVpcAttachment{}

			var nextToken *string
			for {
				t, err := cli.DescribeTransitGatewayVpcAttachments(ctx, &ec2.DescribeTransitGatewayVpcAttachmentsInput{
					Filters: []ec2_types.Filter{
						{
							Name:   aws.String("transit-gateway-id"),
							Values: []string{tgwID},
						},
					},
					MaxResults: aws.Int32(100),
					NextToken:  nextToken,
				})
				if err != nil {
					return false, err
				}

				log.Printf("DescribeTransitGatewayVpcAttachments output: %v", t.TransitGatewayVpcAttachments[0].State)

				if t == nil || t.TransitGatewayVpcAttachments == nil {
					break
				}

				o = append(o, t.TransitGatewayVpcAttachments...)
				if t.NextToken == nil {
					break
				}

				nextToken = t.NextToken
			}

			var attachExists bool

			for _, a := range o {
				log.Printf("checking tgw attachment: %v", a)

				if a.TransitGatewayAttachmentId == nil {
					continue
				}

				if *a.TransitGatewayAttachmentId != tgwAttachID {
					continue
				}

				attachExists = true
				if a.State == ec2_types.TransitGatewayAttachmentStateAvailable {
					return true, nil
				}
			}

			if !attachExists {
				return false, nil
			}
		}
	}
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
}
