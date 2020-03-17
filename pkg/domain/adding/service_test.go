package adding_test

import (
	"fmt"
	"net/http"

	. "github.com/gobuzz/pkg/domain/adding"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// testContent is an internal aggregate for creating tableTest slice.
type testContent struct {
	Fetch
	ServiceValidation
}

var _ = Describe("The adding service", func() {

	Describe("When calling CreateRecord", func() {
		var (
			data     []testContent
			adder    Service
			fetchRep FakeRepositoryAdder
		)

		BeforeEach(func() { // Configuration
			data = []testContent{
				{
					Fetch{"http://httpbin.org/range/15", 10},
					ServiceValidation{0, http.StatusOK, fmt.Sprintf("Record has been insert into fetch db.\n")},
				},
				{
					Fetch{"http://httpbin.org/range/20", 14},
					ServiceValidation{0, http.StatusOK, fmt.Sprintf("Record has been insert into fetch db.\n")},
				},
				{
					Fetch{"http://httpbin.org/delay/150", 15},
					ServiceValidation{0, http.StatusOK, fmt.Sprintf("Record has been insert into fetch db.\n")},
				},
				{
					Fetch{"https://httpbin.org/delay/3000", 16},
					ServiceValidation{0, http.StatusOK, fmt.Sprintf("Record has been insert into fetch db.\n")},
				},
			}
		})

		JustBeforeEach(func() {
			adder = NewService(&fetchRep) // Creation
		})

		Context("When fetch data is valid.", func() {
			It("Should return ID: 0, http.StatusOK, and record add db msg.", func() {
				for _, element := range data {
					serviceVal := adder.CreateRecord(element.Fetch)
					Expect(serviceVal.StorageKeyID).To(Equal(element.ServiceValidation.StorageKeyID))
					Expect(serviceVal.Status).To(Equal(element.ServiceValidation.Status))
					Expect(serviceVal.Msg).To(Equal(element.ServiceValidation.Msg))
				}
			})

			Context("When fetch data is not valid.", func() {
				BeforeEach(func() { // Configuration
					data = []testContent{
						{
							Fetch{"Woops!http://httpbin.org/range/15", 10},
							ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("URL path is not accepted.\n")},
						},
						{
							Fetch{"http://httpbin.Woops!org/range/15", 14},
							ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("URL path is not accepted.\n")},
						},
						{
							Fetch{"http://httpbin.org/range/15Woops!", 15},
							ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("URL path is not accepted.\n")},
						},
						{
							Fetch{"http://httpbin.org/delay/150", 0},
							ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("Interval value must be greater than 0.\n")},
						},
						{
							Fetch{"https://httpbin.org/range/20", -10},
							ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("Interval value must be greater than 0.\n")},
						},
						{
							Fetch{"www.google.com", 12},
							ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("URL path is not accepted.\n")},
						},
						{
							Fetch{"", 10},
							ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("URL path is not accepted.\n")},
						},
						{
							Fetch{"", -1},
							ServiceValidation{-1, http.StatusBadRequest, fmt.Sprintf("Interval and URL path are not accepted.\n")},
						},
					}
				})

				It("Should return ID: -1, http.StatusInternalServerError, and error msg.", func() {
					for _, el := range data {
						serviceVal := adder.CreateRecord(el.Fetch)
						Expect(serviceVal.StorageKeyID).To(Equal(el.ServiceValidation.StorageKeyID))
						Expect(serviceVal.Status).To(Equal(el.ServiceValidation.Status))
						Expect(serviceVal.Msg).To(Equal(el.ServiceValidation.Msg))
					}
				})
			})
		})
	})
})
