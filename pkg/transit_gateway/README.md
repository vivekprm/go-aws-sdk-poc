https://aws.amazon.com/transit-gateway/

AWS Transit Gateway provides a hub and spoke design for connecting VPCs and on-premises networks as a fully managed service without requiring you to provision third-party virtual appliances. No VPN overlay is required, and AWS manages high availability and scalability.

Transit Gateway enables customers to connect thousands of VPCs. You can attach all your hybrid connectivity (VPN and Direct Connect connections) to a single gateway, consolidating and controlling your organization's entire AWS routing configuration in one place (refer to the following figure). Transit Gateway controls how traffic is routed among all the connected spoke networks using route tables. This hub-and-spoke model simplifies management and reduces operational costs because VPCs only connect to the Transit Gateway instance to gain access to the connected networks.

pic

**Transit Gateway is a Regional resource** and can connect thousands of VPCs within the same AWS Region. You can connect multiple gateways over a single Direct Connect connection for hybrid connectivity. Typically, you can use just one Transit Gateway instance connecting all your VPC instances in a given Region, and use **Transit Gateway routing tables** to isolate them wherever needed. 

Note that you do not need additional transit gateways for high availability, because **transit gateways are highly available by design**; for redundancy, use a single gateway in each Region. However, **there is a valid case for creating multiple gateways to limit misconﬁguration blast radius, segregate control plane operations, and administrative ease-of-use**.

With Transit Gateway peering, customers can peer their Transit Gateway instances within same or multiple Regions and route traffic between them. It uses the same underlying infrastructure as VPC peering, and is therefore encrypted. For more information, refer to [Building a global network using AWS Transit Gateway Inter-Region peering](https://aws.amazon.com/blogs/networking-and-content-delivery/building-a-global-network-using-aws-transit-gateway-inter-region-peering/) and [AWS Transit Gateway now supports Intra-Region Peering](https://aws.amazon.com/blogs/networking-and-content-delivery/aws-transit-gateway-now-supports-intra-region-peering/).

Place your organization’s Transit Gateway instance in its **Network Services account**. This enables centralized management by network engineers who manage the Network services account. **Use AWS Resource Access Manager (RAM) to share a Transit Gateway instance for connecting VPCs across multiple accounts in your AWS Organization within the same Region**. AWS RAM enables you to easily and securely share AWS resources with any AWS account, or within your AWS Organization. For more information, refer to the [Automating AWS Transit Gateway attachments to a transit gateway in a central account](https://aws.amazon.com/blogs/networking-and-content-delivery/automating-aws-transit-gateway-attachments-to-a-transit-gateway-in-a-central-account/) blog post.

Transit Gateway also allows you to establish connectivity between SD-WAN infrastructure and AWS using **Transit Gateway Connect**. Use a Transit Gateway Connect attachment with Border Gateway Protocol (BGP) for dynamic routing and **Generic Routing Encapsulation (GRE) tunnel** protocol for high performance, delivering up to 20 Gbps total bandwidth per Connect attachment (up to four Transit Gateway Connect peers per Connect attachment). By using Transit Gateway Connect, you can integrate both on-premises SD-WAN infrastructure or SD-WAN appliances running in the cloud through a VPC attachment or AWS Direct Connect attachment as the underlying transport layer. Refer to [Simplify SD-WAN connectivity with AWS Transit Gateway Connect](https://aws.amazon.com/blogs/networking-and-content-delivery/simplify-sd-wan-connectivity-with-aws-transit-gateway-connect/) for reference architectures and detailed configuration.

# Usecases
- Deliver applications around the world
  - Build, deploy, and manage applications across thousands of Amazon VPCs without having to manage peering connections or update routing tables.
- Rapidly move to global scale
  - Share VPCs, Domain Name System (DNS), Microsoft Active Directory, and IPS/IDS across Regions with inter-Region peering.
- Smoothly respond to spikes in demand
  - Quickly add Amazon VPCs, AWS accounts, virtual private networking (VPN) capacity, or AWS Direct Connect gateways to meet unexpected demand.
- Host multicast applications on AWS
  - Host multicast applications that scale based on demand, without the need to buy and maintain custom hardware.