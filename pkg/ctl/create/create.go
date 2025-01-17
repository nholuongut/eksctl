package create

import (
	"github.com/spf13/cobra"

	"github.com/nholuongut/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `create` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("create", "Create resource(s)", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, createClusterCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, createNodeGroupCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, createIAMServiceAccountCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, createIAMIdentityMappingCmd)

	return verbCmd
}
