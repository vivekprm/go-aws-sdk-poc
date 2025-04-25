# https://docs.aws.amazon.com/cli/v1/userguide/cli-usage-output-format.html
# https://docs.aws.amazon.com/autoscaling/ec2/userguide/examples-launch-templates-aws-cli.html
aws ec2 describe-launch-templates --profile AWS_Keeper --query 'LaunchTemplates[*].[LaunchTemplateId, LaunchTemplateName]' --max-results 15 --output text
aws ec2 describe-launch-templates --launch-template-id lt-01b55c3ecb654e488 --profile AWS_Keeper

# creates a new launch template version based on version 1 of the launch template and specifies a different AMI ID.
aws ec2 create-launch-template --profile AWS_Personal --launch-template-name my-template-for-auto-scaling --version-description version1 --launch-template-data '{"ImageId":"ami-0e449927258d45bc4","InstanceType":"t2.micro"}'

# creates a new launch template version based on version 1 of the launch template and specifies a different AMI ID.
aws ec2 create-launch-template-version --profile AWS_Personal --launch-template-name my-template-for-auto-scaling --version-description version2 --source-version 1 --launch-template-data "ImageId=ami-0b86aaed8ef90e45f"

# deletes the specified launch template version.
aws ec2 delete-launch-template-versions --profile AWS_Personal --launch-template-name my-template-for-auto-scaling --versions 2

# Deleting a launch template deletes all of its versions.
aws ec2 delete-launch-template --profile AWS_Personal --launch-template-name my-template-for-auto-scaling