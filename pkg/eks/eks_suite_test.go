package eks_test

import (
	"testing"

	"github.com/nholuongut/eksctl/pkg/testutils"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}