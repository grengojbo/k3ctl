
```bash
kubectl patch ing/main-ingress --type=json \
  -p='[{"op": "replace", "path": "/spec/rules/0/http/paths/1/backend/serviceName", "value":"some-srvc"}, {"op": "replace", "path": "/spec/rules/1/http/paths/1/backend/serviceName", "value":"some-srvc"}]'
```

```bash
kubectl -n default patch ing/kuard --type=json -p='[{"op": "replace", "path": "/spec/ingressClassName", "value":"nginx"}]'
```

```bash
kubectl -n default patch ing/kuard --type=json -p='[{"op": "replace", "path": "/spec/ingressClassName", "value":"haproxy"}]'
```