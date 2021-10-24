package mock

import (
	"github.com/terdia/greenlight/internal/mailer"
)

type mailerMock struct {
}

func NewMailerMock() mailer.Mailer {

	return &mailerMock{}

}

func (m mailerMock) Send(recipient, templateFile string, data interface{}) error {

	return nil
}
