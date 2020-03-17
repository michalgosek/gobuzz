package memory

import (
	"github.com/gobuzz/pkg/storage/memory/fetch"
	"github.com/gobuzz/pkg/storage/memory/response"
)

// ResponseFetch is an aggregate which keeps fetch and response data in memory
type ResponseFetch struct {
	Fetches   fetch.Storage
	Responses response.Storage
}
