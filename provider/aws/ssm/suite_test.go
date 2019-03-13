package ssm_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSsm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS SSM Suite")
}
