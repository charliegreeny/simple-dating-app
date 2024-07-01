package token

type LoginInput struct {
	ID       string `json:"id"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email"`
}

type LoginOutput struct {
	Token string `json:"token"`
}
