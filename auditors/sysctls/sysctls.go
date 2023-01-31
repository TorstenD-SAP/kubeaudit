package sysctls

import (
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
)

const Name = "sysctls"

const (
	//UnsafeSysctlsUsed occurs, when unsafe Sysctls are used
	UnsafeSysctlsUsed = "UnsafeSysctlsUsed"
	//Message to print if unsafe Sysctls are used
	UnsafeSysctlsMsg = "Unsafe Sysctls are used. Please check the Kubernetes documentation (https://kubernetes.io/docs/tasks/administer-cluster/sysctl-cluster/#enabling-unsafe-sysctls) for safe Sysctls"
	//Message to print if only safe Sysctls are use
	SafeSysctlsUsed = "Only safe Sysctls used according to the Kubernetes documentation"
	//Message to print if Sysctls isn't specified
	NoSysctls = "Sysctls are not specified in the SecurityContext"
)

var defaultAllowedSysctl = []string{
	"kernel.shm_rmid_forced",
	"net.ipv4.ip_local_port_range",
	"net.ipv4.ip_unprivileged_port_start",
	"net.ipv4.tcp_syncookies",
	"net.ipv4.ping_group_range",
}

// Sysctls implements Auditable
type Sysctls struct {
	allowedSysctls []string
}

func New() *Sysctls {
	return &Sysctls{allowedSysctls: defaultAllowedSysctl}
}

// Audit checks if Sysctls is used and if safe Sysctls are listed
func (sysctls *Sysctls) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {

	var auditResults []*kubeaudit.AuditResult

	spec := k8s.GetPodSpec(resource)
	if spec == nil {
		return auditResults, nil
	}

	if !containsSysctls(spec) {
		return auditResults, nil
	}

	unsafeSysctls := auditSysctls(spec, sysctls.allowedSysctls)
	if len(unsafeSysctls) == 0 {
		return auditResults, nil
	}

	auditResults = append(auditResults, &kubeaudit.AuditResult{
		Auditor:  Name,
		Rule:     UnsafeSysctlsUsed,
		Severity: kubeaudit.Error,
		Message:  UnsafeSysctlsMsg,
		Metadata: kubeaudit.Metadata{
			"Sysctls": strings.Join(unsafeSysctls, ","),
		},
	})
	return auditResults, nil
}

func containsSysctls(spec *k8s.PodSpecV1) bool {
	if spec.SecurityContext == nil {
		return false
	}
	return spec.SecurityContext.Sysctls != nil
}

func auditSysctls(spec *k8s.PodSpecV1, allowedSysctls []string) []string {

	var result = []string{}
	for _, element := range spec.SecurityContext.Sysctls {
		if !isValueInList(element.Name, allowedSysctls) {
			result = append(result, element.Name)
		}
	}
	return result
}

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
