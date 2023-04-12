package mailgo

import (
	"bytes"
	"crypto/tls"
	"html/template"

	"gopkg.in/gomail.v2"
)

type Client struct {
	dialer    *gomail.Dialer
	templates *template.Template
}

type DialerConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Templates *template.Template
}

func NewClient(c DialerConfig) (*Client, error) {
	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &Client{
		dialer:    d,
		templates: c.Templates,
	}, nil
}

type TemplateMailParams struct {
	To           []string
	Subject      string
	TemplateName string
	TemplateData any
}

func (m Client) SendTemplateMail(p TemplateMailParams) error {
	body := new(bytes.Buffer)

	if err := m.templates.ExecuteTemplate(body, p.TemplateName, p.TemplateData); err != nil {
		return err
	}

	msg := gomail.NewMessage(gomail.SetCharset("UTF-8"))

	msg.SetHeader("From", m.dialer.Username)
	msg.SetHeader("To", p.To...)
	msg.SetHeader("Subject", p.Subject)
	msg.SetBody("text/html", body.String())

	return m.dialer.DialAndSend(msg)
}
