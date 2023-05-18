package mailgo

import (
	"bytes"
	"crypto/tls"
	"fmt"
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

type SendMailParams struct {
	To      []string
	From    string
	Subject string

	// Add plain text part to the email, will be used as fallback if html is also included.
	PlainText string

	// Add html part to the email as a string, (User template params if you want to use templates passed to the client).
	HTML string

	// Template params to execute templates passed to the client.
	TemaplateParams *SendMailTemplateParams
}

type SendMailTemplateParams struct {
	Name string
	Data any
}

func (m Client) SendMail(p SendMailParams) error {

	msg := gomail.NewMessage(gomail.SetCharset("UTF-8"))

	msg.SetHeader("From", fmt.Sprintf("%q <%s>", p.From, m.dialer.Username))
	msg.SetHeader("To", p.To...)
	msg.SetHeader("Subject", p.Subject)

	if p.PlainText != "" {
		msg.AddAlternative("text/plain", p.PlainText)
	}

	if p.HTML != "" {
		msg.AddAlternative("text/html", p.HTML)
	} else if p.TemaplateParams != nil && m.templates != nil {
		body := new(bytes.Buffer)

		if err := m.templates.ExecuteTemplate(body, p.TemaplateParams.Name, p.TemaplateParams.Data); err != nil {
			return err
		}

		msg.AddAlternative("text/html", body.String())
	}

	return m.dialer.DialAndSend(msg)
}
