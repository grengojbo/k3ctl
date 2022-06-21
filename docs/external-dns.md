NAME: external-dns
LAST DEPLOYED: Tue Jun 21 13:10:01 2022
NAMESPACE: kube-system
STATUS: deployed
REVISION: 3
TEST SUITE: None
NOTES:
CHART NAME: external-dns
CHART VERSION: 6.5.6
APP VERSION: 0.12.0

** Please be patient while the chart is being deployed **

To verify that external-dns has started, run:

```bash
kubectl --namespace=kube-system get pods -l "app.kubernetes.io/name=external-dns,app.kubernetes.io/instance=external-dns"
```