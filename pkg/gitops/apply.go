package gitops

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/nholuongut/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/nholuongut/eksctl/pkg/git"
	"github.com/nholuongut/eksctl/pkg/gitops/flux"
)

// Applier can set up a repo as a gitops repo with flux
type Applier struct {
	UserRepoPath     string
	ClusterConfig    *api.ClusterConfig
	UsersRepoOpts    git.Options
	QuickstartName   string
	FluxInstaller    *flux.Installer
	ProfileGenerator *Profile
	GitClient        *git.Client
}

// Run sets up gitops in a repository and a cluster and installs flux, helm, tiller and a quickstart into the cluster
func (g *Applier) Run(ctx context.Context) error {

	// Install Flux, Helm and Tiller. Clones the user's repo
	userInstructions, err := g.FluxInstaller.Run(context.Background())
	if err != nil {
		return err
	}

	// Clone user's repo to apply Quick Start profile
	options := git.CloneOptions{
		URL:       g.UsersRepoOpts.URL,
		Branch:    g.UsersRepoOpts.Branch,
		Bootstrap: true,
	}
	err = g.GitClient.CloneRepoInPath(g.UserRepoPath, options)
	if err != nil {
		return err
	}

	// Add quickstart components to user's repo. Clones the quickstart base repo
	err = g.ProfileGenerator.Generate(context.Background())
	if err != nil {
		return errors.Wrap(err, "error generating profile")
	}

	// Git add, commit and push component files
	if err = g.GitClient.Add("."); err != nil {
		return err
	}

	commitMsg := fmt.Sprintf("Add %s quickstart components", g.QuickstartName)
	if err = g.GitClient.Commit(commitMsg, g.UsersRepoOpts.User, g.UsersRepoOpts.Email); err != nil {
		return err
	}

	if err = g.GitClient.Push(); err != nil {
		return err
	}

	logger.Info(userInstructions)
	return nil
}
