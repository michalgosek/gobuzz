package responding

// Response stores data coming back from fetchURL routine
// used in background by Gopher.
type Response struct {
	StorageKeyID int
	Content      string
	Duration     float64
}
