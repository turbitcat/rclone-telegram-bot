package rclone

// json struct:
// _async = true
type JobRespons struct {
	JobID int
}

// json struct:
// error
type ErrorRespons struct {
	Error  string
	Status int
}
