package server

import (
	"io"
	"log"
	"strconv"
	"strings"

	"fmt"

	"github.com/0xERR0R/mailcatcher/config"
	"github.com/0xERR0R/mailcatcher/sender"
	gosmtp "github.com/emersion/go-smtp"
	"github.com/veqryn/go-email/email"
)

type Backend struct {
	cfg        config.Configuration
	mailSender *sender.MailSender
}

func (bkd *Backend) NewSession(_ *gosmtp.Conn) (gosmtp.Session, error) {
	return &Session{
		cfg:        bkd.cfg,
		mailSender: bkd.mailSender,
	}, nil
}

type Session struct {
	cfg        config.Configuration
	mailSender *sender.MailSender
	from       string
	to         []string
}

func (s *Session) AuthPlain(username, password string) error {
	return nil
}

func (s *Session) Mail(from string, opts *gosmtp.MailOptions) error {
	s.from = from

	return nil
}

func (s *Session) Rcpt(to string) error {
	s.to = append(s.to, to)

	return nil
}

func (s *Session) Data(r io.Reader) error {
	log.Printf("New message from '%s' to '%s' received", s.from, s.to)

	if !isRecipientValid(s.to, s.cfg.Host) {
		log.Print("ignoring message")

		return nil
	}

	msg, err := email.ParseMessage(r)

	if err != nil {
		log.Print("error", err)

		return err
	}

	msg.Header.SetSubject(fmt.Sprintf("[MAILCATCHER] %s", msg.Header.Subject()))
	msg.Header.SetTo(fmt.Sprintf("\"%s\" <%s>", msg.Header.To()[0], s.cfg.RedirectTo))
	msg.Header.SetFrom(fmt.Sprintf("\"%s\" <%s>", "MAILCATCHER", s.cfg.SenderMail))

	s.mailSender.SendMail(msg)

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func isRecipientValid(recipients []string, host string) bool {
	for _, recipient := range recipients {
		if strings.HasSuffix(recipient, host) {
			return true
		}
	}

	return false
}

func NewServer(configuration config.Configuration) error {
	const KB = 1024

	const maxMessageSizeInMb = 20

	be := &Backend{
		cfg:        configuration,
		mailSender: sender.NewMailSender(configuration),
	}

	s := gosmtp.NewServer(be)

	s.Addr = ":" + strconv.Itoa(configuration.Port)
	s.Domain = configuration.Host

	s.MaxMessageBytes = KB * KB * maxMessageSizeInMb
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	s.AuthDisabled = true

	log.Println("Starting server at", s.Addr)

	return s.ListenAndServe()
}
