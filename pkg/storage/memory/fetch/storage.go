package fetch

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gobuzz/pkg/domain/adding"
)

// Storage represetns global storage for posted fetches
type Storage struct {
	uid  int
	db   map[int][]fetch
	init sync.Once // for mutual exlcusion of critical section
}

// CreateRecord returns an request ID after adding fetch into map storage.
func (f *Storage) CreateRecord(data adding.Fetch) adding.ServiceValidation {

	// Init once
	f.init.Do(func() {
		f.db = make(map[int][]fetch)
		f.uid = 0
	})

	fetchID := f.uid
	record := fetch{
		id:       fetchID,
		url:      data.URL,
		interval: data.Interval,
	}

	f.db[fetchID] = append(f.db[fetchID], record)
	fmt.Println(f.db) // temp for content check
	f.uid++
	return adding.ServiceValidation{StorageKeyID: fetchID, Status: http.StatusOK, Msg: "Record has been insert into fetch db."}
}
