package model

type DeadJobError struct {
	Attempt int
	Error   string
	Trace   string
}

type DeadJob struct {
	ID     int64
	Kind   string
	Args   string
	Errors []DeadJobError
}
