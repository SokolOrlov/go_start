package api

type Todo struct {
	Task     string `json:"task"`
	Complete bool   `json:"complete"`
}
