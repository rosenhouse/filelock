package filelock_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFilelock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filelock Suite")
}
