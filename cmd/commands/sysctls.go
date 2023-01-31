package commands

import (
	"github.com/Shopify/kubeaudit/auditors/sysctls"
	"github.com/spf13/cobra"
)

var sysctlsCmd = &cobra.Command{
	Use:   "sysctls",
	Short: "Audit containers running with sysctls",
	Long: `This command determines which containers are running with sysctls enabled.

An ERROR result is generated when a container has added unsafe sysctls.

Example usage:
kubeaudit sysctls`,
	Run: runAudit(sysctls.New()),
}

func init() {
	RootCmd.AddCommand(sysctlsCmd)
}
