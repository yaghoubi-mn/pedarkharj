package datatypes

type EmptyStruct struct{}

func (e *EmptyStruct) Table() string {
	return ""
}

type ResponseDTO struct {
	UserErr      error
	ServerErr    error
	Data         map[string]any
	ResponseCode string
	Status       int
}
