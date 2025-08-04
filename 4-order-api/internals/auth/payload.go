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
	Phone    string `json:"phone" validate:"required,e164"`
}

type VerifyRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type VerifyResponse struct {
	Session Session
}

type AuthorizationBySMSRequest struct {
	SessionID string `json:"session_id"`
	Code      string `json:"code"`
}

type Session struct {
	SessionID string `json:"session_id"`
	Code      string `json:"code"`
}
