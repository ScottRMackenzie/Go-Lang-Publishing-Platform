package types

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserAccountCreationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
