package sysctls

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditSysctls(t *testing.T) {
	cases := []struct {
		file          string
		expectedError []string
		testLocalMode bool
	}{
		{"pod-wo.yml", nil, false},
		{"pod-safe.yml", nil, false},
		{"pod-unsafe.yml", []string{UnsafeSysctlsUsed}, false},
		{"deployment-wo.yml", nil, false},
		{"deployment-safe.yml", nil, false},
		{"deployment-unsafe.yml", []string{UnsafeSysctlsUsed}, false},
		{"pod-wo.yml", nil, true},
		{"pod-safe.yml", nil, true},
		{"pod-unsafe.yml", []string{UnsafeSysctlsUsed}, true},
		{"deployment-wo.yml", nil, true},
		{"deployment-safe.yml", nil, true},
		{"deployment-unsafe.yml", []string{UnsafeSysctlsUsed}, true},
	}

	for _, tc := range cases {
		// This line is needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			test.AuditManifest(t, fixtureDir, tc.file, New(), tc.expectedError)
			if tc.testLocalMode {
				test.AuditLocal(t, fixtureDir, tc.file, New(), strings.Split(tc.file, ".")[0], tc.expectedError)
			}
		})
	}
}
