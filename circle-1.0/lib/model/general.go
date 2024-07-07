package model

type (
	Headers struct {
		TransactionID string `json:"transaction-id"`
		APIKey        string `json:"api-key"`
		ChannelID     string `json:"channel-id"`
	}

	PathParam struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	QueryParam struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	}

	RequestDetail struct {
		PathParams  []PathParam  `json:"path-params,omitempty"`
		QueryParams []QueryParam `json:"query-params,omitempty"`
		Body        any          `json:"body,omitempty"`
	}

	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	Response struct {
		TransactionID string `json:"transaction-id"`
		ChannelID     string `json:"channel-id"`
		Status        Status `json:"status"`
		Data          any    `json:"data,omitempty"`
	}

	APIDetail struct {
		Hostname  string        `json:"hostname"`
		URL       string        `json:"url"`
		Method    string        `json:"method"`
		Headers   Headers       `json:"headers"`
		Request   RequestDetail `json:"request"`
		Response  any           `json:"response"`
		Message   any           `json:"message"`
		TimeTaken string        `json:"time_taken"`
	}
)
