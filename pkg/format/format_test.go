package format_test

import (
	. "github.com/gobuzz/pkg/format"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// testContent is an internal aggregate for creating tableTest slice.
type testContent struct {
	diff   float64
	prec   float64
	result float64
}

var _ = Describe("When calling Format", func() {
	var data []testContent
	BeforeEach(func() {
		data = []testContent{
			{diff: 0.56642, prec: 100, result: 0.57},
			{diff: 3.3423, prec: 1000, result: 3.342},
			{diff: 4.56642, prec: -10, result: 4.6},
			{diff: 5.012, prec: 10, result: 5.0},
			{diff: 54.764534, prec: -1000, result: 5.0},
		}
	})

	Context("When various diff and precs values are passed.", func() {
		It("Should return 5.0 or rounded diff with set precission.", func() {
			for _, el := range data {
				res := Duration(el.diff, el.prec)
				Expect(res).To(Equal(el.result))
			}
		})
	})
})
