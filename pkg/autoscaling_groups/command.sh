# https://awscli.amazonaws.com/v2/documentation/api/latest/reference/autoscaling/describe-auto-scaling-groups.html

# Creates an Auto Scaling group with the specified name and attributes.
# Every Auto Scaling group has three size properties (DesiredCapacity , MaxSize , and MinSize ).
aws autoscaling create-auto-scaling-group \
    --profile AWS_Personal \
    --auto-scaling-group-name my-asg \
    --launch-template LaunchTemplateName=my-template-for-auto-scaling \
    --min-size 1 \
    --max-size 5 \
    --desired-capacity 1 \
    --vpc-zone-identifier "subnet-0163bd2b8c13cb139,subnet-0747a2a8e6aa1b256,subnet-014251eee77cce8d1"

# Use the latest version of launch template
aws autoscaling update-auto-scaling-group \
    --profile AWS_Personal \
    --auto-scaling-group-name my-asg \
    --desired-capacity 2 \
    --launch-template LaunchTemplateName=my-template-for-auto-scaling,Version='2'

# Delete autoscaling group
aws autoscaling delete-auto-scaling-group --profile AWS_Personal --auto-scaling-group-name my-asg --force-delete

# List autoscaling groups
aws autoscaling describe-auto-scaling-groups --profile AWS_Keeper --query 'AutoScalingGroups[*].[AutoScalingGroupName, LaunchConfigurationName]' --output table

# Get details of autoscaling group
aws autoscaling describe-auto-scaling-groups --auto-scaling-group-names="iheggTWJE" --profile AWS_Keeper --output json