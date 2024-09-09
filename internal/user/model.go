package user

type User struct {
	ID     uint64
	Name   string
	Number string

	IsRegistered bool
}

type Device struct {
	ID   uint64
	Name string

	refresh string // refresh token

}

type VerifyNumberInput struct {
	Number string `json:"number"`
	Code   uint   `json:"code"`
	Token  string `json:"token"` // temporary token
}

type SignupUserInput struct {
	Number   string `json:"number" validate:"require,e164"`
	Name     string `json:"name" validate:"name"`
	Password string `json:"password" validate:"password"`
	Token    string `json:"token" validate:"uuid"`
}

type UserOutput struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}
