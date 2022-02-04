


https://artifacthub.io/packages/helm/cert-manager/cert-manager

```bash
arkade install cert-manager --wait -v v1.7.1
~/.arkade/bin/helm [upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager --wait --values /var/folders/g8/7hd0pg1n3d34dmqz1jrfm4mm0000gp/T/charts/cert-manager/values.yaml --set installCRDs=true]
```

ingress-nginx
https://artifacthub.io/packages/helm/ingress-nginx/ingress-nginx

```bash
arkade install ingress-nginx --wait -n ingress-nginx
/.arkade/bin/helm [upgrade --install ingress-nginx ingress-nginx/ingress-nginx --namespace ingress-nginx --wait --values /var/folders/g8/7hd0pg1n3d34dmqz1jrfm4mm0000gp/T/charts/ingress-nginx/values.yaml
```

The ingress-nginx controller has been installed.
It may take a few minutes for the LoadBalancer IP to be available.
You can watch the status by running 'kubectl --namespace ingress-nginx get services -o wide -w ingress-nginx-controller'

An example Ingress that makes use of the controller:
  apiVersion: networking.k8s.io/v1
  kind: Ingress
  metadata:
    name: example
    namespace: foo
  spec:
    ingressClassName: nginx
    rules:
      - host: www.example.com
        http:
          paths:
            - backend:
                service:
                  name: exampleService
                  port:
                    number: 80
              path: /
    # This section is only required if TLS is to be enabled for the Ingress
    tls:
      - hosts:
        - www.example.com
        secretName: example-tls

If TLS is enabled for the Ingress, a Secret containing the certificate and key must also be provided:

  apiVersion: v1
  kind: Secret
  metadata:
    name: example-tls
    namespace: foo
  data:
    tls.crt: <base64 encoded cert>
    tls.key: <base64 encoded key>
  type: kubernetes.io/tls
=======================================================================
= ingress-nginx has been installed.                                   =
=======================================================================

# If you're using a local environment such as "minikube" or "KinD",
# then try the inlets operator with "arkade install inlets-operator"

# If you're using a managed Kubernetes service, then you'll find
# your LoadBalancer's IP under "EXTERNAL-IP" via:

kubectl get svc ingress-nginx-controller

# Find out more at:
# https://github.com/kubernetes/ingress-nginx/tree/master/charts/ingress-nginx

Thanks for using arkade!

