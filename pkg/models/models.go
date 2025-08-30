package models

type Task struct {
	ID         string `json:"id"`
	Payload    string `json:"payload"`
	MaxRetries int    `json:"max_retries"`
	Status     string `json:"status"`
	Attempts   int    `json:"attempts"`
}
