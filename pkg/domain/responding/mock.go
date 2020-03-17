package responding

import (
	"fmt"
	"net/http"
)

// FakeRepositoryAdder ......
type FakeRepositoryAdder struct{}

// CreateRecord ....
func (f *FakeRepositoryAdder) CreateRecord(record Response) ServiceValidation {
	txt := fmt.Sprintf("Record has been insert into response db.\n")
	return ServiceValidation{StorageKeyID: 0, Status: http.StatusOK, Msg: txt}
}
