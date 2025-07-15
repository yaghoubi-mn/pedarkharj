package app_shared

type ResponseDTO struct {
	UserErr      error
	ServerErr    error
	Data         map[string]any
	ResponseCode string
	Status       int
}
