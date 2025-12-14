package lib

import (
	"bytes"
	"html/template"
	"log"

	"github.com/MhmdEagel/lms-usti-be/env"
	"gopkg.in/gomail.v2"
)


func SendVerificationEmail(email, token string) error {
	// Generate Template
	t, err := template.ParseFiles("/template/index.html")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	var body bytes.Buffer
	t.Execute(&body, struct {
		Token string
	}{
		Token: token,
	})


	// Send Email
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "LMS USTI")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "LMS USTI | Verifikasi Email")
	mailer.SetBody("text/html", body.String())
	dialer := gomail.NewDialer(
		env.CONFIG_SMTP_HOST,
		587,
		env.CONFIG_EMAIL_USERNAME,
		env.CONFIG_EMAIL_PASSWORD,
	)
	dialErr := dialer.DialAndSend(mailer)
	if dialErr != nil {
		return dialErr
	}
	return nil
}
