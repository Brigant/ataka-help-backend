package core

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SetResponse(code int, message string) Response {
	return Response{
		Status:  code,
		Message: message,
	}
}
