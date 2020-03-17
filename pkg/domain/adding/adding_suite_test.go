package adding_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAdding(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adding Service Suite")
}
