NOTES:
HAProxy Kubernetes Ingress Controller has been successfully installed.

Controller image deployed is: "haproxytech/kubernetes-ingress:1.8.0".
Your controller is of a "Deployment" kind. Your controller service is running as a "LoadBalancer" type.
RBAC authorization is enabled.
Controller ingress.class is set to "haproxy" so make sure to use same annotation for
Ingress resource.

Service ports mapped are:
  - name: http
    containerPort: 80
    protocol: TCP
  - name: https
    containerPort: 443
    protocol: TCP
  - name: stat
    containerPort: 1024
    protocol: TCP

Node IP can be found with:
  $ kubectl --namespace haproxy-controller get nodes -o jsonpath="{.items[0].status.addresses[1].address}"

The following ingress resource routes traffic to pods that match the following:
  * service name: web
  * client's Host header: webdemo.com
  * path begins with /

  ---
  apiVersion: networking.k8s.io/v1
  kind: Ingress
  metadata:
    name: web-ingress
    namespace: default
    annotations:
      ingress.class: "haproxy"
  spec:
    rules:
    - host: webdemo.com
      http:
        paths:
        - path: /
          backend:
            serviceName: web
            servicePort: 80

In case that you are using multi-ingress controller environment, make sure to use ingress.class annotation and match it
with helm chart option controller.ingressClass.

For more examples and up to date documentation, please visit:
  * Helm chart documentation: https://github.com/haproxytech/helm-charts/tree/main/kubernetes-ingress
  * Controller documentation: https://www.haproxy.com/documentation/kubernetes/latest/
  * Annotation reference: https://github.com/haproxytech/kubernetes-ingress/tree/master/documentation
  * Image parameters reference: https://github.com/haproxytech/kubernetes-ingress/blob/master/documentation/controller.md
