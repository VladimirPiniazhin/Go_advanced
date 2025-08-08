package verify

type SendLinkRequest struct {
	Address string `json:"address" validate:"required,email"`
}
type VerifyRequest struct {
	Hash string `json:"hash"`
}
type VerifyResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
