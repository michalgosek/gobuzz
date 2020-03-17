package handlers

import (
	"fmt"
	"net/http"

	"github.com/gobuzz/pkg/domain/adding"
	"github.com/gobuzz/pkg/domain/responding"
	"github.com/gobuzz/pkg/http/load"
	"github.com/gobuzz/pkg/http/worker"
)

// HandleFetchCreate creates a single fetch and stores it in fetch repository.
func HandleFetchCreate(adder adding.Service, respsr responding.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var checkStruct load.JSONPostBody
		payloadValidation := load.PostPayloadCheck(w, r, &checkStruct)
		if payloadValidation.Status != http.StatusAccepted {
			http.Error(w, payloadValidation.Msg, payloadValidation.Status)
			return
		}

		url := *checkStruct.URL
		interval := *checkStruct.Interval

		newFetch := adding.Fetch{
			URL:      url,
			Interval: interval,
		}

		validation := adder.CreateRecord(newFetch)
		if validation.Status != http.StatusOK {
			http.Error(w, validation.Msg, validation.Status)
			return
		}

		goph := &worker.Gopher{ // Creating Gopher for background goroutine
			ID:       validation.StorageKeyID,
			URL:      url,
			Interval: interval,
		}

		msg := []byte(fmt.Sprintf(`{"id" : %d }`+"\n", validation.StorageKeyID))
		go worker.GopherRun(goph, respsr)
		w.Write(msg)
	}
}
