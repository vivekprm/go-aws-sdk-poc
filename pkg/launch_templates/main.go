package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/joho/godotenv"
)

type templateData struct {
	templateName string
	imageId      string
	version      string
}

// https://docs.aws.amazon.com/autoscaling/ec2/userguide/examples-launch-templates-aws-cli.html#describe-launch-template-aws-cli
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v\n", err)
	}
	// Using the SDK's default configuration, load additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	// profile := "AWS-Volterra-developer"
	profile := os.Getenv("AWS_PROFILE")
	pprofile := os.Getenv("AWS_PROFILE_PERSONAL")

	// Loading default
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// Loading a profile
	cli := getClient(profile)
	pcli := getClient(pprofile)

	_, err = createTemplate(pcli, &templateData{
		templateName: "my-template-for-auto-scaling",
		imageId:      "ami-0e449927258d45bc4",
		version:      "version1",
	})
	if err != nil {
		log.Fatalf("failed to create launch template, %v", err)
	}
	listTemplates(cli)
	getTemplate(cli, "lt-01b55c3ecb654e488")
	// deleteTemplate(pcli, id)
}

func getClient(profile string) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)
}

func listTemplates(cli *ec2.Client) {
	// Build the request with its input parameters
	resp, err := cli.DescribeLaunchTemplates(context.TODO(), &ec2.DescribeLaunchTemplatesInput{
		MaxResults: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("failed to list templates, %v", err)
	}

	for i, tmpl := range resp.LaunchTemplates {
		fmt.Printf("Template %d\n", i)
		fmt.Println("--------------------------------")
		fmt.Printf("ID: %s\n", *tmpl.LaunchTemplateId)
		fmt.Printf("Name: %s\n", *tmpl.LaunchTemplateName)
		fmt.Printf("Created By: %s\n", *tmpl.CreatedBy)
		fmt.Printf("Created On: %s\n\n", *tmpl.CreateTime)
	}
}

func getTemplate(cli *ec2.Client, templateId string) {
	out, err := cli.DescribeLaunchTemplates(context.TODO(), &ec2.DescribeLaunchTemplatesInput{
		LaunchTemplateIds: []string{templateId},
	})
	if err != nil {
		log.Fatalf("failed to get template data, %v", err)
	}
	fmt.Println("Template detail")
	fmt.Println("--------------------------------")
	fmt.Printf("ID: %s\n", *out.LaunchTemplates[0].LaunchTemplateId)
	fmt.Printf("Name: %s\n", *out.LaunchTemplates[0].LaunchTemplateName)
	fmt.Printf("Created By: %s\n", *out.LaunchTemplates[0].CreatedBy)
	fmt.Printf("Created On: %s\n\n", *out.LaunchTemplates[0].CreateTime)
}

func createTemplate(cli *ec2.Client, templateData *templateData) (string, error) {
	out, err := cli.CreateLaunchTemplate(context.TODO(), &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: aws.String(templateData.templateName),
		LaunchTemplateData: &types.RequestLaunchTemplateData{
			ImageId:      aws.String(templateData.imageId),
			InstanceType: types.InstanceTypeT2Micro,
		},
		VersionDescription: aws.String(templateData.version),
	})
	if err != nil {
		return "", err
	}
	fmt.Printf("Launch template created successfully, id: %s\n", *out.LaunchTemplate.LaunchTemplateId)
	return *out.LaunchTemplate.LaunchTemplateId, nil
}

func deleteTemplate(cli *ec2.Client, id string) {
	resp, err := cli.DeleteLaunchTemplate(context.TODO(), &ec2.DeleteLaunchTemplateInput{
		LaunchTemplateId: aws.String(id),
	})
	if err != nil {
		log.Fatalf("failed to delete launch template, %v", err)
	}
	fmt.Printf("Deletion successful: %v", *resp.LaunchTemplate.LaunchTemplateId)
}
