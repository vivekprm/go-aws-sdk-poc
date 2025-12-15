https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-user-guide.html
https://docs.aws.amazon.com/awscloudtrail/latest/userguide/how-cloudtrail-works.html
https://aws.amazon.com/cloudtrail/pricing/
https://aws.amazon.com/s3/pricing/

AWS CloudTrail is an AWS service that helps you enable operational and risk auditing, governance, and compliance of your AWS account. Actions taken by a user, role, or an AWS service are recorded as events in CloudTrail. Events include actions taken in the AWS Management Console, AWS Command Line Interface, and AWS SDKs and APIs.

CloudTrail provides three ways to record events:

# Event history
The Event history provides a viewable, searchable, downloadable, and immutable record of the past 90 days of management events in an AWS Region. You can search events by filtering on a single attribute. You automatically have access to the Event history when you create your account. For more information, see [Working with CloudTrail event history](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/view-cloudtrail-events.html).

There are no CloudTrail charges for viewing the Event history.

## Limitations of Event History
The following limitations apply to the event history.

- The Event history page on the CloudTrail console only shows management events. It does not show data events, Insights events, or network activity events.
- The event history is limited to the past 90 days of events. For an ongoing record of events in your AWS account, create an [event data store](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/query-event-data-store-cloudtrail.html) or a [trail](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-create-a-trail-using-the-console-first-time.html).
- When you download events from the Event history page on the CloudTrail console, you can download up to 200,000 events in a single file. If you reach the 200,000 event limit, the CloudTrail console will provide the option to download additional files.
- The event history doesn't provide organization level event aggregation. To record events across your organization, create an organization event data store or trail.
- An event history search is limited to a single AWS account, only returns events from a single AWS Region, and cannot query multiple attributes. You can only apply one attribute filter and a time range filter.

You can create a CloudTrail Lake event data store to query across multiple attributes and AWS Regions. You can also query across multiple AWS accounts in an AWS Organizations organization. In CloudTrail Lake, you can query multiple event types, including management events, data events, Insights events, AWS Config configuration items, Audit Manager evidence, and non-AWS events. CloudTrail Lake queries offer a deeper and more customizable view of events than simple key and value lookups on the Event history page, or by running LookupEvents. For more information, see [Working with AWS CloudTrail Lake](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-lake.html) and [Create an event data store for CloudTrail events](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/query-event-data-store-cloudtrail.html) with the console.
- You cannot exclude AWS KMS or Amazon RDS Data API events from event history; settings that you apply to a trail or event data store do not apply to event history.

# CloudTrail Lake
[AWS CloudTrail Lake](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-lake.html) is a managed data lake for capturing, storing, accessing, and analyzing user and API activity on AWS for audit and security purposes. CloudTrail Lake converts existing events in row-based JSON format to [Apache ORC](https://orc.apache.org/) format. 

ORC is a columnar storage format that is optimized for fast retrieval of data. Events are aggregated into event data stores, which are immutable collections of events based on criteria that you select by applying advanced event selectors. You can keep the event data in an event data store for up to 3,653 days (about 10 years) if you choose the One-year extendable retention pricing option, or up to 2,557 days (about 7 years) if you choose the Seven-year retention pricing option. You can create an event data store for a single AWS account or for multiple AWS accounts by using AWS Organizations. 
You can import any existing CloudTrail logs from your S3 buckets into an existing or new event data store. You can also visualize top CloudTrail event trends with [Lake dashboards](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/lake-dashboard.html). For more information, see [Working with AWS CloudTrail Lake](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-lake.html).

CloudTrail Lake event data stores and queries incur charges. When you create an event data store, you choose the [pricing option](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-lake-manage-costs.html#cloudtrail-lake-manage-costs-pricing-option) you want to use for the event data store. The pricing option determines the cost for ingesting and storing events, and the default and maximum retention period for the event data store. When you run queries in Lake, you pay based upon the amount of data scanned. For information about CloudTrail pricing and managing Lake costs, see [AWS CloudTrail Pricing](https://aws.amazon.com/cloudtrail/pricing/) and [Managing CloudTrail Lake costs](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-lake-manage-costs.html).

# Trails
Trails capture a record of AWS activities, delivering and storing these events in an Amazon S3 bucket, with optional delivery to [CloudWatch Logs](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/send-cloudtrail-events-to-cloudwatch-logs.html) and [Amazon EventBridge](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-aws-service-specific-topics.html#cloudtrail-aws-service-specific-topics-eventbridge). You can input these events into your security monitoring solutions. You can also use your own third-party solutions or solutions such as Amazon Athena to search and analyze your CloudTrail logs. You can create trails for a single AWS account or for multiple AWS accounts by using AWS Organizations. You can log [Insights events](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/logging-insights-events-with-cloudtrail.html) to analyze your management events for anomalous behavior in API call rates and error rates. For more information, see [Creating a trail for your AWS account](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-create-and-update-a-trail.html).

You can deliver one copy of your ongoing management events to your S3 bucket at no charge from CloudTrail by creating a trail, [limits may apply](https://docs.aws.amazon.com/general/latest/gr/aws_service_limits.html), however, there are Amazon S3 storage charges. For more information about CloudTrail pricing, see [AWS CloudTrail Pricing](https://aws.amazon.com/cloudtrail/pricing/). For information about Amazon S3 pricing, see [Amazon S3 Pricing](https://aws.amazon.com/s3/pricing/).

When you create a trail, you enable ongoing delivery of events as log files to an Amazon S3 bucket that you specify. Creating a trail has many benefits, including:
- A record of events that extends past 90 days.
- The option to automatically monitor and alarm on specified events by sending log events to Amazon CloudWatch Logs.
- The option to query logs and analyze AWS service activity with Amazon Athena.

## Multi-Region trails
When you create a multi-Region trail, CloudTrail records events in all AWS Regions that are enabled in your AWS account and delivers the CloudTrail event log files to an S3 bucket that you specify. As a best practice, we recommend creating a multi-Region trail because it captures activity in all enabled Regions. All trails created using the CloudTrail console are multi-Region trails. You can convert a single-Region trail to a multi-Region trail by using the AWS CLI. For more information, see [Understanding multi-Region trails and opt-in Regions](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-multi-region-trails.html), [Creating a trail with the console](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-create-a-trail-using-the-console-first-time.html#creating-a-trail-in-the-console), and [Converting a single-Region trail to a multi-Region trail](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-create-and-update-a-trail-by-using-the-aws-cli-update-trail.html#cloudtrail-create-and-update-a-trail-by-using-the-aws-cli-examples-convert).

## Single-Region trails
When you create a single-Region trail, CloudTrail records the events in that Region only. It then delivers the CloudTrail event log files to an Amazon S3 bucket that you specify. You can only create a single-Region trail by using the AWS CLI. If you create additional single trails, you can have those trails deliver CloudTrail event log files to the same S3 bucket or to separate buckets. This is the default option when you create a trail using the AWS CLI or the CloudTrail API. For more information, see [Creating, updating, and managing trails with the AWS CLI](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-create-and-update-a-trail-by-using-the-aws-cli.html).

**Note**
For both types of trails, you can specify an Amazon S3 bucket from any Region.

If you have created an organization in AWS Organizations, you can create an *organization trail* that logs all events for all AWS accounts in that organization. **Organization trails can apply to all AWS Regions, or the current Region**. 

Organization trails must be created using the management account or delegated administrator account, and when specified as applying to an organization, are automatically applied to all member accounts in the organization. Member accounts can see the organization trail, but cannot modify or delete it. By default, member accounts do not have access to the log files for an organization trail in the Amazon S3 bucket.

By default, when you create a trail in the CloudTrail console, your event log files and digest files are encrypted with a KMS key. If you choose not to enable SSE-KMS encryption, your event log files and digest files are encrypted using Amazon S3 server-side encryption (SSE). You can store your log files in your bucket for as long as you want. You can also define Amazon S3 lifecycle rules to archive or delete log files automatically. If you want notifications about log file delivery and validation, you can set up Amazon SNS notifications.

CloudTrail publishes log files multiple times an hour, about every 5 minutes. These log files contain API calls from services in the account that support CloudTrail. For more information, see [CloudTrail supported services and integrations](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-aws-service-specific-topics.html).

**Note**
CloudTrail typically delivers logs within an average of about **5 minutes** of an API call. This time is not guaranteed. Review the [AWS CloudTrail Service Level Agreement](https://aws.amazon.com/cloudtrail/sla) for more information.

If you misconfigure your trail (for example, the S3 bucket is unreachable), CloudTrail will attempt to redeliver the log files to your S3 bucket for 30 days, and these attempted-to-deliver events will be subject to standard CloudTrail charges. **To avoid charges on a misconfigured trail, you need to delete the trail**.

CloudTrail captures actions made directly by the user or on behalf of the user by an AWS service. For example, an CloudFormation CreateStack call can result in additional API calls to Amazon EC2, Amazon RDS, Amazon EBS, or other services as required by the CloudFormation template. This behavior is normal and expected. You can identify if the action was taken by an AWS service with the **invokedby** field in the CloudTrail event.

The following table provides information about tasks you can perform on trails.

| Task                      |	Description                                                |
|---------------------------|--------------------------------------------------------------|
|Logging management events  | Configure your trails to log read-only, write-only, or all   | 
|                           | management events.                                           |
|---------------------------|--------------------------------------------------------------|
|Log data events            | You can use advanced event selectors to create fine-grained  |
|                           | selectors to log only those data events of interest. For     |
|                           | example, you can filter on the eventName field to include or |
|                           | exclude logging of specific API calls, which can help control|
|                           | costs. For more information, see Filtering data events by    |
|                           | using advanced event selectors.                              |
|---------------------------|--------------------------------------------------------------|

https://docs.aws.amazon.com/awscloudtrail/latest/userguide/logging-management-events-with-cloudtrail.html
https://docs.aws.amazon.com/awscloudtrail/latest/userguide/logging-data-events-with-cloudtrail.html

...

**You can deliver one copy of your ongoing management events to your S3 bucket at no charge from CloudTrail by creating a trail, however, there are Amazon S3 storage charges**. For more information about CloudTrail pricing, see [AWS CloudTrail Pricing](https://aws.amazon.com/cloudtrail/pricing/). For information about Amazon S3 pricing, see [Amazon S3 Pricing](https://aws.amazon.com/s3/pricing/).

# CloudTrail Insights events
AWS CloudTrail Insights help AWS users identify and respond to unusual activity associated with API call rates and API error rates by continuously analyzing CloudTrail management events. CloudTrail Insights analyzes your normal patterns of API call volume and API error rates, also called the baseline, and generates Insights events when the call volume or error rates are outside normal patterns. Insights events on API call rate are generated for write management APIs, and Insights events on API error rate are generated for both read and write management APIs.

By default, CloudTrail trails and event data stores don't log Insights events. You must configure your trail or event data store to log Insights events. For more information, see [Logging Insights events with the CloudTrail console](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/insights-events-enable.html) and [Logging Insights events with the AWS CLI](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/insights-events-CLI-enable.html).

Additional charges apply for Insights events. You will be charged separately if you enable Insights for both trails and event data stores. For more information, see [AWS CloudTrail Pricing](https://aws.amazon.com/cloudtrail/pricing/).

## Viewing Insights events for trails and event data stores
CloudTrail supports Insights events for both trails and event data stores, however, there are some differences in how you view and access Insights events.

### Viewing Insights events for trails
If you have Insights events enabled on a trail, and CloudTrail detects unusual activity, Insights events are logged to a different folder or prefix in the destination S3 bucket for your trail. You can also see the type of insight and the incident time period when you view Insights events on the CloudTrail console. For more information, see [Viewing Insights events for trails with the console](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/view-insights-events-console.html).

After you enable CloudTrail Insights for the first time on a trail, CloudTrail may take up to 36 hours to begin delivering Insights events after you enable Insights events on a trail, provided that unusual activity is detected during that time.

### Viewing Insights events for event data stores
To log Insights events in CloudTrail Lake, you need a destination event data store that logs Insights events and a source event data store that enables Insights and logs management events. For more information, see [Create an event data store for Insights events with the console](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/query-event-data-store-insights.html).

After you enable CloudTrail Insights for the first time on the source event data store, CloudTrail may take up to 7 days to begin delivering Insights events, provided that unusual activity is detected during that time.

If you have CloudTrail Insights enabled on a source event data store and CloudTrail detects unusual activity, CloudTrail delivers Insights events to your destination event data store. You can then query your destination event data store to get information about your Insights events and can optionally save the query results to an S3 bucket. For more information, see [Create or edit a query with the CloudTrail console](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/how-cloudtrail-works.html#:~:text=Create%20or%20edit%20a%20query%20with%20the%20CloudTrail%20console) and View sample queries with the CloudTrail console.

You can view the Insights events dashboard to visualize the Insights events in your destination event data store. For more information about Lake dashboards, see [CloudTrail Lake dashboards](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/lake-dashboard.html).

# CloudTrail channels
CloudTrail supports two types of channels:

**Channels for CloudTrail Lake integrations with event sources outside of AWS**
CloudTrail Lake uses channels to bring events from outside of AWS into CloudTrail Lake from external partners that work with CloudTrail, or from your own sources. When you create a channel, you choose one or more event data stores to store events that arrive from the channel source. You can change the destination event data stores for a channel as needed, as long as the destination event data stores are set to log activity events. When you create a channel for events from an external partner, you provide a channel ARN to the partner or source application. The resource policy attached to the channel allows the source to transmit events through the channel. For more information, see Create an integration with an event source outside of AWS and CreateChannel in the AWS CloudTrail API Reference.

**Service-linked channels**
AWS services can create a service-linked channel to receive CloudTrail events on your behalf. The AWS service creating the service-linked channel configures advanced event selectors for the channel and specifies whether the channel applies to all Regions, or the current Region.

You can use the CloudTrail console or AWS CLI to view information about any CloudTrail service-linked channels created by AWS services.