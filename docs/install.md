
```bash
curl -sLS https://get.k3s.io | K3S_AGENT_TOKEN='sBiVyHxIIqYXtMCdLREWmsoCSUOYfPWD' K3S_TOKEN='mFLceYrueuxFigfVqGpPokPALhAGWHdl' INSTALL_K3S_EXEC='server --no-deploy traefik --flannel-backend=host-gw --secrets-encryption  --disable-network-policy=true --tls-san 192.168.192.40 --node-ip 192.168.192.40 --advertise-address 192.168.192.40' INSTALL_K3S_CHANNEL='stable' INSTALL_K3S_SKIP_START=true  sh -

/usr/local/bin/k3s-killall.sh
/usr/local/bin/k3s-uninstall.sh
```