package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailService struct {
	Email    string
	Password string
	Address  string
}

func NewEmailService(email, password, address string) *EmailService {
	return &EmailService{
		Email:    email,
		Password: password,
		Address:  address,
	}
}

func (es *EmailService) SendVerificationEmail(to, hash string) error {
	e := email.NewEmail()
	e.From = "u30390553 <u30390553@gmail.com>"
	e.To = []string{to}
	e.Subject = "Test: Email verification"
	e.HTML = fmt.Appendf(nil, `
            <h1>Перейдите по ссылке чтобы подтвердить ваш адрес электронной почты!</h1>
            <a href="https://localhost:8081/verify/%s">Подтвердить email</a>
        `, hash)

	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", es.Email, es.Password, "smtp.gmail.com"))
	if err != nil {
		return err
	}
	return nil
}

// Deprecated: Use EmailService.SendVerificationEmail instead
func SendEmail(address string, hash string, account string, password string) error {
	e := email.NewEmail()
	e.From = "u30390553 <u30390553@gmail.com>"
	e.To = []string{address}
	e.Subject = "Test: Email verification"
	e.HTML = fmt.Appendf(nil, `
            <h1>Перейдите по ссылке чтобы подтвердить ваш адрес электронной почты!</h1>
            <a href="https://localhost:8081/verify/%s">Подтвердить email</a>
        `, hash)

	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", account, password, "smtp.gmail.com"))
	if err != nil {
		return err
	}
	return nil
}
