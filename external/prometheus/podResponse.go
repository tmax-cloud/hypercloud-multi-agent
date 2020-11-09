package promethues

type Responese struct {
	Status string        `json:"status,omitempty"`
	Data   ResponeseData `json:"data,omitempty"`
}

type ResponeseData struct {
	ResultType string                `json:"resultType,omitempty"`
	Result     []ResponeseDataResult `json:"result,omitempty"`
}

type ResponeseDataResult struct {
	Value []string `json:"value,omitempty"`
}
