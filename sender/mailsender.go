package sender

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/0xERR0R/mailcatcher/config"
	"github.com/veqryn/go-email/email"
)

type MailSender struct {
	cfg config.Configuration
}

func NewMailSender(cfg config.Configuration) *MailSender {
	return &MailSender{
		cfg: cfg,
	}
}

func (m *MailSender) SendMail(msg *email.Message) {
	if err := msg.Save(); err != nil {
		log.Printf("can't save message: %s", err)

		return
	}

	b, err := msg.Bytes()

	if err != nil {
		log.Printf("can't convert message: %s", err)

		return
	}

	err = smtp.SendMail(fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port),
		smtp.PlainAuth("", m.cfg.SMTPUser, m.cfg.SMTPPassword, m.cfg.SMTPHost),
		m.cfg.SenderMail, []string{m.cfg.RedirectTo}, b)

	if err != nil {
		log.Printf("smtp error: %s", err)

		return
	}
}
