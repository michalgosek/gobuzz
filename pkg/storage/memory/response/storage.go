package response

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gobuzz/pkg/domain/responding"
	"github.com/levenlabs/golib/timeutil"
)

// Storage represetns internal storage of fetch request
type Storage struct { // Implements RepositoryAdder interface
	uid  int
	db   map[int][]response
	init sync.Once // for mutual exlcusion of critical section
}

// CreateRecord provides adding record funcionality into response storge
// for each fetch request.
func (s *Storage) CreateRecord(data responding.Response) responding.ServiceValidation {

	// Init once
	s.init.Do(func() {
		s.uid = 0
		s.db = make(map[int][]response)
	})

	record := response{
		response:  data.Content,
		duration:  data.Duration,
		createdAt: fmt.Sprintf("%.5f", timeutil.TimestampNow().Float64()),
	}

	key := data.StorageKeyID
	s.db[key] = append(s.db[key], record)
	fmt.Println(s.db[key]) // temp for content check
	s.uid++
	return responding.ServiceValidation{StorageKeyID: s.uid, Status: http.StatusOK, Msg: "Record has been insert into response db."}
}
