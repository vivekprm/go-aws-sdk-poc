package awsvpc

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type VpcInfo struct {
	VpcID           string `json:"vpcID"`
	SecurityGroupID string `json:"securityGroupID"`
	SubnetID        string `json:"subnetID"`
	RouteTableID    string `json:"routeTableID"`
}

func CreateVpc(cli *ec2.Client, vpcCIDR string) *ec2.CreateVpcOutput {
	ctx := context.Background()
	resp, err := cli.CreateVpc(ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String(vpcCIDR),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeVpc,
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
		log.Fatal(err)
	}
	log.Printf("Creation of VPC successful: %s\n", *resp.Vpc.VpcId)
	return resp
}

func CreateSubnet(cli *ec2.Client, vpcID *string, az, subnetCIDR string) *string {
	resp1, err := cli.CreateSubnet(context.TODO(), &ec2.CreateSubnetInput{
		VpcId:            vpcID,
		AvailabilityZone: aws.String(az),
		CidrBlock:        aws.String(subnetCIDR),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSubnet,
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
		log.Fatalf("Creation of subnet1 failed: %v\n", err)
	}
	log.Printf("Creation of subnet1 successful: %v\n", *resp1.Subnet.SubnetId)
	return resp1.Subnet.SubnetId
}

func CreateRouteTable(cli *ec2.Client, subnetID, vpcID *string) *string {
	ctx := context.Background()
	resp, err := cli.CreateRouteTable(ctx, &ec2.CreateRouteTableInput{
		VpcId: vpcID,
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeRouteTable,
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
		log.Fatalf("Error in route table creation: %v\n", err)
	}
	log.Printf("Route table created successfully: %s\n", *resp.RouteTable.RouteTableId)
	out, err := cli.AssociateRouteTable(ctx, &ec2.AssociateRouteTableInput{
		RouteTableId: resp.RouteTable.RouteTableId,
		SubnetId:     subnetID,
	})
	if err != nil {
		log.Fatalf("Error in route table association with subnet: %v\n", err)
	}
	log.Printf("Route table associated with subnet successfully: %s\n", *out.AssociationId)
	return resp.RouteTable.RouteTableId
	// CreateRoute(ctx, cli, resp.RouteTable.RouteTableId, "0.0.0.0/0")
}

func CreateInstance(cli *ec2.Client, subnet1Id *string, sgID, amiID string) *string {
	ctx := context.Background()
	resp1, err := cli.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:          aws.String(amiID),
		InstanceType:     types.InstanceTypeT2Micro,
		MinCount:         aws.Int32(1),
		MaxCount:         aws.Int32(1),
		SubnetId:         aws.String(*subnet1Id),
		SecurityGroupIds: []string{sgID},
		KeyName:          aws.String("viveknitj06"),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
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
		log.Fatalf("Unable to create instance %v\n", err)
	}
	log.Printf("Instance created successfully: %s", *resp1.Instances[0].InstanceId)
	return resp1.Instances[0].InstanceId
}

func GetVPCData() VpcInfo {
	bytes, err := os.ReadFile("vpcdata.txt")
	if err != nil {
		log.Fatalf("Unable to read file: %v\n", err)
	}

	var vpcInfo VpcInfo
	json.Unmarshal(bytes, &vpcInfo)
	return vpcInfo
}

func CreateEc2InstanceConnect(cli *ec2.Client, subnetID *string, sgIDs []string) *string {
	resp, err := cli.CreateInstanceConnectEndpoint(context.Background(), &ec2.CreateInstanceConnectEndpointInput{
		SubnetId:         subnetID,
		SecurityGroupIds: sgIDs,
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstanceConnectEndpoint,
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
		log.Fatalf("Unable to create instance connect endpoint %v\n", err)
	}
	log.Printf("Instance connect endpoint created successfully: %s", *resp.InstanceConnectEndpoint.InstanceConnectEndpointId)
	return resp.InstanceConnectEndpoint.InstanceConnectEndpointId
}
