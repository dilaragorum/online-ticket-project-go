package auth

type Credentials struct {
	Password string `json:"password"`
	UserName string `json:"user_name"`
}
