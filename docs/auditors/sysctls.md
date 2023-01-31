# Sysctls Auditor (sysctls)

Finds Pods and Deployments using unsafe sysctls

## General Usage

```bash
kubeaudit sysctls [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```bash
$ kubeaudit sysctls -f "auditors/sysctls/fixtures/pod-unsafe.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod-unsafe
    namespace: pod-unsafe

--------------------------------------------

-- [error] UnsafeSysctlsUsed
   Message: Unsafe Sysctls are used. Please check the Kubernetes documentation (https://kubernetes.io/docs/tasks/administer-cluster/sysctl-cluster/#enabling-unsafe-sysctls) for allowed Sysctls
   Metadata:
      Sysctls: net.core.somaxconn,kernel.msgmax
```

## Explanation

Sysctl is a interface which allows an administrator to list and modify kernel parameters. Sysctls are grouped into safe and unsafe sysctls. By far, most of the sysctls are not considered safe. The following sysctls are supported in the safe set:

- kernel.shm_rmid_forced,
- net.ipv4.ip_local_port_range,
- net.ipv4.tcp_syncookies,
- net.ipv4.ping_group_range (since Kubernetes 1.18),
- net.ipv4.ip_unprivileged_port_start (since Kubernetes 1.22).

The listed safe sysctls can be used without impacting the baseline and restricted [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/).

This list is copied from the Kubernetes documentation an could be extended in the future.

The sysctls-directive has to be placed in the SecurityContext of the Pod or Deployment specification and has the following format:

```yaml
spec:
  securityContext:
    sysctls:
      - name: [name of the sysctl]
        value: [value]
```

Example of a resource which passes the `sysctls` audit:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod
  namespace: example
spec:
  containers:
    - name: httpbin
      image: docker.io/kennethreitz/httpbin
  securityContext:
    sysctls:
      - name: kernel.shm_rmid_forced
        value: '1'
      - name: net.ipv4.ip_local_port_range
        value: '32000 64128'
      - name: net.ipv4.ip_unprivileged_port_start
        value: '28256'
      - name: net.ipv4.tcp_syncookies
        value: '2'
      - name: net.ipv4.ping_group_range
        value: '100 100'
```

Example of a resource which doesn't pass the `sysctls` audit:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-unsafe
  namespace: example
spec:
  containers:
    - name: httpbin
      image: docker.io/kennethreitz/httpbin
  securityContext:
    sysctls:
      - name: kernel.shm_rmid_forced
        value: '1'
      - name: net.core.somaxconn
        value: '1024'
      - name: kernel.msgmax
        value: '65536'
```

More information about the usage of sysctls in Kubernetes can be found in the [official Kubernetes documentation](https://kubernetes.io/docs/tasks/administer-cluster/sysctl-cluster/). To get mor information about the sysctl command line tool, reach out to its [manpage](https://www.man7.org/linux/man-pages/man8/sysctl.8.html)

## Override Errors

Overrides are not currently supported for `sysctls`.
