https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-concepts.html
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/monitor-with-cloudtrail.html

An event in CloudTrail is the record of an activity in an AWS account. This activity can be an action taken by an IAM identity, or service that is monitorable by CloudTrail. CloudTrail events provide a history of both API and non-API account activity made through the AWS Management Console, AWS SDKs, command line tools, and other AWS services.

CloudTrail log files aren't an ordered stack trace of the public API calls, so events don't appear in any specific order.

CloudTrail logs four types of events:
- [Management events](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-concepts.html#cloudtrail-concepts-management-events)
- [Data events](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-concepts.html#cloudtrail-concepts-data-events)
- [Network activity events](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-concepts.html#cloudtrail-concepts-network-events)
- [Insights events](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-concepts.html#cloudtrail-concepts-insights-events)

All event types use a CloudTrail JSON log format.

By default, trails and event data stores log management events, but not data or Insights events.

For information about how AWS services integrate with CloudTrail, see [AWS service topics for CloudTrail](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-aws-service-specific-topics.html#cloudtrail-aws-service-specific-topics-list).

# Management events
Management events provide information about management operations that are performed on resources in your AWS account. These are also known as control plane operations.

Example management events include:
- Configuring security (for example, AWS Identity and Access Management AttachRolePolicy API operations).
- Registering devices (for example, Amazon EC2 CreateDefaultVpc API operations).
- Configuring rules for routing data (for example, Amazon EC2 CreateSubnet API operations).
- Setting up logging (for example, AWS CloudTrail CreateTrail API operations).

Management events can also include non-API events that occur in your account. For example, when a user signs in to your account, CloudTrail logs the ConsoleLogin event. For more information, see [Non-API events captured by CloudTrail](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-non-api-events.html).

By default, CloudTrail trails and CloudTrail Lake event data stores log management events. For more information about logging management events, see [Logging management events](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/logging-management-events-with-cloudtrail.html).