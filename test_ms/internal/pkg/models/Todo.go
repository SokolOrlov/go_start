package models

type Todo struct {
	Id       int
	Task     string
	Complete bool
}

type Todos []Todo
