# How VPC peering connections work
https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html

The following steps describe the VPC peering process:
- The owner of the requester VPC sends a request to the owner of the accepter VPC to create the VPC peering connection. The accepter VPC can be owned by you, or another AWS account, and cannot have a CIDR block that overlaps with the CIDR block of the requester VPC.
- The owner of the accepter VPC accepts the VPC peering connection request to activate the VPC peering connection.
- To enable the flow of traffic between the VPCs using private IP addresses, the owner of each VPC in the VPC peering connection must manually add a route to one or more of their VPC route tables that points to the IP address range of the other VPC (the peer VPC).
- If required, update the security group rules that are associated with your EC2 instance to ensure that traffic to and from the peer VPC is not restricted. If both VPCs are in the same Region, you can reference a security group from the peer VPC as a source or destination for inbound or outbound rules in your security group.
- With the default VPC peering connection options, if EC2 instances on either side of a VPC peering connection address each other using a public DNS hostname, the hostname resolves to the public IP address of the EC2 instance. To change this behavior, enable DNS hostname resolution for your VPC connection. After enabling DNS hostname resolution, if EC2 instances on either side of the VPC peering connection address each other using a public DNS hostname, the hostname resolves to the private IP address of the EC2 instance.

For more information, see [VPC peering connections](https://docs.aws.amazon.com/vpc/latest/peering/working-with-vpc-peering.html).

## VPC peering connection lifecycle
A VPC peering connection goes through various stages starting from when the request is initiated. At each stage, there may be actions that you can take, and at the end of its lifecycle, the VPC peering connection remains visible in the Amazon VPC console and API or command line output for a period of time.

pic

- **Initiating-request**: A request for a VPC peering connection has been initiated. At this stage, the peering connection can fail, or can go to ```pending-acceptance```.
- **Failed**: The request for the VPC peering connection has failed. While in this state, it cannot be accepted, rejected, or deleted. The failed VPC peering connection remains visible to the requester for 2 hours.
- **Pending-acceptance**: The VPC peering connection request is awaiting acceptance from the owner of the accepter VPC. During this state, the owner of the requester VPC can delete the request, and the owner of the accepter VPC can accept or reject the request. If no action is taken on the request, it expires after 7 days.
- **Expired**: The VPC peering connection request has expired, and no action can be taken on it by either VPC owner. The expired VPC peering connection remains visible to both VPC owners for 2 days
- **Rejected**: The owner of the accepter VPC has rejected a ```pending-acceptance``` VPC peering connection request. While in this state, the request cannot be accepted. The rejected VPC peering connection remains visible to the owner of the requester VPC for 2 days, and visible to the owner of the accepter VPC for 2 hours. If the request was created within the same AWS account, the rejected request remains visible for 2 hours.
- **Provisioning**: The VPC peering connection request has been accepted, and will soon be in the ```active``` state.
- **Active**: The VPC peering connection is active, and traffic can flow between the VPCs (provided that your security groups and route tables allow the flow of traffic). While in this state, either of the VPC owners can delete the VPC peering connection, but cannot reject it.

**Note**
If an event in a Region in which a VPC resides prevents the flow of traffic, the status of the VPC peering connection remains ```Active```.

- **Deleting**: Applies to an inter-Region VPC peering connection that is in the process of being deleted. The owner of either VPC has submitted a request to delete an ```active``` VPC peering connection, or the owner of the requester VPC has submitted a request to delete a ```pending-acceptance``` VPC peering connection request.
- **Deleted**: An ```active``` VPC peering connection has been deleted by either of the VPC owners, or a ```pending-acceptance``` VPC peering connection request has been deleted by the owner of the requester VPC. While in this state, the VPC peering connection cannot be accepted or rejected. The VPC peering connection remains visible to the party that deleted it for 2 hours, and visible to the other party for 2 days. If the VPC peering connection was created within the same AWS account, the deleted request remains visible for 2 hours.

## Multiple VPC peering connections
A VPC peering connection is a one to one relationship between two VPCs. You can create multiple VPC peering connections for each VPC that you own, but transitive peering relationships are not supported. You do not have any peering relationship with VPCs that your VPC is not directly peered with.

The following diagram is an example of one VPC peered to two different VPCs. There are two VPC peering connections: VPC A is peered with both VPC B and VPC C. VPC B and VPC C are not peered, and you cannot use VPC A as a transit point for peering between VPC B and VPC C. If you want to enable routing of traffic between VPC B and VPC C, you must create a unique VPC peering connection between them.

pic

## VPC peering limitations
Consider the following limitations for VPC peering connections. In some cases, you can use a transit gateway attachment instead of a VPC peering connection. For more information, see [Example transit gateway scenarios](https://docs.aws.amazon.com/vpc/latest/tgw/how-transit-gateways-work.html#TGW_Scenarios) in Amazon VPC Transit Gateways.

### Connections
- There is a quota on the number of active and pending VPC peering connections per VPC. For more information, see [VPC peering connection quotas for an account](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-connection-quotas.html).
- You cannot have more than one VPC peering connection between two VPCs at the same time.
- Any tags that you create for your VPC peering connection are only applied in the account or Region in which you create them.
- You cannot connect to or query the Amazon DNS server in a peer VPC.
- If the IPv4 CIDR block of a VPC in a VPC peering connection falls outside of the private IPv4 address ranges specified by [RFC 1918](http://www.faqs.org/rfcs/rfc1918.html), private DNS hostnames for that VPC cannot be resolved to private IP addresses. To resolve private DNS hostnames to private IP addresses, you can enable DNS resolution support for the VPC peering connection. For more information, see [Enable DNS resolution for a VPC peering connection](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-dns.html).
- You can enable resources on either side of a VPC peering connection to communicate over IPv6. You must associate an IPv6 CIDR block with each VPC, enable the instances in the VPCs for IPv6 communication, and route IPv6 traffic intended for the peer VPC to the VPC peering connection.
- Unicast reverse path forwarding in VPC peering connections is not supported. For more information, see [Routing for response traffic](https://docs.aws.amazon.com/vpc/latest/peering/peering-configurations-partial-access.html#peering-incorrect-response-routing).

### Overlapping CIDR blocks
- You cannot create a VPC peering connection between VPCs that have matching or overlapping IPv4 or IPv6 CIDR blocks.
- If you have multiple IPv4 CIDR blocks, you can't create a VPC peering connection if any of the CIDR blocks overlap, even if you intend to use only the non-overlapping CIDR blocks or only IPv6 CIDR blocks.

### Transitive peering
- VPC peering does not support transitive peering relationships. For example, if there are VPC peering connections between VPC A and VPC B, and between VPC A and VPC C, you can't route traffic from VPC B to VPC C through VPC A. To route traffic between VPC B and VPC C, you must create a VPC peering connection between them. For more information, see Three VPCs peered together.

### Edge to edge routing through a gateway or private connection
- If VPC A has an internet gateway, resources in VPC B can't use the internet gateway in VPC A to access the internet.
- If VPC A has an NAT device that provides internet access to subnets in VPC A, resources in VPC B can't use the NAT device in VPC A to access the internet.
- If VPC A has a VPN connection to a corporate network, resources in VPC B can't use the VPN connection to communicate with the corporate network.
- If VPC A has an AWS Direct Connect connection to a corporate network, resources in VPC B can't use the AWS Direct Connect connection to communicate with the corporate network.
- If VPC A has a gateway endpoint that provides connectivity to Amazon S3 to private subnets in VPC A, resources in VPC B can't use the gateway endpoint to access Amazon S3.

### Inter-Region VPC peering connections
- For jumbo frames, the Maximum Transmission Unit (MTU) between VPC peering connections within the same Region is 9001 bytes. The MTU for inter-region VPC peering connections is 8500 bytes. For more information about jumbo frames, see [Jumbo frames (9001 MTU)](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/network_mtu.html#jumbo_frame_instances) in the Amazon EC2 User Guide.
- You must enable DNS resolution support for the VPC peering connection to resolve private DNS hostnames of the peered VPC to private IP addresses, even if the IPv4 CIDR for the VPC falls into the private IPv4 address ranges specified by RFC 1918.

### Shared VPCs and subnets
Only VPC owners can work with (describe, create, accept, reject, modify, or delete) peering connections. Participants cannot work with peering connections. For more information see, [Share your VPC with other accounts](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-sharing.html) in the Amazon VPC User Guide.

# VPC peering connections
https://docs.aws.amazon.com/vpc/latest/peering/working-with-vpc-peering.html

VPC peering enables you to connect two VPCs in the same or different AWS Regions. This enables instances in one VPC to communicate with instances in the other VPC as if they were all part of the same network.

**VPC peering creates a direct network route between the two VPCs** using private IPv4 addresses or IPv6 addresses. **Traffic sent between the connected VPCs does not traverse the internet, a VPN connection, or an AWS Direct Connect connection**. This makes VPC peering a secure way to share resources, such as databases or web servers, across VPC boundaries.

To establish a VPC peering connection, you create a peering connection request from one VPC and the owner of the other VPC accepts the request. After the connection is established, you can update your route tables to route traffic between the VPCs. This allows instances in one VPC to access resources in the other VPC.

VPC peering is an important tool for building multi-VPC architectures and sharing resources across organizational boundaries in AWS. It provides a simple, low-latency way to connect VPCs without the complexity of configuring a VPN or other networking service.

Use the following procedures to create and work with VPC peering connections.

## Create a VPC peering connection
To create a VPC peering connection, first create a request to peer with another VPC. To activate the request, the owner of the accepter VPC must accept the request. The following peering connections are supported:
- Between VPCs in the same account and Region
- Between VPCs in the same account and different Regions
- Between VPCs in different accounts and the same Region
- Between VPCs in different accounts and Regions

For an inter-Region VPC peering connection, the request must be made from the Region of the requester VPC, and the request must be accepted from the Region of the accepter VPC. For more information, see [Accept or reject a VPC peering connection](https://docs.aws.amazon.com/vpc/latest/peering/accept-vpc-peering-connection.html).

### Prerequisites
- Review the [limitations](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-limitations) for VPC peering connections.
- Ensure that the VPCs do not have overlapping IPv4 CIDR blocks. If they overlap, the status of the VPC peering connection immediately goes to failed. This limitation applies even if the VPCs have unique IPv6 CIDR blocks.

### Create a peering connection using the console
Use the following procedure to create a VPC peering connection.

#### To create a peering connection using the console
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- In the navigation pane, choose **Peering connections**.
- Choose Create peering connection.
- (Optional) For **Name**, specify a name the VPC peering connection. This creates a tag with a key of Name and the value that you specify.
- For **VPC ID (Requester)**, select a VPC from the current account.
- Under Select another VPC to peer with, do the following:
  - For **Account**, to peer with a VPC in another account, choose **Another account** and enter the account ID. Otherwise, keep **My account**.
  - For **Region**, to peer with a VPC in another Region, choose **Another Region** and choose the Region. Otherwise, keep **This Region**.
  - For **VPC ID (Accepter)**, select a VPC from the specified account and Region.
- (Optional) To add a tag, choose **Add new tag** and enter the tag key and tag value.
- Choose **Create peering connection**.
- The owner of the accepter account must accept the peering connection. For more information, see [Accept or reject a VPC peering connection](https://docs.aws.amazon.com/vpc/latest/peering/accept-vpc-peering-connection.html).
- Update the route tables for both VPCs to enable communication between them. For more information, see [Update your route tables for a VPC peering connection](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-routing.html).

### Create a peering connection using the command line
You can create a VPC peering connection using the following commands:
- [create-vpc-peering-connection](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/create-vpc-peering-connection.html) (AWS CLI)
- [New-EC2VpcPeeringConnection](https://docs.aws.amazon.com/powershell/latest/reference/items/New-EC2VpcPeeringConnection.html) (AWS Tools for Windows PowerShell)

## Accept or Reject a VPC Peering Connection
A VPC peering connection that's in the ```pending-acceptance``` state must be accepted by the owner of the accepter VPC to be activated. For more information about the ```Deleted``` peering connection status, see [VPC peering connection lifecycle](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-lifecycle). You can't accept a VPC peering connection request that you sent to another AWS account. To create a VPC peering connection between VPCs in the same AWS account, you can both create and accept the request yourself.

You can reject any VPC peering connection request that you've received that's in the ```pending-acceptance``` state. You should only accept VPC peering connections from AWS accounts that you know and trust; you can reject any unwanted requests. For more information about the ```Rejected``` peering connection status, see [VPC peering connection lifecycle](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-lifecycle).

**Important**
Do not accept VPC peering connections from unknown AWS accounts. A malicious user may have sent you a VPC peering connection request to gain unauthorized network access to your VPC. This is known as peer phishing. You can safely reject unwanted VPC peering connection requests without any risk of the requester gaining access to any information about your AWS account or your VPC. For more information, see Accept or reject a VPC peering connection. You can also ignore the request and let it expire; by default, requests expire after 7 days.

### To accept or reject a peering connection using the console
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- Use the Region selector to choose the Region of the accepter VPC.
- In the navigation pane, choose **Peering connections**.
- To reject a peering connection, select the VPC peering connection, and choose **Actions, Reject request**. When prompted for confirmation, choose **Reject request**.
- To accept a peering connection, select the pending VPC peering connection (the status is pending-acceptance), and choose **Actions, Accept request**. For more information about peering connection lifecycle statuses, see [VPC peering connection lifecycle](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-lifecycle). If there is no pending VPC peering connection, verify that you selected the Region of the accepter VPC.
- When prompted for confirmation, choose **Accept request**.
- Choose **Modify my route tables now** to add a route to the VPC route table so that you can send and receive traffic across the peering connection. For more information, see [Update your route tables for a VPC peering connection](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-routing.html).

### To accept a peering connection using the command line
- [accept-vpc-peering-connection](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/accept-vpc-peering-connection.html) (AWS CLI)
- [Approve-EC2VpcPeeringConnection](https://docs.aws.amazon.com/powershell/latest/reference/items/Approve-EC2VpcPeeringConnection.html) (AWS Tools for Windows PowerShell)

### To reject a peering connection using the command line
- [reject-vpc-peering-connection](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/reject-vpc-peering-connection.html) (AWS CLI)
- [Deny-EC2VpcPeeringConnection](https://docs.aws.amazon.com/powershell/latest/reference/items/Deny-EC2VpcPeeringConnection.html) (AWS Tools for Windows PowerShell)

## Update your route tables for a VPC peering connection
To enable private IPv4 traffic between instances in peered VPCs, you must add a route to the route tables associated with the subnets for both instances. The route destination is the CIDR block (or portion of the CIDR block) of the peer VPC and the target is the ID of the VPC peering connection. For more information, see [Configure route tables](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html) in the Amazon VPC User Guide.

The following is an example of the route tables that enables communication between instances in two peered VPCs, VPC A and VPC B. Each table has a local route and a route that sends traffic for the peer VPC to the VPC peering connection.

| Route table | Destination | Target       |
| ----------- | ----------- | ------------ |
| VPC A       | VPC A CIDR  | Local        |
|             | VPC B CIDR  | pcx-11112222 |
| VPC B       | VPC B CIDR  | Local        |
|             | VPC A CIDR  | pcx-11112222 |

Similarly, if the VPCs in the VPC peering connection have associated IPv6 CIDR blocks, you can add routes that enable communication with the peer VPC over IPv6.

For more information about supported route table configurations for VPC peering connections, see [Common VPC peering connection configurations](https://docs.aws.amazon.com/vpc/latest/peering/peering-configurations.html).

**Considerations**
- If you have a VPC peered with multiple VPCs that have overlapping or matching IPv4 CIDR blocks, ensure that your route tables are configured to avoid sending response traffic from your VPC to the incorrect VPC. AWS currently does not support unicast reverse path forwarding in VPC peering connections that checks the source IP of packets and routes reply packets back to the source. For more information, see [Routing for response traffic](https://docs.aws.amazon.com/vpc/latest/peering/peering-configurations-partial-access.html#peering-incorrect-response-routing).
- Your account has a [quota](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html) on the number of entries you can add per route table. If the number of VPC peering connections in your VPC exceeds the route table entry quota for a single route table, consider using multiple subnets that are each associated with a custom route table.
- You can add a route for a VPC peering connection that's in the ```pending-acceptance``` state. However, the route has a state of ```blackhole```, and has no effect until the VPC peering connection is in the ```active``` state.

### To add an IPv4 route for a VPC peering connection
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- In the navigation pane, choose **Route tables**.
- Select the check box next to the route table that's associated with the subnet in which your instance resides. If you do not have a route table explicitly associated with that subnet, the main route table for the VPC is implicitly associated with the subnet.
- Choose **Actions, Edit routes**.
- Choose **Add route**.
- For **Destination**, enter the IPv4 address range to which the network traffic in the VPC peering connection must be directed. You can specify the entire IPv4 CIDR block of the peer VPC, a specific range, or an individual IPv4 address, such as the IP address of the instance with which to communicate. For example, if the CIDR block of the peer VPC is ```10.0.0.0/16```, you can specify a portion ```10.0.0.0/24```, or a specific IP address ```10.0.0.7/32```.
- For **Target**, select the VPC peering connection.
- Choose **Save changes**.

The owner of the peer VPC must also complete these steps to add a route to direct traffic back to your VPC through the VPC peering connection.

If you have resources in different AWS Regions that use IPv6 addresses, you can create an inter-Region peering connection. You can then add an IPv6 route for communication between the resources.

### To add an IPv6 route for a VPC peering connection
- Open the Amazon VPC console at https://console.aws.amazon.com/vpc/.
- In the navigation pane, choose **Route tables**.
- Select the check box next to the route table that's associated with the subnet in which your instance resides.

**Note**
If you do not have a route table associated with that subnet, select the main route table for the VPC, as the subnet then uses this route table by default.

- Choose **Actions, Edit routes**.
- Choose **Add route**.
- For **Destination**, enter the IPv6 address range for the peer VPC. You can specify the entire IPv6 CIDR block of the peer VPC, a specific range, or an individual IPv6 address. For example, if the CIDR block of the peer VPC is ```2001:db8:1234:1a00::/56```, you can specify a portion ```2001:db8:1234:1a00::/64```, or a specific IP address ```2001:db8:1234:1a00::123/128```.
- For **Target**, select the VPC peering connection.
- Choose **Save changes**.

For more information, see [Route tables](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html) in the Amazon VPC User Guide.

### To add or replace a route using the command line
- [create-route](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/create-route.html) and [replace-route](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/replace-route.html)(AWS CLI)
- [New-EC2Route](https://docs.aws.amazon.com/powershell/latest/reference/items/New-EC2Route.html) and [Set-EC2Route](https://docs.aws.amazon.com/powershell/latest/reference/items/Set-EC2Route.html)(AWS Tools for Windows PowerShell)

# VPC to VPC connectivity
Customers can use two different VPC connectivity patterns to set up multi-VPC environments: many to many, or hub and spoke. In the many-to-many approach, the traffic between each VPC is managed individually between each VPC. In the hub-and-spoke model, all inter-VPC traffic flows through a central resource, which routes traffic based on established rules.

# VPC Peering
The first way to connect two VPCs is to use VPC peering. In this setup, a connection enables full bidirectional connectivity between the VPCs. This peering connection is used to route traffic between the VPCs. **VPCs in different accounts and AWS Regions can also be peered together**. 

**All data transfer over a VPC peering connection that stays within an Availability Zone is free**. All data transfer over a VPC peering connection that crosses Availability Zones is charged at the standard in-region data transfer rates. If the VPCs are peered across Regions, standard inter-Region data transfer charges will apply.

VPC peering is point-to-point connectivity, and it does not support [transitive routing](https://docs.aws.amazon.com/vpc/latest/peering/invalid-peering-configurations.html#transitive-peering). For example, if you have a [VPC peering](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-peering.html) connection between VPC A and VPC B and between VPC A and VPC C, an instance in VPC B cannot transit through VPC A to reach VPC C. To route packets between VPC B and VPC C, you are required to create a direct VPC peering connection.

At scale, when you have tens or hundreds of VPCs, interconnecting them with peering can result in a mesh of hundreds or thousands of peering connections. A large number of connections can be difficult to manage and scale. For example, if you have 100 VPCs and you want to setup a full mesh peering between them, it will take 4,950 peering connections [n(n-1)/2] where n is the total number of VPCs. **There is a maximum limit of 125 active peering connections per VPC**.

![image](https://github.com/user-attachments/assets/052738ab-e755-46e4-bc4c-99832c1f74c1)

If you are using VPC peering, on-premises connectivity (VPN and/or Direct Connect) must be made to each VPC. Resources in a VPC cannot reach on-premises using the hybrid connectivity of a peered VPC, as shown in the preceding figure.

VPC peering is best used when resources in one VPC must communicate with resources in another VPC, the environment of both VPCs is controlled and secured, and the number of VPCs to be connected is less than 10 (to allow for the individual management of each connection). VPC peering offers the lowest overall cost and highest aggregate performance when compared to other options for inter-VPC connectivity.
