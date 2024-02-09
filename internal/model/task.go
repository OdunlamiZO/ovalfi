package model

type Status string

const (
	TODO        Status = "Todo"
	IN_PROGRESS Status = "In Progress"
	COMPLETED   Status = "Completed"
)

type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
	User        string `json:"-"`
}
