package datatypes

type EmptyStruct struct{}

func (e *EmptyStruct) Table() string {
	return ""
}
