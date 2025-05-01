package awssg

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func CreateSecurityGroup(cli *ec2.Client, name, vpcID *string) *string {
	ctx := context.Background()
	resp, err := cli.CreateSecurityGroup(ctx, &ec2.CreateSecurityGroupInput{
		GroupName:   name,
		VpcId:       vpcID,
		Description: aws.String("custom security group"),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSecurityGroup,
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
		log.Fatalf("Error in security group creation: %v\n", err)
	}
	log.Printf("Security group created successfully: %s\n", *resp.GroupId)

	out, err := cli.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: resp.GroupId,
		IpPermissions: []types.IpPermission{
			{
				FromPort:   aws.Int32(22),
				ToPort:     aws.Int32(22),
				IpProtocol: aws.String("tcp"),
				IpRanges: []types.IpRange{
					{
						CidrIp: aws.String("0.0.0.0/0"),
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Error in modifying security group rules: %v\n", err)
	}
	log.Printf("Security group rules modification successful: %v\n", out.ResultMetadata)
	return resp.GroupId
}

func AllowPort(cli *ec2.Client, protocolNum, sgID, cidrIp string, fromPort, toPort int32) {
	out, err := cli.AuthorizeSecurityGroupIngress(context.Background(), &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(sgID),
		IpPermissions: []types.IpPermission{{
			IpProtocol: aws.String(protocolNum), // ICMP
			IpRanges: []types.IpRange{{
				CidrIp: aws.String(cidrIp),
			}},
			FromPort: aws.Int32(fromPort),
			ToPort:   aws.Int32(toPort),
		}},
	})
	if err != nil {
		log.Fatalf("error in modifying security group %s, %v\n", sgID, err)
	}
	log.Printf("Modification of security group %s successful: %v\n", sgID, out)
}
