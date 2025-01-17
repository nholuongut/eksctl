package get

import (
	"github.com/spf13/cobra"

	"github.com/nholuongut/eksctl/pkg/ctl/cmdutils"
)

type getCmdParams struct {
	chunkSize int
	output    string
}

// Command will create the `get` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("get", "Get resource(s)", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getClusterCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getNodeGroupCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getIAMServiceAccountCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getIAMIdentityMappingCmd)

	return verbCmd
}
