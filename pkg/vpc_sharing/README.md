Sharing VPCs is useful when network isolation between teams does not need to be strictly managed by the VPC owner, but the account level users and permissions must be. 

With [Shared VPC](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-sharing.html), multiple AWS accounts create their application resources (such as Amazon EC2 instances) in shared, centrally managed Amazon VPCs. 

In this model, the account that owns the VPC (owner) shares one or more subnets with other accounts (participants). After a subnet is shared, the participants can view, create, modify, and delete their application resources in the subnets shared with them. Participants cannot view, modify, or delete resources that belong to other participants or the VPC owner. 

Security between resources in shared VPCs is managed using security groups, network access control lists (NACLs), or through a firewall between the subnets.

VPC sharing benefits:
- Simplified design — no complexity around inter-VPC connectivity
- Fewer managed VPCs
- Segregation of duties between network teams and application owners
- Better IPv4 address utilization
- Lower costs — no data transfer charges between instances belonging to different accounts within the same Availability Zone.

**Note**
When you share a subnet with multiple accounts, your participants should have some level of cooperation since they're sharing IP space and network resources. If necessary, you can choose to share a different subnet for each participant account. One subnet per participant enables network ACL to provide network isolation in addition to security groups.

Most customer architectures will contain multiple VPCs, many of which will be shared with two or more accounts. Transit Gateway and VPC peering can be used to connect the shared VPCs. For example, suppose you have 10 applications. Each application requires its own AWS account. The apps can be categorized into two application portfolios (apps within the same portfolio have similar networking requirements, App 1–5 in ‘Marketing’ and App 6–10 in ‘Sales’).

You can have one VPC per application portfolio (two VPCs total), and the VPC is shared with the different application owner accounts within that portfolio. App owners deploy apps into their respective shared VPC (in this case, in the different subnets for network route segmentation and isolation using NACLs). The two shared VPCs are connected via the Transit Gateway. With this setup, you could go from having to connect 10 VPCs to just two, as seen in the following figure.

pic

**Note**
VPC sharing participants cannot create all AWS resources in a shared subnet. For more information, refer to the [Limitations](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-sharing.html#vpc-share-limitations) section in the VPC Sharing documentation.

For more information about the key considerations and best practices for VPC sharing, refer to the [VPC sharing: key considerations and best practices](https://aws.amazon.com/blogs/networking-and-content-delivery/vpc-sharing-key-considerations-and-best-practices/) blog post.