package manager

import (
	"fmt"

	api "github.com/nholuongut/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/nholuongut/eksctl/pkg/iam/oidc"
	"github.com/nholuongut/eksctl/pkg/kubernetes"
)

// NewTasksToCreateClusterWithNodeGroups defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works
func (c *StackCollection) NewTasksToCreateClusterWithNodeGroups(nodeGroups []*api.NodeGroup) *TaskTree {
	tasks := &TaskTree{Parallel: false}

	tasks.Append(
		&taskWithoutParams{
			info: fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
			call: c.createClusterTask,
		},
	)

	nodeGroupTasks := c.NewTasksToCreateNodeGroups(nodeGroups)
	if nodeGroupTasks.Len() > 0 {
		nodeGroupTasks.IsSubTask = true
		tasks.Append(nodeGroupTasks)
	}

	return tasks
}

// NewTasksToCreateNodeGroups defines tasks required to create all of the nodegroups
func (c *StackCollection) NewTasksToCreateNodeGroups(nodeGroups []*api.NodeGroup) *TaskTree {
	tasks := &TaskTree{Parallel: true}

	for _, ng := range nodeGroups {
		tasks.Append(&taskWithNodeGroupSpec{
			info:      fmt.Sprintf("create nodegroup %q", ng.NameString()),
			nodeGroup: ng,
			call:      c.createNodeGroupTask,
		})
		// TODO: move authconfigmap tasks here using kubernetesTask and kubernetes.CallbackClientSet
	}

	return tasks
}

// NewTasksToCreateIAMServiceAccounts defines tasks required to create all of the IAM ServiceAccounts
func (c *StackCollection) NewTasksToCreateIAMServiceAccounts(serviceAccounts []*api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) *TaskTree {
	tasks := &TaskTree{Parallel: true}

	for i := range serviceAccounts {
		sa := serviceAccounts[i]
		saTasks := &TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}

		saTasks.Append(&taskWithClusterIAMServiceAccountSpec{
			info:           fmt.Sprintf("create IAM role for serviceaccount %q", sa.NameString()),
			serviceAccount: sa,
			oidc:           oidc,
			call:           c.createIAMServiceAccountTask,
		})

		saTasks.Append(&kubernetesTask{
			info:       fmt.Sprintf("create serviceaccount %q", sa.NameString()),
			kubernetes: clientSetGetter,
			call: func(clientSet kubernetes.Interface) error {
				sa.SetAnnotations()
				return kubernetes.MaybeCreateServiceAccountOrUpdateMetadata(clientSet, sa.ObjectMeta)
			},
		})

		tasks.Append(saTasks)
	}
	return tasks
}
