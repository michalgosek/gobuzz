package responding

import (
	"fmt"
	"net/http"
)

// RepositoryAdder provides adding functionality to requests response repository.
type RepositoryAdder interface {
	CreateRecord(record Response) ServiceValidation
}

// Service defines RepositoryAdder operation.
type Service struct {
	reqsRep RepositoryAdder
}

// ServiceValidation represetns response body sending to client
// when validation check fails.
type ServiceValidation struct {
	StorageKeyID int
	Status       int
	Msg          string
}

// CreateRecord provides adding request into Service repository.
func (s *Service) CreateRecord(record Response) ServiceValidation {

	//.. Validation logic
	switch {
	case record.StorageKeyID < 0:
		txt := fmt.Sprintf("StorageKeyID must be greater or equal 0.\n")
		return ServiceValidation{StorageKeyID: -1, Status: http.StatusBadRequest, Msg: txt}
	case record.Duration > 5.0 && record.Content != "null":

		txt := fmt.Sprintf("Response duration longer than 5s should return null as content.\n")
		return ServiceValidation{StorageKeyID: -1, Status: http.StatusBadRequest, Msg: txt}

	case len(record.Content) <= 0 || len(record.Content) > 102402:
		txt := fmt.Sprintf("Response string must be in range (0, 102402) characters.\n")
		return ServiceValidation{StorageKeyID: -1, Status: http.StatusBadRequest, Msg: txt}

	case record.Duration > 5.0:
		txt := fmt.Sprintf("Response duration cannot be longer than 5s.\n")
		return ServiceValidation{StorageKeyID: -1, Status: http.StatusBadRequest, Msg: txt}
	}

	return s.reqsRep.CreateRecord(record)
}

// NewService creates an adding service with the necessary dependencies
func NewService(r RepositoryAdder) Service {
	return Service{r}
}
