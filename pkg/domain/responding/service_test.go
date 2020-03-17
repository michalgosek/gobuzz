package responding_test

import (
	"fmt"
	"net/http"

	. "github.com/gobuzz/pkg/domain/responding"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// testContent is an internal aggregate for creating tableTest slice.
type testContent struct {
	Response
	ServiceValidation
}

var _ = Describe("Responding Service", func() {
	Describe("When calling CreateRecord", func() {
		var (
			fakeRep FakeRepositoryAdder
			respsr  Service
			data    []testContent
		)

		JustBeforeEach(func() {
			respsr = NewService(&fakeRep) //Creation
		})

		BeforeEach(func() { // Configuration
			data = []testContent{
				{
					Response{0, "abcdefegh", 0.342},
					ServiceValidation{0, http.StatusOK, fmt.Sprintf("Record has been insert into response db.\n")},
				},
				{
					Response{1, "abcdefghij", 0.560},
					ServiceValidation{0, http.StatusOK, fmt.Sprintf("Record has been insert into response db.\n")},
				},
				{
					Response{1, "abcdefghij", 0.560},
					ServiceValidation{0, http.StatusOK, fmt.Sprintf("Record has been insert into response db.\n")},
				},
			}
		})

		Context("When response content is valid.", func() {
			It("Should return StorageKeyID, http.StatusOK, and record add db msg.", func() {
				for _, element := range data {
					serviceVal := respsr.CreateRecord(element.Response)
					Expect(serviceVal.StorageKeyID).To(Equal(element.ServiceValidation.StorageKeyID))
					Expect(serviceVal.Status).To(Equal(element.ServiceValidation.Status))
					Expect(serviceVal.Msg).To(Equal(element.ServiceValidation.Msg))
				}
			})
		})

		Describe("When response content is invalid.", func() {
			BeforeEach(func() { // Configuration
				data = []testContent{
					{
						Response{-1, "abcdefegh", 0.342},
						ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("StorageKeyID must be greater or equal 0.\n")},
					},
					{
						Response{1, "", 0.342},
						ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("Response string must be in range (0, 102402) characters.\n")},
					},
					{
						Response{1, "null", 5.1},
						ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("Response duration cannot be longer than 5s.\n")},
					},
					{
						Response{1, "abcdefgh", 5.1},
						ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("Response duration longer than 5s should return null as content.\n")},
					},
				}
			})

			Context("When response content is invalid.", func() {
				It("Should return StorageKeyID as -1, http.StatusBadRequest, and record add db msg.", func() {
					for _, element := range data {
						serviceVal := respsr.CreateRecord(element.Response)
						Expect(serviceVal.StorageKeyID).To(Equal(element.ServiceValidation.StorageKeyID))
						Expect(serviceVal.Status).To(Equal(element.ServiceValidation.Status))
						Expect(serviceVal.Msg).To(Equal(element.ServiceValidation.Msg))
					}
				})
			})

		})

	})
})
