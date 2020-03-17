package adding

import (
	"fmt"
	"net/http"
)

//FakeRepositoryAdder defines FetchCreate mock.
type FakeRepositoryAdder struct{}

//CreateRecord implements RepositoryAdder interface.
func (f *FakeRepositoryAdder) CreateRecord(record Fetch) ServiceValidation {
	txt := fmt.Sprintf("Record has been insert into fetch db.\n")
	return ServiceValidation{StorageKeyID: 0, Status: http.StatusOK, Msg: txt}
}
