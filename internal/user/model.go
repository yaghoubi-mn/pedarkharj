package user

type User struct {
	ID     uint64
	Name   string
	Number string

	IsRegistered bool
}

type Device struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`

	refresh string // refresh token

}

type VerifyNumberInput struct {
	Number string `json:"number"`
	Code   uint   `json:"code"`
	Token  string `json:"token"` // temporary token
}

type UserInput struct {
	Name string `json:"name"`
}

type UserOutput struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}
