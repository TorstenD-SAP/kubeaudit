# Setup test environment

Setup a Kubernetes cluster using K3d and allow unsafe sysctls (net.core.somaxconn, kernel.msg*). Setting up a cluster in that way works better for me because otherwise it is not possible to deploy the unsafe specifications successfully. For details on k3d visit [k3d.io](https://k3d.io).

```bash
k3d cluster create test --agents 1 --servers 1 --k3s-arg "--kubelet-arg=allowed-unsafe-sysctls=net.core.somaxconn,kernel.msg*@server:*;agents:*" --port 8081:80@loadbalancer --port 8444:443@loadbalancer --verbose
```
