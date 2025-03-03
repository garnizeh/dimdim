package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"net/url"

	"github.com/garnizeH/dimdim/embeded"
)

var (
	ErrFailedParseTemplate   = errors.New("failed to parse template")
	ErrFailedExecuteTemplate = errors.New("failed to execute template")
)

type mail struct {
	subject string
	to      []string
	data    any
}

func NewMailSignup(baseURL, email, name, token string) *mail {
	const endpoint = "/auth/signup/"
	url := baseURL + endpoint + url.QueryEscape(token)
	data := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  url,
	}

	const subject = "Confirm your email address"

	return &mail{
		subject: subject,
		to:      []string{email},
		data:    data,
	}
}

func NewMailPassword(baseURL, email, name, token string) *mail {
	const endpoint = "/auth/reset-password/"
	url := baseURL + endpoint + url.QueryEscape(token)
	data := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  url,
	}

	const subject = "Change your password"

	return &mail{
		subject: subject,
		to:      []string{email},
		data:    data,
	}
}

type Mailer struct {
	auth      smtp.Auth
	addr      string
	from      string
	templates *embeded.Template
}

func New(addr, host, identity, username, password string) *Mailer {
	auth := smtp.PlainAuth(identity, username, password, host)

	templates := embeded.Templates()
	templates.NewEmail("signup", "signup.tmpl")

	return &Mailer{
		auth:      auth,
		addr:      addr,
		from:      identity,
		templates: templates,
	}
}

func (m *Mailer) SendMailSignup(mail *mail) error {
	buf := new(bytes.Buffer)
	if err := m.templates.RenderEmail(buf, "signup", mail.data); err != nil {
		return fmt.Errorf("failed to render email template signup: %w", err)
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + mail.subject + "!\n"
	msg := []byte(subject + mime + "\n" + buf.String())

	return smtp.SendMail(m.addr, m.auth, m.from, mail.to, msg)
}
