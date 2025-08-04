package auth

type AuthorizationResponse struct {
	Token string `json:"token"`
	Msg   string `json:"msg"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name"`
}

type VerifyRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type VerifyResponse struct {
	SessionID string `json:"sessionID"`
}

type AuthorizationRequest struct {
	SessionID string `json:"sessionID"`
	Code      string `json:"code"`
}
