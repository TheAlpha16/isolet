package response

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type APIResponse[T any] struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func Success[T any](message string, data T) APIResponse[T] {
	return newAPIResponse(StatusSuccess, message, data)
}

func Error[T any](message string, data T) APIResponse[T] {
	return newAPIResponse(StatusError, message, data)
}

func newAPIResponse[T any](status Status, message string, data T) APIResponse[T] {
	return APIResponse[T]{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
