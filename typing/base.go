package typing

type Home struct {
	Status      bool   `json:"status"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Twitter     string `json:"twitter"`
}

type Health struct {
	Name    string `json:"name"`
	Status  bool   `json:"status"`
	Version string `json:"version"`
}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
