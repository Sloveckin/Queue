package queue

const (
	StatusOK    = "Ok"
	StatusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Ok() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(err string) Response {
	return Response{
		Status: StatusError,
		Error:  err,
	}
}
