package op_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOpGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OP Suite")
}
