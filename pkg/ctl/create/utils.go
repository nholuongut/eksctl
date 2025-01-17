package create

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"

	api "github.com/nholuongut/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/nholuongut/eksctl/pkg/ctl/cmdutils"
	"github.com/nholuongut/eksctl/pkg/eks"
	"github.com/nholuongut/eksctl/pkg/ssh"
	"github.com/nholuongut/eksctl/pkg/utils/file"
)

func checkSubnetsGivenAsFlags(params *createClusterCmdParams) bool {
	return len(*params.subnets[api.SubnetTopologyPrivate])+len(*params.subnets[api.SubnetTopologyPublic]) != 0
}

func checkVersion(cmd *cmdutils.Cmd, ctl *eks.ClusterProvider, meta *api.ClusterMeta) error {
	switch meta.Version {
	case "auto":
		break
	case "":
		meta.Version = "auto"
	case "default":
		meta.Version = api.DefaultVersion
		logger.Info("will use default version (%s) for new nodegroup(s)", meta.Version)
	case "latest":
		meta.Version = api.LatestVersion
		logger.Info("will use latest version (%s) for new nodegroup(s)", meta.Version)
	default:
		if !isValidVersion(meta.Version) {
			if isDeprecatedVersion(meta.Version) {
				return fmt.Errorf("invalid version, %s is now deprecated, supported values: auto, default, latest, %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", meta.Version, strings.Join(api.SupportedVersions(), ", "))
			}
			return fmt.Errorf("invalid version %s, supported values: auto, default, latest, %s", meta.Version, strings.Join(api.SupportedVersions(), ", "))
		}
	}

	if v := ctl.ControlPlaneVersion(); v == "" {
		return fmt.Errorf("unable to get control plane version")
	} else if meta.Version == "auto" {
		meta.Version = v
		logger.Info("will use version %s for new nodegroup(s) based on control plane version", meta.Version)
	} else if meta.Version != v {
		hint := "--version=auto"
		if cmd.ClusterConfigFile != "" {
			hint = "metadata.version: auto"
		}
		logger.Warning("will use version %s for new nodegroup(s), while control plane version is %s; to automatically inherit the version use %q", meta.Version, v, hint)
	}

	return nil
}

func isValidVersion(version string) bool {
	for _, v := range api.SupportedVersions() {
		if version == v {
			return true
		}
	}
	return false
}

func isDeprecatedVersion(version string) bool {
	for _, v := range api.DeprecatedVersions() {
		if version == v {
			return true
		}
	}
	return false
}

// loadSSHKey loads the ssh public key specified in the NodeGroup. The key should be specified
// in only one way: by name (for a key existing in EC2), by path (for a key in a local file)
// or by its contents (in the config-file). It also assumes that if ssh is enabled (SSH.Allow
// == true) then one key was specified
func loadSSHKey(ng *api.NodeGroup, clusterName string, provider api.ClusterProvider) error {
	sshConfig := ng.SSH
	if sshConfig.Allow == nil || *sshConfig.Allow == false {
		return nil
	}

	switch {

	// Load Key by content
	case sshConfig.PublicKey != nil:
		keyName, err := ssh.LoadKeyByContent(sshConfig.PublicKey, clusterName, ng.Name, provider)
		if err != nil {
			return err
		}
		sshConfig.PublicKeyName = &keyName

	// Use key by name in EC2
	case sshConfig.PublicKeyName != nil && *sshConfig.PublicKeyName != "":
		if err := ssh.CheckKeyExistsInEC2(*sshConfig.PublicKeyName, provider); err != nil {
			return err
		}
		logger.Info("using EC2 key pair %q", *sshConfig.PublicKeyName)

	// Local ssh key file
	case file.Exists(*sshConfig.PublicKeyPath):
		keyName, err := ssh.LoadKeyFromFile(*sshConfig.PublicKeyPath, clusterName, ng.Name, provider)
		if err != nil {
			return err
		}
		sshConfig.PublicKeyName = &keyName

	// A keyPath, when specified as a flag, can mean a local key (checked above) or a key name in EC2
	default:
		err := ssh.CheckKeyExistsInEC2(*sshConfig.PublicKeyPath, provider)
		if err != nil {
			return err
		}
		sshConfig.PublicKeyName = sshConfig.PublicKeyPath
		sshConfig.PublicKeyPath = nil
		logger.Info("using EC2 key pair %q", *ng.SSH.PublicKeyName)
	}

	return nil
}
