package adding

import (
	"fmt"
	"net/http"
	"regexp"
)

// RepositoryAdder provides adding functionality into fetch repository.
type RepositoryAdder interface {
	CreateRecord(fetch Fetch) ServiceValidation
}

// Service defines RepositoryAdder operation.
type Service struct {
	fetchRep RepositoryAdder
}

// ServiceValidation represetns response body sending to client
// when validation check fails.
type ServiceValidation struct {
	StorageKeyID int
	Status       int
	Msg          string
}

// CreateRecord provides adding fetch into Service repository.
func (s *Service) CreateRecord(record Fetch) ServiceValidation {

	// Validation logic...
	// pattern matching: http|https://httpbin.org/range|delay/upTo6Digits, 1st other than 0
	pattern := `^https://httpbin.org/(range|delay)/([1-9]|[1-9]{1,5})$`
	invalidPath, _ := regexp.MatchString(pattern, record.URL)

	switch {
	case record.Interval <= 0 && !invalidPath:
		txt := fmt.Sprintf("Interval and URL path are not accepted.\n")
		return ServiceValidation{StorageKeyID: -1, Status: http.StatusBadRequest, Msg: txt}
	case !invalidPath:
		txt := fmt.Sprintf("URL path is not accepted.\n")
		return ServiceValidation{StorageKeyID: -1, Status: http.StatusBadRequest, Msg: txt}
	case record.Interval <= 0:
		txt := fmt.Sprintf("Interval value must be greater than 0.\n")
		return ServiceValidation{StorageKeyID: -1, Status: http.StatusBadRequest, Msg: txt}
	}

	return s.fetchRep.CreateRecord(record)
}

// NewService creates an adding service with the necessary dependencies.
func NewService(r RepositoryAdder) Service {
	return Service{r}
}
