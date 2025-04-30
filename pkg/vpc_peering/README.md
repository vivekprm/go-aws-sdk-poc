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
