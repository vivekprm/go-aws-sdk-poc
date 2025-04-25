Taken From: https://docs.aws.amazon.com/autoscaling/ec2/userguide/launch-templates.html

## Introduction
A launch template is similar to a [launch configuration](https://docs.aws.amazon.com/autoscaling/ec2/userguide/launch-configurations.html), in that it specifies instance configuration information. It includes:
- the ID of the Amazon Machine Image (AMI), 
- the instance type, 
- a key pair, 
- security groups, 
- and other parameters used to launch EC2 instances. 

However, defining a launch template instead of a launch configuration allows you to have **multiple versions of a launch template**.

With versioning of launch templates, you can create a subset of the full set of parameters. Then, you can reuse it to create other versions of the same launch template. 

For example, you can create a launch template that defines a base configuration without an AMI or user data script. 

After you create your launch template, you can create a new version and add the AMI and user data that has the latest version of your application for testing. This results in two versions of the launch template. 

Storing a base configuration helps you to maintain the required general configuration parameters. You can create a new version of your launch template from the base configuration whenever you want. You can also delete the versions used for testing your application when you no longer need them.

We recommend that you use launch templates to ensure that you're accessing the latest features and improvements. 

Not all Amazon EC2 Auto Scaling features are available when you use launch configurations. For example, **you cannot create an Auto Scaling group that launches both Spot and On-Demand Instances or that specifies multiple instance types**. 

You must use a launch template to configure these features. For more information, see [Auto Scaling groups with multiple instance types and purchase options](https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-mixed-instances-groups.html).

With launch templates, you can also use newer features of Amazon EC2. This includes:
- Systems Manager parameters (AMI ID), 
- the current generation of EBS Provisioned IOPS volumes (io2), 
- EBS volume tagging, 
- T2 Unlimited instances, 
- Capacity Reservations, 
- Capacity Blocks, 
- and Dedicated Hosts, 

to name a few.

When you create a launch template, all parameters are optional. However, if a launch template does not specify an AMI, you cannot add the AMI when you create your Auto Scaling group. If you specify an AMI but no instance type, you can add one or more instance types when you create your Auto Scaling group.

## Permissions To Work With Launch Templates
The procedures in this section assume that you already have the required permissions to create launch templates. For information about how an administrator grants you permissions, see [Control access to launch templates with IAM permissions](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/permissions-for-launch-templates.html) in the Amazon EC2 User Guide.

Note that if you do not have sufficient permissions to use and create resources specified in a launch template, you receive an error that you're not authorized to use the launch template when you try to specify it for an Auto Scaling group. For more information, see [Troubleshoot Amazon EC2 Auto Scaling: Launch templates](https://docs.aws.amazon.com/autoscaling/ec2/userguide/ts-as-launch-template.html).

For examples of IAM policies that let you call the ```CreateAutoScalingGroup```, ```UpdateAutoScalingGroup```, and ```RunInstances``` API operations with a launch template, see [Control Amazon EC2 launch template usage in Auto Scaling groups](https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-launch-template-permissions.html).

## API operations supported by launch templates
For a list of API operations supported by launch templates, see [Amazon EC2 actions](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/OperationList-query-ec2.html) in the [Amazon EC2 API Reference](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/).

# Creating a Launch Template Using CLI
```sh
aws ec2 create-launch-template --launch-template-name my-template-for-auto-scaling --version-description version1 \
  --launch-template-data '{"ImageId":"ami-04d5cc9b88example","InstanceType":"t2.micro"}'
```

Or use can pass launch template data in a file:
```sh
aws ec2 create-launch-template --launch-template-name my-template-for-auto-scaling --version-description version1 \
  --launch-template-data file://config.json
```

config.json
```json
{
    "LaunchTemplateName": "my-template-for-auto-scaling",
    "VersionDescription": "test description",
    "LaunchTemplateData": {
        "ImageId": "ami-04d5cc9b88example",
        "IamInstanceProfile": {
            "Name":"my-instance-profile"
        },
        "InstanceType": "t2.micro",
        "TagSpecifications": [
            {
                "ResourceType":"instance",
                "Tags": [
                    {"Key":"purpose","Value":"webserver"}
                ]
            }
        ],
        "NetworkInterfaces": [
            {
                "DeviceIndex":0,
                "AssociatePublicIpAddress":true,
                "Groups":["sg-903004f88example"],
                "DeleteOnTermination":true
            }
        ],
        "UserData": "IyEvYmluL2Jhc...",
        "BlockDeviceMappings": [
            {
                "DeviceName":"/dev/xvdcz",
                "Ebs":{
                    "VolumeSize":22,
                    "VolumeType":"gp2",
                    "DeleteOnTermination":true
                }
            }
        ],
        "SecurityGroupIds": [
            "sg-903004f88example"
        ], 
        "KeyName": "MyKeyPair",
        "Monitoring": {
            "Enabled": true
        },
        "Placement": {
            "Tenancy": "dedicated"
        },
        "CreditSpecification": {
            "CpuCredits": "unlimited"
        },
        "MetadataOptions": {
            "HttpTokens": "required",
            "HttpPutResponseHopLimit": 1,
            "HttpEndpoint": "enabled"
        }
    }
}
```

The above config example specifies:
- The name of the instance profile associated with the IAM role to pass to instances at launch.
- Adds a tag (for example, purpose=webserver) to instances at launch.
- Configures the launch template to assign public addresses to instances launched in a nondefault VPC
- Specifies a user data script as a base64-encoded string that configures instances at launch.
- Creates a launch template with a block device mapping: a 22-gigabyte EBS volume mapped to ```/dev/xvdcz```. The ```/dev/xvdcz``` volume uses the General Purpose SSD (gp2) volume type and is deleted when terminating the instance it is attached to.

Look at below links for more configuration options:
https://docs.aws.amazon.com/autoscaling/ec2/userguide/examples-launch-templates-aws-cli.html
