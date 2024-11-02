// +build integration

package integration_test

import (
	"fmt"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/nholuongut/eksctl/integration/matchers"
	. "github.com/nholuongut/eksctl/integration/runner"

	"github.com/nholuongut/eksctl/pkg/ctl/cmdutils"
)

const (
	pollInterval = 15   //seconds
	timeOut      = 1200 //seconds = 20 minutes
)

var _ = Describe("(Integration) Create & Delete before Active", func() {
	const initNG = "ng-0"
	var delBeforeActiveName string

	// initialize delBeforeActiveName (and possibly clusterName) for this test suite
	if clusterName == "" {
		clusterName = cmdutils.ClusterName("", "")
	}
	if delBeforeActiveName == "" {
		delBeforeActiveName = clusterName + "-delb4active"
	}

	Context("when creating a new cluster", func() {
		It("should not return an error", func() {
			cmd := eksctlCreateCmd.WithArgs(
				"cluster",
				"--verbose", "2",
				"--name", delBeforeActiveName,
				"--tags", "alpha.eksctl.io/description=eksctl delete before active test",
				"--without-nodegroup",
				"--version", version,
			)
			cmd.Start()
			awsSession := NewSession(region)
			Eventually(awsSession, timeOut, pollInterval).Should(
				HaveExistingCluster(delBeforeActiveName, awseks.ClusterStatusCreating, version))
		})
	})

	Context("when deleting the cluster in process of being created", func() {
		It("deleting cluster should have a zero exitcode", func() {
			cmd := eksctlDeleteClusterCmd.WithArgs(
				"--name", delBeforeActiveName,
			)
			Expect(cmd).To(RunSuccessfully())
		})
	})

	Context("after the delete of the cluster in progress has been initiated", func() {
		It("should eventually delete the EKS cluster and both CloudFormation stacks", func() {
			awsSession := NewSession(region)
			Eventually(awsSession, timeOut, pollInterval).ShouldNot(
				HaveExistingCluster(delBeforeActiveName, awseks.ClusterStatusActive, version))
			Eventually(awsSession, timeOut, pollInterval).ShouldNot(
				HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", delBeforeActiveName)))
			Eventually(awsSession, timeOut, pollInterval).ShouldNot(
				HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", delBeforeActiveName, initNG)))
		})
	})

	Context("when trying to delete the cluster again", func() {
		It("should return an a non-zero exit code", func() {
			cmd := eksctlDeleteClusterCmd.WithArgs(
				"--name", delBeforeActiveName,
			)
			Expect(cmd).ToNot(RunSuccessfully())
		})
	})
})
