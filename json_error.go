package guard

type JsonError struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Type   string `json:"type"`
	Status int    `json:"status"`
}
