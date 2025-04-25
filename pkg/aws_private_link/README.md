[AWS PrivateLink](https://aws.amazon.com/privatelink/) provides private connectivity between VPCs, AWS services, and your on-premises networks without exposing your traffic to the public internet. 

Interface VPC endpoints, powered by AWS PrivateLink, make it easy to connect to AWS and other services across diﬀerent accounts and VPCs to signiﬁcantly simplify your network architecture. This allows customers who may want to privately expose a service/application residing in one VPC (service provider) to other VPCs (consumer) within an AWS Region in a way that only consumer VPCs initiate connections to the service provider VPC. **An example of this is the ability for your private applications to access service provider APIs**.

To use AWS PrivateLink, create a **Network Load Balancer** for your application in your VPC, and create a VPC endpoint service configuration pointing to that load balancer. A service consumer then creates an **interface endpoint** to your service. This creates an elastic network interface (ENI) in the consumer subnet with a private IP address that serves as an entry point for traffic destined for the service. The consumer and service are not required to be in the same VPC. **If the VPC is different, the consumer and service provider VPCs can have overlapping IP address ranges**. In addition to creating the interface VPC endpoint to access services in other VPCs, you can create interface VPC endpoints to privately access [supported AWS services](https://docs.aws.amazon.com/vpc/latest/userguide/vpce-interface.html) through AWS PrivateLink, as shown in the following figure.

pic

With Application Load Balancer (ALB) as target of NLB, you can now combine ALB advanced routing capabilities with AWS PrivateLink. Refer to [Application Load Balancer-type Target Group for Network Load Balancer](https://aws.amazon.com/blogs/networking-and-content-delivery/application-load-balancer-type-target-group-for-network-load-balancer/) for reference architectures and detailed configuration.

The choice between Transit Gateway, VPC peering, and AWS PrivateLink is dependent on connectivity.
- **AWS PrivateLink** — Use AWS PrivateLink when you have a client/server set up where you want to allow one or more consumer VPCs unidirectional access to a specific service or set of instances in the service provider VPC or certain AWS services. Only the clients with access in the consumer VPC can initiate a connection to the service in the service provider VPC or AWS service. This is also a good option when client and servers in the two VPCs have overlapping IP addresses because **AWS PrivateLink uses ENIs** within the client VPC in a manner that ensures that are no IP conflicts with the service provider. You can access AWS PrivateLink endpoints over VPC peering, VPN, Transit Gateway, Cloud WAN, and AWS Direct Connect.

- **VPC peering and Transit Gateway** — Use VPC peering and Transit Gateway when you want to enable layer-3 IP connectivity between VPCs.

Your architecture will contain a mix of these technologies in order to fulfill different use cases. All of these services can be combined and operated with each other. For example, 
- AWS PrivateLink handling API style client-server connectivity, 
- VPC peering for handling direct connectivity requirements where placement groups may still be desired within the Region or inter-Region connectivity is needed, 
- And Transit Gateway to simplify connectivity of VPCs at scale as well as edge consolidation for hybrid connectivity.