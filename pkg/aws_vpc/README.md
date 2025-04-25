https://docs.aws.amazon.com/vpc/latest/userguide/what-is-amazon-vpc.html

With Amazon Virtual Private Cloud (Amazon VPC), you can launch AWS resources in a logically isolated virtual network that you've defined. This virtual network closely resembles a traditional network that you'd operate in your own data center, with the benefits of using the scalable infrastructure of AWS.

The following diagram shows an example VPC. The VPC has one subnet in each of the Availability Zones in the Region, EC2 instances in each subnet, and an internet gateway to allow communication between the resources in your VPC and the internet.

pic

For more information, see [Amazon Virtual Private Cloud (Amazon VPC)](https://aws.amazon.com/vpc/).

# Features
The following features help you configure a VPC to provide the connectivity that your applications need:

## Virtual private clouds (VPC)
A VPC is a virtual network that closely resembles a traditional network that you'd operate in your own data center. After you create a VPC, you can add subnets.

## Subnets
A subnet is a range of IP addresses in your VPC. **A subnet must reside in a single Availability Zone**. After you add subnets, you can deploy AWS resources in your VPC.

## IP addressing
You can assign [IP addresses](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-ip-addressing.html), both IPv4 and IPv6, to your VPCs and subnets. You can also bring your public IPv4 addresses and IPv6 GUA addresses to AWS and allocate them to resources in your VPC, such as EC2 instances, NAT gateways, and Network Load Balancers.

## Routing
Use [route tables](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html) to determine where network traffic from your subnet or gateway is directed.

## Gateways and endpoints
A [gateway](https://docs.aws.amazon.com/vpc/latest/userguide/extend-intro.html) connects your VPC to another network. For example, use an [internet gateway](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Internet_Gateway.html) to connect your VPC to the internet. 

Use a [VPC endpoint](https://docs.aws.amazon.com/vpc/latest/privatelink/privatelink-access-aws-services.html) to connect to AWS services privately, without the use of an internet gateway or NAT device.

## Peering connections
Use a [VPC peering connection](https://docs.aws.amazon.com/vpc/latest/peering/) to route traffic between the resources in two VPCs.

## Traffic Mirroring
[Copy network traffic](https://docs.aws.amazon.com/vpc/latest/mirroring/) from network interfaces and send it to security and monitoring appliances for deep packet inspection.

## Transit gateways
Use a [transit gateway](https://docs.aws.amazon.com/vpc/latest/userguide/extend-tgw.html), which acts as a central hub, to route traffic between your VPCs, VPN connections, and AWS Direct Connect connections.

## VPC Flow Logs
A [flow log](https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs.html) captures information about the IP traffic going to and from network interfaces in your VPC.

## VPN connections
Connect your VPCs to your on-premises networks using [AWS Virtual Private Network (AWS VPN)](https://docs.aws.amazon.com/vpc/latest/userguide/vpn-connections.html).

# Getting started with Amazon VPC
Your AWS account includes a [default VPC](https://docs.aws.amazon.com/vpc/latest/userguide/default-vpc.html) in each AWS Region. Your default VPCs are configured such that you can immediately start launching and connecting to EC2 instances. For more information, see [Plan your VPC](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-getting-started.html).

You can choose to create additional VPCs with the subnets, IP addresses, gateways and routing that you need. For more information, see [Create a VP](https://docs.aws.amazon.com/vpc/latest/userguide/create-vpc.html)C.

# Working with Amazon VPC
You can create and manage your VPCs using any of the following interfaces:
- **AWS Management Console** — Provides a web interface that you can use to access your VPCs.
- **AWS Command Line Interface (AWS CLI)** — Provides commands for a broad set of AWS services, including Amazon VPC, and is supported on Windows, Mac, and Linux. For more information, see [AWS Command Line Interface](https://aws.amazon.com/cli/).
- **AWS SDKs** — Provides language-specific APIs and takes care of many of the connection details, such as calculating signatures, handling request retries, and error handling. For more information, see [AWS SDKs](https://aws.amazon.com/developer/tools/).
- **Query API** — Provides low-level API actions that you call using HTTPS requests. Using the Query API is the most direct way to access Amazon VPC, but it requires that your application handle low-level details such as generating the hash to sign the request, and error handling. For more information, see [Amazon VPC actions](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/OperationList-query-vpc.html) in the *Amazon EC2 API Reference*.

# Pricing for Amazon VPC
There's no additional charge for using a VPC. There are, however, charges for some VPC components, such as NAT gateways, IP Address Manager, traffic mirroring, Reachability Analyzer, and Network Access Analyzer. For more information, see [Amazon VPC Pricing](https://aws.amazon.com/vpc/pricing/).

Nearly all resources that you launch in your virtual private cloud (VPC) provide you with an IP address for connectivity. The vast majority of resources in your VPC use private IPv4 addresses. Resources that require direct access to the internet over IPv4, however, use public IPv4 addresses.

Amazon VPC enables you to launch managed services, such as Elastic Load Balancing, Amazon RDS, and Amazon EMR, without having a VPC set up beforehand. It does this by using the [default VPC](https://docs.aws.amazon.com/vpc/latest/userguide/default-vpc.html) in your account if you have one. Any public IPv4 addresses provisioned to your account by the managed service will be charged. These charges will be associated with Amazon VPC service in your AWS Cost and Usage Report.

## Pricing for public IPv4 addresses
A *public IPv4 address* is an IPv4 address that is routable from the internet. A public IPv4 address is necessary for a resource to be directly reachable from the internet over IPv4.

If you are an existing or new [AWS Free Tier](https://aws.amazon.com/free/) customer, you get 750 hours of public IPv4 address usage with the EC2 service at no charge. If you are not using the EC2 service in the AWS Free Tier, Public IPv4 addresses are charged. For specific pricing information, see the *Public IPv4 address tab* in [Amazon VPC Pricing](https://aws.amazon.com/vpc/pricing/).

Private IPv4 addresses ([RFC 1918](https://datatracker.ietf.org/doc/html/rfc1918)) are not charged. For more information about how public IPv4 addresses are charged for shared VPCs, see [Billing and metering for the owner and participants](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-sharing.html#vpc-sharing-permissions).

Public IPv4 addresses have the following types:
- **Elastic IP addresses (EIPs)**: Static, public IPv4 addresses provided by Amazon that you can associate with an EC2 instance, elastic network interface, or AWS resource.
- **EC2 public IPv4 addresses**: Public IPv4 addresses assigned to an EC2 instance by Amazon (if the EC2 instance is launched into a default subnet or if the instance is launched into a subnet that’s been configured to automatically assign a public IPv4 address).
- **BYOIPv4 addresses**: Public IPv4 addresses in the IPv4 address range that you’ve brought to AWS using [Bring your own IP addresses (BYOIP)](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-byoip.html).
- **Service-managed IPv4 addresses**: Public IPv4 addresses automatically provisioned on AWS resources and managed by an AWS service. For example, public IPv4 addresses on Amazon ECS, Amazon RDS, or Amazon WorkSpaces.

The following list shows the most common AWS services that can use public IPv4 addresses.
- Amazon AppStream 2.0
- [AWS Client VPN](https://docs.aws.amazon.com/vpn/latest/clientvpn-admin/what-is.html#what-is-pricing)
- AWS Database Migration Service
- Amazon EC2
- Amazon Elastic Container Service
- Amazon EKS
- Amazon EMR
- Amazon GameLift Servers
- AWS Global Accelerator
- AWS Mainframe Modernization
- Amazon Managed Streaming for Apache Kafka
- Amazon MQ
- Amazon RDS
- Amazon Redshift
- AWS Site-to-Site VPN
- Amazon VPC NAT gateway
- Amazon WorkSpaces
- Elastic Load Balancing

# How Amazon VPC works
With Amazon Virtual Private Cloud (Amazon VPC), you can launch AWS resources in a logically isolated virtual network that you've defined. This virtual network closely resembles a traditional network that you'd operate in your own data center, with the benefits of using the scalable infrastructure of AWS.

The following is a visual representation of a VPC and its resources from the Preview pane shown when you create a VPC using the AWS Management Console. For an existing VPC, you can access this visualization on the [Resource map](https://docs.aws.amazon.com/vpc/latest/userguide/view-vpc-resource-map.html) tab. This example shows the resources that are initially selected on the Create VPC page when you choose to create the VPC plus other networking resources. This VPC is configured with an IPv4 CIDR and an Amazon-provided IPv6 CIDR, subnets in two Availability Zones, three route tables, an internet gateway, and a gateway endpoint. Because we've selected the internet gateway, the visualization indicates that traffic from the public subnets is routed to the internet because the corresponding route table sends the traffic to the internet gateway.

## VPCs and subnets
A *virtual private cloud* (VPC) is a virtual network dedicated to your AWS account. It is logically isolated from other virtual networks in the AWS Cloud. You can specify an IP address range for the VPC, add subnets, add gateways, and associate security groups.

A *subnet* is a range of IP addresses in your VPC. You launch AWS resources, such as Amazon EC2 instances, into your subnets. You can connect a subnet to the internet, other VPCs, and your own data centers, and route traffic to and from your subnets using route tables.

### Learn more
- [IP addressing for your VPCs and subnets](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-ip-addressing.html)
- [Configure a virtual private cloud](https://docs.aws.amazon.com/vpc/latest/userguide/configure-your-vpc.html)
- [Subnets for your VPC](https://docs.aws.amazon.com/vpc/latest/userguide/configure-subnets.html)

## Default and nondefault VPCs
If your account was created after December 4, 2013, it comes with a default VPC in each Region. A default VPC is configured and ready for you to use. For example, it has a default subnet in each Availability Zone in the Region, an attached internet gateway, a route in the main route table that sends all traffic to the internet gateway, and DNS settings that automatically assign public DNS hostnames to instances with public IP addresses and enable DNS resolution through the Amazon-provided DNS server (see [DNS attributes for your VPC](https://docs.aws.amazon.com/vpc/latest/userguide/AmazonDNS-concepts.html#vpc-dns-support)). 

Therefore, an EC2 instance that is launched in a default subnet automatically has access to the internet. If you have a default VPC in a Region and you don't specify a subnet when you launch an EC2 instance into that Region, we choose one of the default subnets and launch the instance into that subnet.

You can also create your own VPC, and configure it as you need. This is known as a *nondefault VPC*. Subnets that you create in your nondefault VPC and additional subnets that you create in your default VPC are called *nondefault subnets*.

### Learn more
- [Default VPCs](https://docs.aws.amazon.com/vpc/latest/userguide/default-vpc.html)
- [Create a VPC](https://docs.aws.amazon.com/vpc/latest/userguide/create-vpc.html)

## Route tables
A route table contains a set of rules, called routes, that are used to determine where network traffic from your VPC is directed. You can explicitly associate a subnet with a particular route table. Otherwise, the subnet is implicitly associated with the main route table.

Each route in a route table specifies the range of IP addresses where you want the traffic to go (the destination) and the gateway, network interface, or connection through which to send the traffic (the target).

### Learn more
[Configure route tables](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html)

## Access the internet
You control how the instances that you launch into a VPC access resources outside the VPC.

A default VPC includes an internet gateway, and **each default subnet is a public subnet**. Each instance that you launch into a default subnet has a private IPv4 address and a public IPv4 address. These instances can communicate with the internet through the internet gateway. An internet gateway enables your instances to connect to the internet through the Amazon EC2 network edge.

By default, each instance that you launch into a nondefault subnet has a private IPv4 address, but no public IPv4 address, unless you specifically assign one at launch, or you modify the subnet's public IP address attribute. These instances can communicate with each other, but can't access the internet.

You can enable internet access for an instance launched into a nondefault subnet by attaching an internet gateway to its VPC (if its VPC is not a default VPC) and associating an Elastic IP address with the instance.

Alternatively, to allow an instance in your VPC to initiate outbound connections to the internet but prevent unsolicited inbound connections from the internet, you can use a **network address translation (NAT)** device. **NAT maps multiple private IPv4 addresses to a single public IPv4 address**. 

You can configure the NAT device with an Elastic IP address and connect it to the internet through an internet gateway. This makes it possible for an instance in a private subnet to connect to the internet through the NAT device, routing traffic from the instance to the internet gateway and any responses to the instance.

If you associate an IPv6 CIDR block with your VPC and assign IPv6 addresses to your instances, instances can connect to the internet over IPv6 through an internet gateway. Alternatively, instances can initiate outbound connections to the internet over IPv6 using an egress-only internet gateway. IPv6 traffic is separate from IPv4 traffic; your route tables must include separate routes for IPv6 traffic.

### Learn more
- [Enable internet access for a VPC using an internet gateway](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Internet_Gateway.html)
- [Enable outbound IPv6 traffic using an egress-only internet gateway](https://docs.aws.amazon.com/vpc/latest/userguide/egress-only-internet-gateway.html)
- [Connect to the internet or other networks using NAT devices](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat.html)

## Access a corporate or home network
You can optionally connect your VPC to your own corporate data center using an IPsec AWS Site-to-Site VPN connection, making the AWS Cloud an extension of your data center.

A Site-to-Site VPN connection consists of two VPN tunnels between a virtual private gateway or transit gateway on the AWS side, and a customer gateway device located in your data center. A customer gateway device is a physical device or software appliance that you configure on your side of the Site-to-Site VPN connection.

### Learn more
- [AWS Site-to-Site VPN User Guide](https://docs.aws.amazon.com/vpn/latest/s2svpn/)
- [Amazon VPC Transit Gateways](https://docs.aws.amazon.com/vpc/latest/tgw/)

## Connect VPCs and networks
You can create a *VPC peering connection* between two VPCs that enables you to route traffic between them privately. Instances in either VPC can communicate with each other as if they are within the same network.

You can also create a *transit gateway* and use it to interconnect your VPCs and on-premises networks. The transit gateway acts as a *Regional virtual router* for traffic flowing between its attachments, which can include VPCs, VPN connections, AWS Direct Connect gateways, and transit gateway peering connections.

### Learn more
- [Amazon VPC Peering Guide](https://docs.aws.amazon.com/vpc/latest/peering/)
- [Amazon VPC Transit Gateways](https://docs.aws.amazon.com/vpc/latest/tgw/)

## AWS private global network
AWS provides a high-performance, and low-latency private global network that delivers a secure cloud computing environment to support your networking needs. **AWS Regions are connected to multiple Internet Service Providers (ISPs) as well as to a private global network backbone**, which provides improved network performance for cross-Region traffic sent by customers.

Packets that originate in the private global network with a destination in the private global network stay in the private global network and do not traverse the public internet. This is true whether the destination is a private IP address or a public IP address. For example, if EC2 instances in two VPCs communicate using public IP addresses, the traffic stays in the private global network. The destination can be in the same Availability Zone, a different Availability Zone in the same Region, or a different Region, except for the China Regions.

Network packet loss can be caused by a number of factors, including network flow collisions, lower level (Layer 2) errors, and other network failures. We engineer and operate our networks to minimize packet loss. We measure packet-loss rate (PLR) across the global backbone that connects the AWS Regions. We operate our backbone network to target a p99 of the hourly PLR of less than 0.0001%.

# Plan Your VPC
