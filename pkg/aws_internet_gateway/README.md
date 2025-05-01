https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Internet_Gateway.html

An internet gateway is a horizontally scaled, redundant, and highly available VPC component that allows communication between your VPC and the internet. It supports IPv4 and IPv6 traffic. It does not cause availability risks or bandwidth constraints on your network traffic.

An internet gateway enables resources in your **public subnets** (such as EC2 instances) to connect to the internet **if the resource has a public IPv4 address or an IPv6 address**. Similarly, resources on the internet can initiate a connection to resources in your subnet using the public IPv4 address or IPv6 address. For example, an internet gateway enables you to connect to an EC2 instance in AWS using your local computer.

An internet gateway provides a target in your VPC route tables for internet-routable traffic. For communication using IPv4, the internet gateway also performs network address translation (NAT). For more information, see [IP addresses and NAT](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Internet_Gateway.html#ip-addresses-and-nat).

# Pricing
There is no charge for an internet gateway, but there are data transfer charges for EC2 instances that use internet gateways. For more information, see [Amazon EC2 On-Demand Pricing](https://aws.amazon.com/ec2/pricing/on-demand/).

# Internet gateway basics
To use an internet gateway, you must attach it to a VPC and configure routing.

## Routing configuration
**If a subnet is associated with a route table that has a route to an internet gateway, it's known as a public subnet**. If a subnet is associated with a route table that does not have a route to an internet gateway, it's known as a private subnet.

In your public subnet's route table, you can specify a route for the internet gateway to all destinations not explicitly known to the route table (0.0.0.0/0 for IPv4 or ::/0 for IPv6). Alternatively, you can scope the route to a narrower range of IP addresses; for example, the public IPv4 addresses of your company’s public endpoints outside of AWS, or the Elastic IP addresses of other Amazon EC2 instances outside your VPC.

## Internet gateway diagram
In the following diagram, the subnet in Availability Zone A is a public subnet because its route table has a route that sends all internet-bound IPv4 traffic to the internet gateway. The instances in the public subnet must have public IP addresses or Elastic IP addresses to enable communication with the internet over the internet gateway. For comparison, the subnet in Availability Zone B is a private subnet because its route table does not have a route to the internet gateway. Because there is no route to the internet gateway, instances in the private subnet can't communicate with the internet, even if they have public IP addresses.

pic

## IP addresses and NAT
To enable communication over the internet for IPv4, your instance must have a public IPv4 address. You can either configure your VPC to automatically assign public IPv4 addresses to your instances, or you can assign Elastic IP addresses to your instances. Your instance is only aware of the private (internal) IP address space defined within the VPC and subnet. The internet gateway logically provides the one-to-one NAT on behalf of your instance, so that when traffic leaves your VPC subnet and goes to the internet, the **reply address field is set to the public IPv4 address or Elastic IP address of your instance, and not its private IP address**. Conversely, traffic that's destined for the public IPv4 address or Elastic IP address of your instance has **its destination address translated into the instance's private IPv4 address before the traffic is delivered to the VPC**.

To enable communication over the internet for IPv6, your VPC and subnet must have an associated IPv6 CIDR block, and your instance must be assigned an IPv6 address from the range of the subnet. IPv6 addresses are globally unique, and therefore public by default.

### Internet access for default and nondefault VPCs
The following table provides an overview of whether your VPC automatically comes with the components required for internet access over IPv4 or IPv6.

| Component                                       | Default VPC          | Nondefault VPC         |
| ----------------------------------------------- | -------------------- | ---------------------- |
| Internet gateway                                | Yes                  | No                     |
| Route table with route to internet gateway      | Yes                  | No                     |
| for IPv4 traffic (0.0.0.0/0)                    |                      |                        |
| Route table with route to internet gateway      | No                   | No                     |
| for IPv6 traffic (::/0)                         |                      |                        |
| Public IPv4 address automatically assigned to   | Yes (default subnet) | No (nondefault subnet) |
| instance launched into subnet                   |                      |                        |
| IPv6 address automatically assigned to instance | No (default subnet)  | No (nondefault subnet) |
| launched into subnet                            |                      |                        |
		
# Add internet access to a subnet
To support internet access from a subnet in a nondefault VPC using an internet gateway. You must create the internet gateway, attach it to the VPC, and configure routing for the subnet.

After you configure internet access for your subnet, you must ensure that resources in the subnet can access the internet. For example, your EC2 instances must have a public IPv4 or IPv6 address, and the security groups for your instances must allow specific traffic to and from the internet.

Alternatively, to provide your instances with internet access without assigning them a public IP address, use a NAT device instead. For more information, see [NAT devices](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat.html).

To remove internet access, you can detach the internet gateway from your VPC and then delete it. For more information, see [Delete an internet gateway](https://docs.aws.amazon.com/vpc/latest/userguide/delete-igw.html).

## Step 1: Create an internet gateway
Use the following procedure to create an internet gateway.

### To create an internet gateway using the console
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- In the navigation pane, choose **Internet gateways**.
- Choose **Create internet gateway**.
- (Optional) Enter a name for your internet gateway.
- (Optional) To add a tag, choose **Add new tag** and enter the tag key and value.
- Choose **Create internet gateway**.
- (Optional) To attach the internet gateway to a VPC now, choose **Attach to a VPC** from the banner at the top of the screen, select an available VPC, and then choose Attach internet gateway. Otherwise, you can attach your internet gateway to a VPC at another time.

### To create an internet gateway using the command line
- [create-internet-gateway](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/create-internet-gateway.html) (AWS CLI)
- [New-EC2InternetGateway](https://docs.aws.amazon.com/powershell/latest/reference/items/New-EC2InternetGateway.html) (AWS Tools for Windows PowerShell)

## Step 2: Attach the internet gateway to the VPC
To use an internet gateway, you must attach it to a VPC.

### To attach an internet gateway to a VPC using the console
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- In the navigation pane, choose **Internet gateways**.
- Select the check box for the internet gateway.
- To attach it, choose **Actions**, Attach to VPC, select an available VPC, and choose **Attach internet gateway**.
- To detach it, choose **Actions**, Detach from VPC and choose **Detach internet gateway**. When prompted for confirmation, choose **Detach internet gateway**.

### To attach an internet gateway to a VPC using the command line
- [attach-internet-gateway](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/attach-internet-gateway.html) (AWS CLI)
- [Add-EC2InternetGateway](https://docs.aws.amazon.com/powershell/latest/reference/items/Add-EC2InternetGateway.html) (AWS Tools for Windows PowerShell)

## Step 3: Add a route to the subnet route table
The route table for the subnet must have a route that sends internet traffic to the internet gateway.

### To configure the subnet route table using the console
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- In the navigation pane, choose **Route tables**.
- Select the route table for the subnet. By default, a subnet uses the main route table for the VPC. Alternatively, you can [create a custom route table](https://docs.aws.amazon.com/vpc/latest/userguide/WorkWithRouteTables.html#CustomRouteTable) and then [associate the subnet with the new route table](https://docs.aws.amazon.com/vpc/latest/userguide/WorkWithRouteTables.html#AssociateSubnet).
- On the **Routes** tab, choose **Edit routes** and then choose **Add route**.
- Enter 0.0.0.0/0 for **Destination** and select the internet gateway for **Target**.
- Choose **Save changes**.

### To configure the subnet route table using the command line
- [create-route](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/create-route.html) (AWS CLI)
- [New-EC2Route](https://docs.aws.amazon.com/powershell/latest/reference/items/New-EC2Route.html) (AWS Tools for Windows PowerShell)