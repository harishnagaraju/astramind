package response

type Health struct {
	Status string `json:"status"`
}

type Status struct {
	Status   string `json:"status"`
	Provider string `json:"provider"`
	Model    string `json:"model,omitempty"`
}

type Version struct {
	Version string `json:"version"`
}
