package adding

// Fetch defines incoming fetch request JSON data
type Fetch struct {
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}
