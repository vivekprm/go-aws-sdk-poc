# Nat Gateways
https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-gateway.html

A NAT gateway is a Network Address Translation (NAT) service. You can use a NAT gateway so that instances in a private subnet can connect to services outside your VPC but external services can't initiate a connection with those instances.

When you create a NAT gateway, you specify one of the following connectivity types:
- **Public** – (Default) Instances in private subnets can connect to the internet through a public NAT gateway, but the instances can't receive unsolicited inbound connections from the internet. You create a public NAT gateway in a public subnet and must associate an elastic IP address with the NAT gateway at creation. You route traffic from the NAT gateway to the internet gateway for the VPC. Alternatively, you can use a public NAT gateway to connect to other VPCs or your on-premises network. In this case, you route traffic from the NAT gateway through a transit gateway or a virtual private gateway.
- **Private** – Instances in private subnets can connect to other VPCs or your on-premises network through a private NAT gateway, but the instances can't receive unsolicited inbound connections from the other VPCs or the on-premises network. You can route traffic from the NAT gateway through a transit gateway or a virtual private gateway. You can't associate an elastic IP address with a private NAT gateway. You can attach an internet gateway to a VPC with a private NAT gateway, but **if you route traffic from the private NAT gateway to the internet gateway, the internet gateway drops the traffic**.

A NAT gateway is for use with IPv4 or IPv6 traffic (using [DNS64 and NAT64](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-nat64-dns64.html) ). Another option for enabling outbound-only internet communication over IPv6 is using an [egress-only internet gateway](https://docs.aws.amazon.com/vpc/latest/userguide/egress-only-internet-gateway.html).

Both private and public NAT gateways map the source private IPv4 address of the instances to the private IPv4 address of the NAT gateway, but in the case of a public NAT gateway, the internet gateway then maps the private IPv4 address of the public NAT gateway to the Elastic IP address associated with the NAT gateway. When sending response traffic to the instances, whether it's a public or private NAT gateway, the NAT gateway translates the address back to the original source IP address.

**Important**
Connections must always be initiated from within the VPC containing the NAT Gateway.

You can use either a public or private NAT gateway to route traffic to transit gateways and virtual private gateways.

If you use a private NAT gateway to connect to a transit gateway or virtual private gateway, traffic to the destination will come from the private IP address of the private NAT gateway.

If you use a public NAT gateway to connect to a transit gateway or virtual private gateway, traffic to the destination will come from the private IP address of the public NAT gateway. The public NAT gateway will only use its EIP as the source IP address when used in conjunction with an internet gateway in the same VPC.

NAT gateways support traffic with a maximum transmission unit (MTU) of 8500. For more information, see [NAT gateway basic](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-basics.html)s.

# [NAT gateway basic](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-basics.html)
Each NAT gateway is created in a specific Availability Zone and implemented with redundancy in that zone. There is a quota on the number of NAT gateways that you can create in each Availability Zone. For more information, see [Amazon VPC quotas](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html).

If you have resources in multiple Availability Zones and they share one NAT gateway, and if the NAT gateway’s Availability Zone is down, resources in the other Availability Zones lose internet access. To improve resiliency, create a NAT gateway in each Availability Zone, and configure your routing to ensure that resources use the NAT gateway in the same Availability Zone.

The following characteristics and rules apply to NAT gateways:
- A NAT gateway supports the following protocols: TCP, UDP, and ICMP.
- NAT gateways are supported for IPv4 or IPv6 traffic. For IPv6 traffic, NAT gateway performs NAT64. By using this in conjunction with DNS64 (available on Route 53 resolver), your IPv6 workloads in a subnet in Amazon VPC can communicate with IPv4 resources. These IPv4 services may be present in the same VPC (in a separate subnet) or a different VPC, on your on-premises environment or on the internet
- A NAT gateway supports 5 Gbps of bandwidth and automatically scales up to 100 Gbps. If you require more bandwidth, you can split your resources into multiple subnets and create a NAT gateway in each subnet.
- A NAT gateway can process one million packets per second and automatically scales up to ten million packets per second. Beyond this limit, a NAT gateway will drop packets. To prevent packet loss, split your resources into multiple subnets and create a separate NAT gateway for each subnet.
- Each IPv4 address can support up to 55,000 simultaneous connections to each unique destination. A unique destination is identified by a unique combination of destination IP address, the destination port, and protocol (TCP/UDP/ICMP). You can increase this limit by associating up to 8 IPv4 addresses to your NAT gateways (1 primary IPv4 address and 7 secondary IPv4 addresses). You are limited to associating 2 Elastic IP addresses to your public NAT gateway by default. You can increase this limit by requesting a quota adjustment. For more information, see [Elastic IP addresses](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html#vpc-limits-eips).
- You can pick the private IPv4 address to assign to the NAT gateway or have it automatically assigned from the IPv4 address range of the subnet. The assigned private IPv4 address persists until you delete the private NAT gateway. You can't detach the private IPv4 address and you can't attach additional private IPv4 addresses.
- You can't associate a security group with a NAT gateway. You can associate security groups with your instances to control inbound and outbound traffic.
- You can use a network ACL to control the traffic to and from the subnet for your NAT gateway. NAT gateways use ports 1024–65535. For more information, see [Control subnet traffic with network access control lists](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-network-acls.html).
- A NAT gateway receives a network interface. You can pick the private IPv4 address to assign to the interface or have it automatically assigned from the IPv4 address range of the subnet. You can view the network interface for the NAT gateway using the Amazon EC2 console. For more information, see [Viewing details about a network interface](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-eni.html#view_eni_details). You can't modify the attributes of this network interface.
- You can't route traffic to a NAT gateway through a VPC peering connection.
- You can't route traffic to a NAT gateway from Site-to-Site VPN or Direct Connect using a virtual private gateway. You can route traffic to a NAT gateway from Site-to-Site VPN or Direct Connect if you use a transit gateway instead of a virtual private gateway.
- NAT gateways support traffic with a **maximum transmission unit (MTU) of 8500**, but it's important to note the following:
  - The MTU of a network connection is the size, in bytes, of the largest permissible packet that can be passed over the connection. The larger the MTU of a connection, the more data that can be passed in a single packet.
  - Packets larger than 8500 bytes that arrive at the NAT gateway are dropped (or fragmented, if applicable).
  - To prevent potential packet loss when communicating with resources over the internet using a public NAT gateway, the MTU setting for your EC2 instances should not exceed 1500 bytes. For more information about checking and setting the MTU on an instance, see [Check and set the MTU on your Linux instance](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/network_mtu.html#set_mtu) in the Amazon EC2 User Guide.
  - NAT gateways support Path MTU Discovery (PMTUD) via FRAG_NEEDED ICMPv4 packets and Packet Too Big (PTB) ICMPv6 packets.
  - NAT gateways enforce Maximum Segment Size (MSS) clamping for all packets. For more information, see [RFC879](https://datatracker.ietf.org/doc/html/rfc879).

# Work with NAT gateways
## Control the use of NAT gateways
By default, users do not have permission to work with NAT gateways. You can create an IAM role with a policy attached that grants users permissions to create, describe, and delete NAT gateways. For more information, see [Identity and access management for Amazon VPC](https://docs.aws.amazon.com/vpc/latest/userguide/security-iam.html).

## Create a NAT gateway
Use the following procedure to create a NAT gateway.

### Related quotas
- You won't be able to create a public NAT gateway if you've exhausted the number of EIPs allocated to your account. For more information on EIP quotas and how to adjust them, see [Elastic IP addresses](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html#vpc-limits-eips).
- You can assign up to 8 private IPv4 addresses to your private NAT gateway. This limit is not adjustable.
- You are limited to associating 2 Elastic IP addresses to your public NAT gateway by default. You can increase this limit by requesting a quota adjustment. For more information, see [Elastic IP addresses](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html#vpc-limits-eips).

### To create a NAT gateway
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- In the navigation pane, choose **NAT gateways**.
- Choose **Create NAT gateway**.
- (Optional) Specify a name for the NAT gateway. This creates a tag where the key is ```Name``` and the value is the name that you specify.
- Select the subnet in which to create the NAT gateway.
- For Connectivity type, leave the default Public selection to create a public NAT gateway or choose Private to create a private NAT gateway. For more information about the difference between a public and private NAT gateway, see [NAT gateways](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-gateway.html).
- If you chose Public, do the following; otherwise, skip to step 8:
  - Choose an **Elastic IP allocation ID** to assign an EIP to the NAT gateway or choose **Allocate Elastic IP** to automatically allocate an EIP for the public NAT gateway. You are limited to associating 2 Elastic IP addresses to your public NAT gateway by default. You can increase this limit by requesting a quota adjustment. For more information, see [Elastic IP addresses](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html#vpc-limits-eips).
    - **Important**: When you assign an EIP to a public NAT gateway, the network border group of the EIP must match the network border group of the Availability Zone (AZ) that you're launching the public NAT gateway into. If it's not the same, the NAT gateway will fail to launch. You can see the network border group for the subnet's AZ by viewing the details of the subnet. Similarly, you can view the network border group of an EIP by viewing the details of the EIP address. For more information about network border groups and EIPs, see [1. Allocate an Elastic IP address](https://docs.aws.amazon.com/vpc/latest/userguide/WorkWithEIPs.html#allocate-eip).
  - (Optional) Choose **Additional settings** and, under **Private IP address - optional**, enter a private IPv4 address for the NAT gateway. If you don't enter an address, AWS will automatically assign a private IPv4 address to your NAT gateway at random from the subnet that your NAT gateway is in.
  - Skip to step 11.
- If you chose **Private**, for **Additional settings**, **Private IPv4 address assigning method**, choose one of the following:
  - **Auto-assign**: AWS chooses the primary private IPv4 address for the NAT gateway. For Number of auto-assigned private IPv4 addresses, you can optionally specify the number of secondary private IPv4 addresses for the NAT gateway. AWS chooses these IP addresses at random from the subnet for your NAT gateway.
  - **Custom**: For **Primary private IPv4 address**, choose the primary private IPv4 address for the NAT gateway. For **Secondary private IPv4 addresses**, you can optionally specify up to 7 secondary private IPv4 addresses for the NAT gateway.
- If you chose **Custom** in Step 8, skip this step. If you chose **Auto-assign**, under **Number of auto-assigned private IP addresses**, choose the number of secondary IPv4 addresses that you want AWS assign to this private NAT gateway. You can choose up to 7 IPv4 addresses.
  - **Note**: Secondary IPv4 addresses are optional and should be assigned or allocated when your workloads that use a NAT gateway exceed 55,000 concurrent connections to a single destination (the same destination IP, destination port, and protocol). Secondary IPv4 addresses increase the number of available ports, and therefore they increase the limit on the number of concurrent connections that your workloads can establish using a NAT gateway.
- If you chose **Auto-assign** in Step 9, skip this step. If you chose **Custom**, do the following:
  - Under **Primary private IPv4 address**, enter a private IPv4 address.
  - Under **Secondary private IPv4 address**, enter up to 7 secondary private IPv4 addresses
- (Optional) To add a tag to the NAT gateway, choose **Add new tag** and enter the key name and value. You can add up to 50 tags.
- Choose **Create a NAT gateway**.
- The initial status of the NAT gateway is ```Pending```. After the status changes to Available, the NAT gateway is ready for you to use. Be sure to update your route tables as needed. For examples, see [NAT gateway use cases](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-scenarios.html).

If the status of the NAT gateway changes to ```Failed```, there was an error during creation. For more information, see [NAT gateway creation fails](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-troubleshooting.html#nat-gateway-troubleshooting-failed).

## Edit secondary IP address associations
Each IPv4 address can support up to 55,000 simultaneous connections to each unique destination. A unique destination is identified by a unique combination of destination IP address, the destination port, and protocol (TCP/UDP/ICMP). You can increase this limit by associating up to 8 IPv4 addresses to your NAT gateways (1 primary IPv4 address and 7 secondary IPv4 addresses). You are limited to associating 2 Elastic IP addresses to your public NAT gateway by default. You can increase this limit by requesting a quota adjustment. For more information, see [Elastic IP addresses](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html#vpc-limits-eips).

You can use the [NAT gateway CloudWatch metrics](https://docs.aws.amazon.com/vpc/latest/userguide/metrics-dimensions-nat-gateway.html) *ErrorPortAllocation* and *PacketsDropCount* to determine if your NAT gateway is generating port allocation errors or dropping packets. To resolve this issue, add secondary IPv4 addresses to your NAT gateway.

### Considerations
- You can add secondary private IPv4 addresses when you create a private NAT gateway or after you create the NAT gateway using the procedure in this section. You can add secondary EIP addresses to public NAT gateways only after you create the NAT gateway by using the procedure in this section.
- Your NAT gateway can have up to 8 IPv4 addresses associated with it (1 primary IPv4 address and 7 secondary IPv4 addresses). You can assign up to 8 private IPv4 addresses to your private NAT gateway. You are limited to associating 2 Elastic IP addresses to your public NAT gateway by default. You can increase this limit by requesting a quota adjustment. For more information, see [Elastic IP addresses](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html#vpc-limits-eips).

### To edit secondary IPv4 address associations
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- In the navigation pane, choose **NAT gateways**.
- Select the NAT gateway whose secondary IPv4 address associations you want to edit.
- Choose **Actions**, and then choose **Edit secondary IP address associations**.
- If you are editing the secondary IPv4 address associations of a private NAT gateway, under **Action**, choose **Assign new IPv4 addresses** or **Unassign existing IPv4 addresses**. If you are editing the secondary IPv4 address associations of a public NAT gateway, under **Action**, choose **Associate new IPv4 addresses** or **Disassociate existing IPv4 addresses**.
- Do one of the following:
  - If you chose to assign or associate new IPv4 addresses, do the following:
    - This step is required. You must select a private IPv4 address. Choose the **Private IPv4 address assigning method**:
      - **Auto-assign**: AWS automatically chooses a primary private IPv4 address and you choose if you want AWS to assign up to 7 secondary private IPv4 addresses to assign to the NAT gateway. AWS automatically chooses and assigns them for you at random from the subnet that your NAT gateway is in.
      - **Custom**: Choose the primary private IPv4 address and up to 7 secondary private IPv4 addresses to assign to the NAT gateway.
    - Under **Elastic IP allocation ID**, choose an EIP to add as a secondary IPv4 address. This step is required. You must select an EIP along with a private IPv4 address. If you chose **Custom** for the **Private IP address assigning method**, you also must enter a private IPv4 address for each EIP that you add.
      - **Important**: When you assign a secondary EIP to a public NAT gateway, the network border group of the EIP must match the network border group of the Availability Zone (AZ) that the public NAT gateway is in. If it's not the same, the EIP will fail to assign. You can see the network border group for the subnet's AZ by viewing the details of the subnet. Similarly, you can view the network border group of an EIP by viewing the details of the EIP address. For more information about network border groups and EIPs, see [1. Allocate an Elastic IP address](https://docs.aws.amazon.com/vpc/latest/userguide/WorkWithEIPs.html#allocate-eip). 
      - Your NAT gateway can have up to 8 IP addresses associated with it. If this is a public NAT gateway, there is a default quota limit for EIPs per Region. For more information, see [Elastic IP addresses](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html#vpc-limits-eips).
  - If you chose to unassign or disassociate new IPv4 addresses, complete the following:
    - Under **Existing secondary IP address to unassign**, select the secondary IP addresses that you want to unassign.
    - (optional) Under **Connection drain duration**, enter the maximum amount of time to wait (in seconds) before forcibly releasing the IP addresses if connections are still in progress. If you don't enter a value, the default value is 350 seconds.
- Choose Save changes.

If the status of the NAT gateway changes to ```Failed```, there was an error during creation. For more information, see [NAT gateway creation fails](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-troubleshooting.html#nat-gateway-troubleshooting-failed).

## Tag a NAT gateway
You can tag your NAT gateway to help you identify it or categorize it according to your organization's needs. For information about working with tags, see [Tagging your Amazon EC2 resources](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html) in the *Amazon EC2 User Guide*.

Cost allocation tags are supported for NAT gateways. Therefore, you can also use tags to organize your AWS bill and reflect your own cost structure. For more information, see [Using cost allocation tags](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/cost-alloc-tags.html) in the *AWS Billing User Guide*. For more information about setting up a cost allocation report with tags, see [Monthly cost allocation report](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/configurecostallocreport.html) in *About AWS Account Billing*.

### To tag a NAT gateway
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- In the navigation pane, choose **NAT gateways**.
- Select the NAT gateway that you want to tag and choose **Actions**. Then choose **Manage tags**
- Choose **Add new tag**, and define a **Key** and **Value** for the tag. You can add up to 50 tags
- Choose Save.

## Delete a NAT gateway
If you no longer need a NAT gateway, you can delete it. After you delete a NAT gateway, its entry remains visible in the Amazon VPC console for about an hour, after which it's automatically removed. You can't remove this entry yourself.

**Deleting a NAT gateway disassociates its Elastic IP address, but does not release the address from your account**. If you delete a NAT gateway, the NAT gateway routes remain in a ```blackhole``` status until you delete or update the routes.

### To delete a NAT gateway
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- In the navigation pane, choose **NAT gateways**.
- Select the radio button for the NAT gateway, and then choose **Actions, Delete NAT gateway**.
- When prompted for confirmation, enter ```delete``` and then choose **Delete**.
- If you no longer need the Elastic IP address that was associated with a public NAT gateway, we recommend that you release it. For more information, see [5. Release an Elastic IP address](https://docs.aws.amazon.com/vpc/latest/userguide/WorkWithEIPs.html#release-eip).

## Command line overview
You can perform the tasks described on this page using the command line.

### Assign a private IPv4 address to a private NAT gateway
- [assign-private-nat-gateway-address](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/assign-private-nat-gateway-address.html) (AWS CLI)
- [Register-EC2PrivateNatGatewayAddress](https://docs.aws.amazon.com/powershell/latest/reference/items/Register-EC2PrivateNatGatewayAddress.html) (AWS Tools for Windows PowerShell)

### Associate Elastic IP addresses (EIPs) and private IPv4 addresses with a public NAT gateway
- [associate-nat-gateway-address](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/associate-nat-gateway-address.html) (AWS CLI)
- [Register-EC2NatGatewayAddress](https://docs.aws.amazon.com/powershell/latest/reference/items/Register-EC2NatGatewayAddress.html) (AWS Tools for Windows PowerShell)

### Create a NAT gateway
- [create-nat-gateway](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/create-nat-gateway.html) (AWS CLI)
- [New-EC2NatGateway](https://docs.aws.amazon.com/powershell/latest/reference/items/New-EC2NatGateway.html) (AWS Tools for Windows PowerShell)

### Delete a NAT gateway
- [delete-nat-gateway](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/delete-nat-gateway.html) (AWS CLI)
- [Remove-EC2NatGateway](https://docs.aws.amazon.com/powershell/latest/reference/items/Remove-EC2NatGateway.html) (AWS Tools for Windows PowerShell)

### Describe a NAT gateway
- [describe-nat-gateways](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/describe-nat-gateways.html) (AWS CLI)
- [Get-EC2NatGateway](https://docs.aws.amazon.com/powershell/latest/reference/items/Get-EC2NatGateway.html) (AWS Tools for Windows PowerShell)

### Disassociate secondary Elastic IP addresses (EIPs) from a public NAT gateway
- [disassociate-nat-gateway-address](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/disassociate-nat-gateway-address.html) (AWS CLI)
- [Unregister-EC2NatGatewayAddress](https://docs.aws.amazon.com/powershell/latest/reference/items/Unregister-EC2NatGatewayAddress.html) (AWS Tools for Windows PowerShell)

### Tag a NAT gateway
- [create-tags](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/create-tags.html) (AWS CLI)
- [New-EC2Tag](https://docs.aws.amazon.com/powershell/latest/reference/items/New-EC2Tag.html) (AWS Tools for Windows PowerShell)

### Unassign secondary IPv4 addresses from a private NAT gateway
- [unassign-private-nat-gateway-address](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/unassign-private-nat-gateway-address.html) (AWS CLI)
- [Unregister-EC2PrivateNatGatewayAddress](https://docs.aws.amazon.com/powershell/latest/reference/items/Unregister-EC2PrivateNatGatewayAddress.html) (AWS Tools for Windows PowerShell)

# NAT Gateway Usecases
https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-scenarios.html

The following are example use cases for public and private NAT gateways.

## Scenarios
- [Access the internet from a private subnet](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-scenarios.html#public-nat-internet-access)
- [Access your network using allow-listed IP addresses](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-scenarios.html#private-nat-allowed-range)
- [Enable communication between overlapping networks](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-scenarios.html#private-nat-overlapping-networks)

## Access the internet from a private subnet
You can use a public NAT gateway to enable instances in a private subnet to send outbound traffic to the internet, while preventing the internet from establishing connections to the instances.

### Overview
The following diagram illustrates this use case. There are two Availability Zones, with two subnets in each Availability Zone. The route table for each subnet determines how traffic is routed. In Availability Zone A, the instances in the public subnet can reach the internet through a route to the internet gateway, while the instances in the private subnet have no route to the internet. In Availability Zone B, the public subnet contains a NAT gateway, and the instances in the private subnet can reach the internet through a route to the NAT gateway in the public subnet. Both private and public NAT gateways map the source private IPv4 address of the instances to the private IPv4 address of the private NAT gateway, but in the case of a public NAT gateway, the internet gateway then maps the private IPv4 address of the public NAT gateway to the Elastic IP address associated with the NAT gateway. When sending response traffic to the instances, whether it's a public or private NAT gateway, the NAT gateway translates the address back to the original source IP address.

pic

Note that if the instances in the private subnet in Availability Zone A also need to reach the internet, you can create a route from this subnet to the NAT gateway in Availability Zone B. Alternatively, you can improve resiliency by creating a NAT gateway in each Availability Zone that contains resources that require internet access. For an example diagram, see [Example: VPC with servers in private subnets and NAT](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-example-private-subnets-nat.html).

### Routing
The following is the route table associated with the public subnet in Availability Zone A. The first entry is the local route; it enables the instances in the subnet to communicate with other instances in the VPC using private IP addresses. The second entry sends all other subnet traffic to the internet gateway, which enables the instances in the subnet to access the internet.

| Destination | Target              |
| ----------- | ------------------- |
| VPC CIDR    | local               |
| 0.0.0.0/0   | internet-gateway-id |

The following is the route table associated with the private subnet in Availability Zone A. The entry is the local route, which enables the instances in the subnet to communicate with other instances in the VPC using private IP addresses. The instances in this subnet have no access to the internet.

| Destination | Target              |
| ----------- | ------------------- |
| VPC CIDR    | local               |

The following is the route table associated with the public subnet in Availability Zone B. The first entry is the local route, which enables the instances in the subnet to communicate with other instances in the VPC using private IP addresses. The second entry sends all other subnet traffic to the internet gateway, which enables the NAT gateway in the subnet to access the internet.

| Destination | Target              |
| ----------- | ------------------- |
| VPC CIDR    | local               |
| 0.0.0.0/0   | internet-gateway-id |

The following is the route table associated with the private subnet in Availability Zone B. The first entry is the local route; it enables the instances in the subnet to communicate with other instances in the VPC using private IP addresses. The second entry sends all other subnet traffic to the NAT gateway.

| Destination | Target         |
| ----------- | -------------- |
| VPC CIDR    | local          |
| 0.0.0.0/0   | nat-gateway-id |

For more information, see [Change a subnet route table](https://docs.aws.amazon.com/vpc/latest/userguide/WorkWithRouteTables.html).

### Test the public NAT gateway
After you've created your NAT gateway and updated your route tables, you can ping remote addresses on the internet from an instance in your private subnet to test whether it can connect to the internet. For an example of how to do this, see [Test the internet connection](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-scenarios.html#nat-gateway-testing-example).

If you can connect to the internet, you can also test whether internet traffic is routed through the NAT gateway:
- Trace the route of traffic from an instance in your private subnet. To do this, run the ```traceroute``` command from a Linux instance in your private subnet. In the output, you should see the private IP address of the NAT gateway in one of the hops (usually the first hop).
- Use a third-party website or tool that displays the source IP address when you connect to it from an instance in your private subnet. The source IP address should be the elastic IP address of the NAT gateway.

If these tests fail, see [Troubleshoot NAT gateways](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-troubleshooting.html).

### Test the internet connection
The following example demonstrates how to test whether an instance in a private subnet can connect to the internet.
- Launch an instance in your public subnet (use this as a bastion host). In the launch wizard, ensure that you select an **Amazon Linux AMI**, and assign a public IP address to your instance. Ensure that your security group rules allow inbound SSH traffic from the range of IP addresses for your local network, and outbound SSH traffic to the IP address range of your private subnet (you can also use 0.0.0.0/0 for both inbound and outbound SSH traffic for this test).
- Launch an instance in your private subnet. In the launch wizard, ensure that you select an Amazon Linux AMI. Do not assign a public IP address to your instance. Ensure that your security group rules allow inbound SSH traffic from the private IP address of your instance that you launched in the public subnet, and all outbound ICMP traffic. You must choose the same key pair that you used to launch your instance in the public subnet.
- Configure SSH agent forwarding on your local computer, and connect to your bastion host in the public subnet. For more information, see [To configure SSH agent forwarding for Linux or macOS](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-scenarios.html#ssh-forwarding-linux) or To configure [SSH agent forwarding for Windows](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-scenarios.html#ssh-forwarding-windows).
- From your bastion host, connect to your instance in the private subnet, and then test the internet connection from your instance in the private subnet. For more information, see To [test the internet connection](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-scenarios.html#test-internet-connection).

#### To configure SSH agent forwarding for Linux or macOS
From your local machine, add your private key to the authentication agent.

For Linux, use the following command.
```sh
ssh-add -c mykeypair.pem
```

For macOS, use the following command.
```sh
ssh-add -K mykeypair.pem
```

Connect to your instance in the public subnet using the ```-A``` option to enable SSH agent forwarding, and use the instance's public address, as shown in the following example.

```sh
ssh -A ec2-user@54.0.0.123
```

#### To test the internet connection
From your instance in the public subnet, connect to your instance in your private subnet by using its private IP address as shown in the following example.

```sh
ssh ec2-user@10.0.1.123
```
From your private instance, test that you can connect to the internet by running the ping command for a website that has ICMP enabled.

```sh
ping ietf.org
PING ietf.org (4.31.198.44) 56(84) bytes of data.
64 bytes from mail.ietf.org (4.31.198.44): icmp_seq=1 ttl=47 time=86.0 ms
64 bytes from mail.ietf.org (4.31.198.44): icmp_seq=2 ttl=47 time=75.6 ms
...
```

Press Ctrl+C on your keyboard to cancel the ping command. If the ```ping``` command fails, see [Instances cannot access the internet](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-troubleshooting.html#nat-gateway-troubleshooting-no-internet-connection).

(Optional) If you no longer require your instances, terminate them. For more information, see [Terminate your instance](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/terminating-instances.html) in the Amazon EC2 User Guide.

## Access your network using allow-listed IP addresses
You can use a private NAT gateway to enable communication from your VPCs to your on-premises network using a pool of allow-listed addresses. Instead of assigning each instance a separate IP address from the allow-listed IP address range, you can route traffic from the subnet that is destined for the on-premises network through a private NAT gateway with an IP address from the allow-listed IP address range.

### Overview
The following diagram shows how instances can access on-premises resources through AWS VPN. Traffic from the instances is routed to a virtual private gateway, over the VPN connection, to the customer gateway, and then to the destination in the on-premises network. However, suppose that the destination allows traffic only from a specific IP address range, such as 100.64.1.0/28. This would prevent traffic from these instances from reaching the on-premises network.

pic

The following diagram shows the key components of the configuration for this scenario. The VPC has its original IP address range plus the allowed IP address range. The VPC has a subnet from the allowed IP address range with a private NAT gateway. Traffic from the instances that is destined for the on-premises network is sent to the NAT gateway before being routed to the VPN connection. The on-premises network receives the traffic from the instances with the source IP address of the NAT gateway, which is from the allowed IP address range.

### Resources
Create or update resources as follows:

- Associate the allowed IP address range with the VPC.
- Create a subnet in the VPC from the allowed IP address range.
- Create a private NAT gateway in the new subnet.
- Update the route table for the subnet with the instances to send traffic destined for the on-premises network to the NAT gateway. Add a route to the route table for the subnet with the private NAT gateway that sends traffic destined for the on-premises network to the virtual private gateway.

### Routing
The following is the route table associated with the first subnet. There is a local route for each VPC CIDR. Local routes enable resources in the subnet to communicate with other resources in the VPC using private IP addresses. The third entry sends traffic destined for the on-premises network to the private NAT gateway.

| Destination    | Target         |
| -------------- | -------------- |
| 10.0.0.0/16    | local          |
| 100.64.1.0/24  | local          |
| 192.168.0.0/16 | nat-gateway-id |

The following is the route table associated with the second subnet. There is a local route for each VPC CIDR. Local routes enable resources in the subnet to communicate with other resources in the VPC using private IP addresses. The third entry sends traffic destined for the on-premises network to the virtual private gateway.

| Destination    | Target |
| -------------- | ------ |
| 10.0.0.0/16    | local  |
| 100.64.1.0/24  | local  |
| 192.168.0.0/16 | vgw-id |

## Enable communication between overlapping networks
You can use a private NAT gateway to enable communication between networks even if they have overlapping CIDR ranges. For example, suppose that the instances in VPC A need to access the services provided by the instances in VPC B.

pic

### Overview
The following diagram shows the key components of the configuration for this scenario. First, your IP management team determines which address ranges can overlap (non-routable address ranges) and which can't (routable address ranges). The IP management team allocates address ranges from the pool of routable address ranges to projects on request.

Each VPC has its original IP address range, which is non-routable, plus the routable IP address range assigned to it by the IP management team. VPC A has a subnet from its routable range with a private NAT gateway. The private NAT gateway gets its IP address from its subnet. VPC B has a subnet from its routable range with an Application Load Balancer. The Application Load Balancer gets its IP addresses from its subnets.

Traffic from an instance in the non-routable subnet of VPC A that is destined for the instances in the non-routable subnet of VPC B is sent through the private NAT gateway and then routed to the transit gateway. The transit gateway sends the traffic to the Application Load Balancer, which routes the traffic to one of the target instances in the non-routable subnet of VPC B. The traffic from the transit gateway to the Application Load Balancer has the source IP address of the private NAT gateway. Therefore, response traffic from the load balancer uses the address of the private NAT gateway as its destination. The response traffic is sent to the transit gateway and then routed to the private NAT gateway, which translates the destination to the instance in the non-routable subnet of VPC A.

pic

### Resources
Create or update resources as follows:

- Associate the assigned routable IP address ranges with their respective VPCs.
- Create a subnet in VPC A from its routable IP address range, and create a private NAT gateway in this new subnet.
- Create a subnet in VPC B from its routable IP address range, and create an Application Load Balancer in this new subnet. Register the instances in the non-routable subnet with the target group for the load balancer.
- Create a transit gateway to connect the VPCs. Be sure to disable route propagation. When you attach each VPC to the transit gateway, use the routable address range of the VPC.
- Update the route table of the non-routable subnet in VPC A to send all traffic destined for the routable address range of VPC B to the private NAT gateway. Update the route table of the routable subnet in VPC A to send all traffic destined for the routable address range of VPC B to the transit gateway.
- Update the route table of the routable subnet in VPC B to send all traffic destined for the routable address range of VPC A to the transit gateway.

### Routing
The following is the route table for the non-routable subnet in VPC A.

| Destination   | Target         |
| ------------- | -------------- |
| 10.0.0.0/16   | local          |
| 100.64.1.0/24 | local          |
| 100.64.2.0/24 | nat-gateway-id |

The following is the route table for the routable subnet in VPC A.

| Destination   | Target             |
| ------------- | ------------------ |
| 10.0.0.0/16   | local              |
| 100.64.1.0/24 | local              |
| 100.64.2.0/24 | transit-gateway-id |

The following is the route table for the non-routable subnet in VPC B.

| Destination   | Target |
| ------------- | ------ |
| 10.0.0.0/16   | local  |
| 100.64.2.0/24 | local  |

The following is the route table for the routable subnet in VPC B.
| Destination   | Target             |
| ------------- | ------------------ |
| 10.0.0.0/16   | local              |
| 100.64.2.0/24 | local              |
| 100.64.1.0/24 | transit-gateway-id |

The following is the transit gateway route table.
| CIDR          | Attachment           | Route Type |
| ------------- | -------------------- | ---------- |
| 100.64.1.0/24 | Attachment for VPC A | static     |
| 100.64.2.0/24 | Attachment for VPC B | static     |


