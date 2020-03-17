package responding_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestResponding(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Responding Service Suite")
}
