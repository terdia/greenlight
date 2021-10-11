package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"

	"github.com/terdia/greenlight/config"
)

//go:embed "templates"
var templateFS embed.FS

const (
	dialerTimeout = 5 * time.Second
	retries       = 3
	retryAfter    = 20000 * time.Millisecond
)

type Mailer interface {
	Send(recipient, templateFile string, data interface{}) error
}

type mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(cfg config.Smtp) Mailer {

	dialer := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	dialer.Timeout = dialerTimeout

	return &mailer{
		dialer: dialer,
		sender: cfg.Sender,
	}

}

func (m mailer) Send(recipient, templateFile string, data interface{}) error {

	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	for i := 1; i <= retries; i++ {
		err = m.dialer.DialAndSend(msg)
		if nil == err {
			return nil
		}

		//rety after 20 seconds
		time.Sleep(retryAfter)
	}

	return nil
}
