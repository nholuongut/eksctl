package gitops

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nholuongut/eksctl/pkg/testutils"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	testutils.RegisterAndRun(t)
}
