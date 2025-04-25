Resources that you provision in your AMS Advanced environment automatically include the installation of an endpoint security (EPS) monitoring client. This process ensures that the AMS Advanced-managed resources are monitored and supported 24x7. In addition, AMS Advanced monitors all agent activity, and an incident is created if any security event is detected.

**Note**
Security incidents are handled as incidents; for more information, see [Incident response](https://docs.aws.amazon.com/managedservices/latest/userguide/sec-incident-response.html).

Endpoint security provides anti-malware protection, specifically, the following actions are supported:
- EC2 instances register with EPS
- EC2 instances deregister from EPS
- EC2 instances real-time anti-malware protection
- EPS agent-initiated heartbeat
- EPS restore quarantined file
- EPS event notification
- EPS reporting

AMS Advanced uses Trend Micro for endpoint security (EPS). These are the default EPS settings. To learn more about Trend Micro, see the [Trend Micro Deep Security Help Center](https://help.deepsecurity.trendmicro.com/aws/welcome.html?redirected=true); note that non-Amazon links may change without notice to us.

AMS Advanced Multi-Account Landing Zone (MALZ) default settings are described in the following sections; for non-default AMS multi-account landing zone EPS settings, see [AMS Advanced Multi-Account Landing Zone EPS non-default settings](https://docs.aws.amazon.com/managedservices/latest/userguide/security-mgmt.html#malz-eps-settings).

**Note**
You can bring your own EPS, see [AMS bring your own EPS](https://docs.aws.amazon.com/managedservices/latest/userguide/ams-byoeps.html).

# General EPS settings
https://docs.aws.amazon.com/managedservices/latest/onboardingguide/eps-defaults.html#general-eps-defaults

# Base policy
https://docs.aws.amazon.com/managedservices/latest/onboardingguide/eps-defaults.html#base-eps-policy

# Anti-malware
https://docs.aws.amazon.com/managedservices/latest/onboardingguide/eps-defaults.html#eps-anti-malware-defaults

