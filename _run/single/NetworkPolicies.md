K8s Network Policies
--------------------

#### `akash-services` Namespace

Declared to house all of the Akash services[Node + Provider] 

#### Calico network policy

Make targets for `_run/single` and `_run/kube` scripts/YAML to install `calico` as the network manager for the KinD test environments.





### Testing Connection from workload Namespace
```
kubectl run --namespace=ndvr8b3np4qpr887arpdrbbgckfusji5k6rvnkhcreg4e access --rm -ti --image busybox /bin/sh

If you don't see a command prompt, try pressing enter.

/ # ping 192.168.51.74
PING 192.168.51.74 (192.168.51.74): 56 data bytes
^C
--- 192.168.51.74 ping statistics ---
27 packets transmitted, 0 packets received, 100% packet loss
/ # wget 192.168.51.74:26657/
Connecting to 192.168.51.74:26657 (192.168.51.74:26657)
```

### Testing connection from `akash-services` namespace

This should work because connections are allowed to ingress between namespace'd containers.

```
$kubectl run --namespace=akash-services access --rm -ti --image busybox /bin/sh
If you don't see a command prompt, try pressing enter.
/ # wget 192.168.51.76:8080/
Connecting to 192.168.51.76:8080 (192.168.51.76:8080)
wget: server returned error: HTTP/1.1 404 Not Found
/ # wget akash-provider:8080/
Connecting to akash-provider:8080 (10.98.95.196:8080)
wget: server returned error: HTTP/1.1 404 Not Found
```
HTTP connection was able to be established with the akash-provider container.

