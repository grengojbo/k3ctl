
NAME: netdata
LAST DEPLOYED: Tue May 17 21:54:11 2022
NAMESPACE: monitoring
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
1. netdata will be available on http://netdata.k8s.local/, on the exposed port of your ingress controller

In a production environment, you 
 You can get that port via `kubectl get services`. e.g. in the following example, the http exposed port is 31737, the https one is 30069.
 The hostname netdata.k8s.local will need to be added to /etc/hosts, so that it resolves to the exposed IP. That IP depends on how your cluster is set up: 
	- When no load balancer is available (e.g. with minikube), you get the IP shown on `kubectl cluster-info`
	- In a production environment, the command `kubectl get services` will show the IP under the EXTERNAL-IP column

The port can be retrieved in both cases from `kubectl get services`

NAME                                         TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
exiled-tapir-nginx-ingress-controller        LoadBalancer   10.98.132.169    <pending>     80:31737/TCP,443:30069/TCP   11h