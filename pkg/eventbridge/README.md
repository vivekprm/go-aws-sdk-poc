https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-what-is.html

EventBridge is a serverless service that uses events to connect application components together, making it easier for you to build scalable event-driven applications. Event-driven architecture is a style of building loosely-coupled software systems that work together by emitting and responding to events. Event-driven architecture can help you boost agility and build reliable, scalable applications.

EventBridge provides simple and consistent ways to ingest, filter, transform, and deliver events so you can build applications quickly.

EventBridge includes two ways to process and deliver events: *event buses* and *pipes*.

[Event buses](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-event-bus.html) are routers that receive [events](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-events.html) and delivers them to zero or more targets. Use EventBridge to route events from sources such as home-grown applications, AWS services, and third-party software to consumer applications across your organization.

Event buses are well-suited for routing events from many sources to many targets, with optional transformation of events prior to delivery to a target.

[Pipes](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-pipes.html) EventBridge Pipes is intended for point-to-point integrations; each pipe receives events from a single source for processing and delivery to a single target. Pipes also include support for advanced transformations and enrichment of events prior to delivery to a target.

Pipes and event buses are often used together. A common use case is to create a pipe with an event bus as its target; the pipe sends events to the event bus, which then sends those events on to multiple targets. For example, you could create a pipe with a DynamoDB stream for a source, and an event bus as the target. The pipe receives events from the DynamoDB stream and sends them to the event bus, which then sends them on to multiple targets according to the rules you've specified on the event bus.

In addition, EventBridge provides [EventBridge Scheduler](https://docs.aws.amazon.com/eventbridge/latest/userguide/using-eventbridge-scheduler.html), a serverless scheduler that allows you to create, run, and manage tasks from one central, managed service. With EventBridge Scheduler, you can create schedules using cron and rate expressions for recurring patterns, or configure one-time invocations. You can set up flexible time windows for delivery, define retry limits, and set the maximum retention time for failed API invocations.

# Getting started
To get familiar with EventBridge rules and their capabilities, we'll use a **CloudFormation template** to set up an event bus rule and associated components, including an event source, event pattern, and target. Then we can explore how rules work to select the events you want.

The template creates a rule on the default event bus. This rule uses an event pattern to filter for events from a specific Amazon S3 bucket. The rule sends matching events to the specified target, an Amazon SNS topic. Every time an object is created in the bucket, the rule sends a notification to the topic, which then sends an email to your specified email address.

The deployed resources consist of:
- An Amazon S3 bucket with EventBridge notifications enabled to act as the event source.
- An Amazon SNS topic and email subscription as the target for notifications.
- An execution role that grants EventBridge the necessary permissions to publish to the Amazon SNS topic.
- The rule itself, which:
  - Defines an event pattern that matches only Object Created events from the specific Amazon S3 bucket.
  - Specifies the Amazon SNS topic as a target to which EventBridge delivers matching events.

For specific technical details of the template, see [Template details](https://docs.aws.amazon.com/eventbridge/latest/userguide/event-bus-rule-get-started.html#event-bus-rule-get-started-template-details).

## Creating the rule using CloudFormation
To create the rule and its associated resources, we'll create a CloudFormation template and use it to create a stack containing a sample rule, complete with source and target.

**Important**
You will be billed for the Amazon resources used if you create a stack from this template.

### Creating the template
First, create the CloudFormation template.
- In the [Template section](https://docs.aws.amazon.com/eventbridge/latest/userguide/event-bus-rule-get-started.html#event-bus-rule-get-started-template), click the copy icon on the JSON or YAML tab to copy the template contents.
- Paste the template contents into a new file.
- Save the file locally.

### Creating the stack
Next, use the template you've saved to provision a CloudFormation stack.

#### Create the stack using CloudFormation (console)
- Open the CloudFormation console at https://console.aws.amazon.com/cloudformation/.
- On the Stacks page, from the **Create stack** menu, choose with new resources (standard)
- Specify the template:
  - Under Prerequisite, choose Choose an existing template.
  - Under Specify template, choose Upload a template file.
  - Choose Choose file, navigate to the template file, and choose it.
  - Choose Next.
- Specify the stack details:
  - Enter a stack name.
  - For parameters, accept the default values for BucketName, SNSTopicDisplayName, SNSTopicName, and RuleName, or enter your own.
  - For EmailAddress, enter a valid email address where you want to receive notifications.
  - Choose Next.
- Configure the stack options:
  - Under Stack failure options, choose Delete all newly created resources.
  - Accept all other default values.
  - Under Capabilities, check the box to acknowledge that CloudFormation might create IAM resources in your account.
  - Choose Next.
- Review the stack details and choose Submit.

### Create the stack using CloudFormation (AWS CLI)
You can also use the AWS CLI to create the stack.

- Use the [create-stack](https://docs.aws.amazon.com/cli/latest/reference/cloudformation/create-stack.html) command.
  - Accept the default template parameter values, specifying the stack name and your email address. Use the template-body parameter to pass the template contents, or template-url to specify a URL location.
```sh
aws cloudformation create-stack \
  --stack-name eventbridge-rule-tutorial \
  --template-body template-contents \
  --parameters ParameterKey=EmailAddress,ParameterValue=your.email@example.com \
  --capabilities CAPABILITY_IAM
```
  -  Override the default value(s) of one or more template parameters. For example:

```sh
aws cloudformation create-stack \
  --stack-name eventbridge-rule-tutorial \
  ----template-body template-contents \
  --parameters \
    ParameterKey=EmailAddress,ParameterValue=your.email@example.com \
    ParameterKey=BucketName,ParameterValue=my-custom-bucket-name \
    ParameterKey=RuleName,ParameterValue=my-custom-rule-name \
  --capabilities CAPABILITY_IAM
```

CloudFormation creates the stack. Once the stack creation is complete, the stack resources are ready to use. You can use the Resources tab on the stack detail page to view the resources that were provisioned in your account.

After the stack is created, you will receive a subscription confirmation email at the address you provided. You must confirm this subscription to receive notifications.

# Exploring rule capabilities
Once the rule has been created, you can use the EventBridge console to observe rule operation and test event delivery.

- Open the EventBridge console at https://console.aws.amazon.com/events/home?#/rules.
- Choose the rule you created.
  - On the rule detail page, the Rule details section displays information about the rule, including its event pattern and targets.

# Examining the event pattern
Before we test the rule operation, let's examine the event pattern we've specified to control which events are sent to the target. The rule will only send events that match the pattern criteria to the target. In this case, we only want the event that Amazon S3 generates when an object is created in our specific bucket.

- On the rule detail page, under Event pattern, you can see the event pattern selects only events where:
  - The source is the Amazon S3 service (aws.s3)
  - The detail-type is Object Created
  - The bucket name matches the name of the bucket we created

```json
{
  "source": ["aws.s3"],
  "detail-type": ["Object Created"],
  "detail": {
    "bucket": {
      "name": ["eventbridge-rule-example-source"]
    }
  }
}
```

# Sending events through the rule
Next, we'll generate events in the event source to test that the rule matching and delivery is operating correctly. To do this, we'll upload an object to the S3 bucket we specified as the event source.

- Open the Amazon S3 console at https://console.aws.amazon.com/s3/.
- In the Buckets list, choose the bucket you created with the template (default name: eventbridge-rule-example-source).
- Choose Upload.
- Upload a test file to generate an ```Object Created``` event:
  - Choose Add files and select a file from your computer.
  - Choose Upload.
- Wait a few moments for the event to be processed by EventBridge and for the notification to be sent.
- Check your email for a notification about the object creation event. The email will contain details about the S3 event, including the bucket name and the object key.

# Viewing rule metrics
You can view metrics for your rule to confirm that events are being processed correctly.

- In the EventBridge console, choose your rule.
- Choose the Metrics tab.
- You can view metrics such as:
  - Invocations: the number of times the rule was triggered.
  - TriggeredRules: the number of rules that were triggered by matching events.

# CloudFormation template details
This template creates resources and grants permissions in your account.

## Resources
The CloudFormation template for this tutorial will create the following resources in your account:

- AWS::S3::Bucket: An Amazon S3 bucket that acts as the event source for the rule, with EventBridge notifications enabled.
- AWS::SNS::Topic: An Amazon SNS topic that acts as the target for the events matched by the rule.
- AWS::SNS::Subscription: An email subscription to the SNS topic.
- AWS::IAM::Role: An IAM execution role granting permissions to the EventBridge service in your account.
- AWS::Events::Rule: The rule connecting the Amazon S3 bucket events to the Amazon SNS topic.

## Permissions
The template includes an AWS::IAM::Role resource that represents an execution role. This role grants the EventBridge service (events.amazonaws.com) the following permissions in your account.

The following permissions are granted through the managed policy AmazonSNSFullAccess:

Full access to Amazon SNS resources and operations

# CloudFormation template
Save the following YAML code as a separate file to use as the CloudFormation template for this tutorial.

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Description: '[AWSDocs] EventBridge: event-bus-rule-get-started'

Parameters:
  BucketName:
    Type: String
    Description: Name of the S3 bucket
    Default: eventbridge-rule-example-source

  SNSTopicDisplayName:
    Type: String
    Description: Display name for the SNS topic
    Default: eventbridge-rule-example-target

  SNSTopicName:
    Type: String
    Description: Name for the SNS topic
    Default: eventbridge-rule-example-target

  RuleName:
    Type: String
    Description: Name for the EventBridge rule
    Default: eventbridge-rule-example

  EmailAddress:
    Type: String
    Description: Email address to receive notifications
    AllowedPattern: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+[a-zA-Z0-9-]*(\\.[a-zA-Z0-9-]+)*$'

Resources:
  # S3 Bucket with notifications enabled
  S3Bucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Ref BucketName
      NotificationConfiguration:
        EventBridgeConfiguration:
          EventBridgeEnabled: true

  # SNS Topic for email notifications
  SNSTopic:
    Type: AWS::SNS::Topic
    Properties:
      DisplayName: !Ref SNSTopicDisplayName
      TopicName: !Ref SNSTopicName

  # SNS Subscription for email
  SNSSubscription:
    Type: AWS::SNS::Subscription
    Properties:
      Protocol: email
      Endpoint: !Ref EmailAddress
      TopicArn: !Ref SNSTopic

  # EventBridge Rule to match S3 object creation events and send them to the SNS topic
  EventBridgeRule:
    Type: AWS::Events::Rule
    Properties:
      Name: !Ref RuleName
      Description: "Rule to detect S3 object creation and send email notification"
      EventPattern:
        source:
          - aws.s3
        detail-type:
          - "Object Created"
        detail:
          bucket:
            name:
              - !Ref BucketName
      State: ENABLED
      Targets:
        - Id: SendToSNS
          Arn: !Ref SNSTopic
          RoleArn: !GetAtt EventBridgeRole.Arn

  # IAM Role for EventBridge to publish to SNS
  EventBridgeRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17		 	 	 '
        Statement:
          - Effect: Allow
            Principal:
              Service: events.amazonaws.com
            Action: sts:AssumeRole
      ManagedPolicyArns:
        - arn:aws:iam::aws:policy/AmazonSNSFullAccess

Outputs:
  BucketName:
    Description: Name of the S3 bucket
    Value: !Ref S3Bucket
  SNSTopicARN:
    Description: ARN of the SNS topic
    Value: !Ref SNSTopic
  EmailSubscription:
    Description: Email address for notifications
    Value: !Ref EmailAddress
```

# Updating a default event bus using AWS CloudFormation in EventBridge
CloudFormation enables you to configure and manage your AWS resources across accounts and regions in a centralized and repeatable manner by treating infrastructure as code. CloudFormation does this by letting you create *templates*, which define the resources you want to provision and manage.

Because EventBridge provisions the default event bus into your account automatically, you cannot create it using a CloudFormation template, as you normally would for any resource you wanted to include in a CloudFormation stack. To include the default event bus in a CloudFormation stack, you must first import it into a stack. Once you have imported the default event bus into a stack, you can then update the event bus properties as desired.

To import an existing resource into a new or existing CloudFormation stack, you need the following information:

- A unique identifier for the resource to import.
  - For default event buses, the identifier is Name and then identifier value is default.
- A template that accurately describes the current properties of the existing resource.
  - The template snippet below contains an ```AWS::Events::EventBus``` resource that describes the current properties of a default event bus. In this example, the event bus has been configured to use a customer managed key and DLQ for encryption at rest.

Also, the ```AWS::Events::EventBus``` resource that describes the default event bus you want to import should include a ```DeletionPolicy``` property set to ```Retain```.

```json
{
    "AWSTemplateFormatVersion": "2010-09-09",
    "Description": "Default event bus import example",
    "Resources": {
        "defaultEventBus": {
            "Type" : "AWS::Events::EventBus",
            "DeletionPolicy": "Retain",
            "Properties" : {
                "Name" : "default",
                "KmsKeyIdentifier" : "KmsKeyArn",
                "DeadLetterConfig" : {
                    "Arn" : "DLQ_ARN"
                }
            }
        }
    }
}
```

For more information, see [Bringing existing resources into CloudFormation management](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/resource-import.html) in the CloudFormation User Guide.

