package worker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gobuzz/pkg/domain/responding"
	"github.com/gobuzz/pkg/format"
)

// Gopher definies task rules for working goroutine
type Gopher struct {
	ID       int
	URL      string
	Interval int
}

// GopherValidationStatus represents data stream body sending back
// to GopherRun about fetchURL state during its execution.
type GopherValidationStatus struct {
	Status int
	Msg    string
}

// fetchURL gorotuine for Gopher internal usage. Fetch the conent from URL
// mesure elapsed time from start till end of the request and pass these data to
// repository of responding service.
func fetchURL(ctxParent context.Context, goph *Gopher, respsr responding.Service, dataStream chan<- GopherValidationStatus) {

	log.Println()
	log.Printf("fetchURL[worker id:%d] - Start.\n", goph.ID)
	defer log.Printf("fetchURL[worker id:%d] - Stop.\n", goph.ID)
	defer close(dataStream)

	timeout := 5 * time.Second // GET request cancellation after 5s
	ctxChild, cancel := context.WithTimeout(ctxParent, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxChild, http.MethodGet, goph.URL, nil)
	if err != nil {
		log.Println("Error: ", err.Error())
		record := responding.Response{
			StorageKeyID: goph.ID,
			Content:      "null",
			Duration:     0,
		}
		respsr.CreateRecord(record)
		fault := GopherValidationStatus{Status: http.StatusBadRequest, Msg: err.Error()}
		dataStream <- fault
		return
	}

	start := time.Now()
	res, err := http.DefaultClient.Do(req)
	end := time.Now()
	diff := end.Sub(start).Seconds()
	elapsed := format.Duration(diff, 1000)

	if err != nil {
		log.Println("Request failed: ", err.Error())
		record := responding.Response{
			StorageKeyID: goph.ID,
			Content:      "null",
			Duration:     0,
		}
		respsr.CreateRecord(record)
		fault := GopherValidationStatus{Status: http.StatusBadRequest, Msg: err.Error()}
		dataStream <- fault
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		record := responding.Response{
			StorageKeyID: goph.ID,
			Content:      "null",
			Duration:     0,
		}
		respsr.CreateRecord(record)
		fault := GopherValidationStatus{Status: http.StatusNotFound, Msg: http.StatusText(http.StatusNotFound)}
		dataStream <- fault
		return
	}

	if res.StatusCode == http.StatusOK {
		var (
			reader  io.Reader
			resData bytes.Buffer
		)

		maxByteSize := int64(1 << 20) // 1MB limit
		reader = io.LimitReader(res.Body, maxByteSize)

		n, err := resData.ReadFrom(reader)
		if err != nil {
			log.Printf("Error reading the body: %v\n", err)
			record := responding.Response{
				StorageKeyID: goph.ID,
				Content:      "null",
				Duration:     0,
			}
			respsr.CreateRecord(record)
			fault := GopherValidationStatus{Status: http.StatusBadRequest, Msg: http.StatusText(http.StatusBadRequest)}
			dataStream <- fault
			return
		}

		record := responding.Response{
			StorageKeyID: goph.ID,
			Content:      resData.String(),
			Duration:     elapsed,
		}

		servValid := respsr.CreateRecord(record)

		log.Printf("Data content: %s read bytes: %d\n", resData.String(), n)
		log.Println("DefaultClient response recived, status code:", res.StatusCode)
		log.Println("Responser service:")
		log.Println("Status code:", servValid.Status)
		log.Printf("Validation msg: %s | response db key = %d\n", servValid.Msg, goph.ID)
		log.Println("Added record key:", servValid.StorageKeyID)
		fault := GopherValidationStatus{Status: http.StatusAccepted, Msg: "Adding record into resp db was succeed."}
		dataStream <- fault
		return
	}
	log.Println(err.Error())
	fault := GopherValidationStatus{Status: http.StatusBadRequest, Msg: "Something goes wrong."}
	dataStream <- fault
	return
}

// GopherRun is a background goroutine for fetching data for individual requests
func GopherRun(goph *Gopher, respsr responding.Service) GopherValidationStatus {

	log.Printf("Worker[id:%d] - Start\n", goph.ID)
	defer log.Printf("Worker[id:%d] - Stop\n", goph.ID)

	interval := time.Duration(goph.Interval) * time.Second
	halt := 20 * time.Minute

	ctxParent := context.Background()

	dataStream := make(chan GopherValidationStatus)
	var dataRecived GopherValidationStatus

loop:
	for {
		select {
		case <-time.After(interval):
			go fetchURL(ctxParent, goph, respsr, dataStream) // TO-DO: Think about using wg here.
		case res := <-dataStream:
			if res.Status != http.StatusAccepted {
				dataRecived = GopherValidationStatus{Status: res.Status, Msg: res.Msg}
				break loop
			}
		case <-time.After(halt):
			log.Printf("Worker[id:%d]: Timeout because of 20 min halt.", goph.ID)
			dataRecived = GopherValidationStatus{Status: http.StatusRequestTimeout, Msg: http.StatusText(http.StatusRequestTimeout)}
			break loop
		}
	}

	fmt.Println("\nChannel data: ")
	fmt.Println(dataRecived.Status)
	fmt.Println(dataRecived.Msg)
	return dataRecived
}
